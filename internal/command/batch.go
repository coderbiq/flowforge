package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"flowforge/internal/core"
)

type batchManifest struct {
	Proposal string      `yaml:"proposal,omitempty"`
	Cards    []batchCard `yaml:"cards"`
}

type batchCard struct {
	Ref    string   `yaml:"ref,omitempty"`
	Type   string   `yaml:"type"`
	Title  string   `yaml:"title"`
	Status string   `yaml:"status,omitempty"`
	Body   string   `yaml:"body,omitempty"`
	Links  []string `yaml:"links,omitempty"`
	Tags   []string `yaml:"tags,omitempty"`
	Domain string   `yaml:"domain,omitempty"`
}

type batchCreatedCard struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type batchResult struct {
	Created []batchCreatedCard `json:"created"`
	Errors  []batchError       `json:"errors,omitempty"`
}

type batchError struct {
	Index int    `json:"index"`
	Error string `json:"error"`
}

func newCardCreateBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch <file>",
		Short: "Create multiple cards from a YAML manifest",
		Long: `Create multiple cards from a YAML manifest file.

Manifest format:
  proposal: "CR26062001"        # optional, auto-resolves if omitted
  cards:
    - ref: "str-core"            # optional cross-reference name
      type: structure
      title: "Architecture Index"
      status: active
      body: |
        Multi-line body content.
      links:
        - "FIND-xxx:references"

    - type: convention
      title: "Naming Rules"
      status: draft
      body: |
        Rules content.
      links:
        - "FIND-xxx:references"
        - "@str-core:indexes"     # cross-ref + auto structure add

Use @ref in links to reference another card in the same batch.
The indexes relation automatically performs structure add.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			batchFile := args[0]

			data, err := os.ReadFile(batchFile)
			if err != nil {
				return fmt.Errorf("reading batch file: %w", err)
			}

			var manifest batchManifest
			if err := yaml.Unmarshal(data, &manifest); err != nil {
				return fmt.Errorf("parsing batch YAML: %w", err)
			}

			if len(manifest.Cards) == 0 {
				return fmt.Errorf("batch manifest must contain at least one card")
			}

			if err := validateBatchManifest(&manifest); err != nil {
				return err
			}

			store, err := currentCardStore()
			if err != nil {
				return err
			}

			resolvedRefs := map[string]string{}
			var createdCards []batchCreatedCard
			var errors []batchError

			for i, card := range manifest.Cards {
				ct := core.CardType(card.Type)
				newCard := core.NewCard(ct, card.Title)
				if card.Status != "" {
					cs := core.CardStatus(card.Status)
					if !cs.Valid() {
						errors = append(errors, batchError{Index: i, Error: fmt.Sprintf("invalid status: %s", card.Status)})
						continue
					}
					newCard.Status = cs
				}
				newCard.Body = card.Body
				newCard.Tags = card.Tags
				newCard.Domain = card.Domain

				resolvedProposalID, err := resolveDefaultProposalID(manifest.Proposal, ct)
				if err != nil {
					errors = append(errors, batchError{Index: i, Error: err.Error()})
					continue
				}
				proposalTs := proposalTimestamp(resolvedProposalID)
				newCard.ID = core.GenerateCardID(ct, proposalTs)

				addProposalOwnershipLink(newCard, resolvedProposalID)

				for _, link := range card.Links {
					resolvedLink := resolveRef(link, resolvedRefs)
					parts := strings.SplitN(resolvedLink, ":", 2)
					target := parts[0]
					relation := "references"
					if len(parts) == 2 {
						relation = parts[1]
					}
					if _, err := store.ReadCard(target); err != nil {
						errors = append(errors, batchError{Index: i, Error: fmt.Sprintf("link target %s not found: %v", target, err)})
						continue
					}
					newCard.AddLink(target, relation)
				}

				if len(errors) > 0 {
					continue
				}

				if len(newCard.Links) == 0 && resolvedProposalID == "" {
					errors = append(errors, batchError{Index: i, Error: "card requires at least one outbound link; add --links or set proposal"})
					continue
				}

				upsertLinksSection(store, newCard)

				_, err = store.CreateCard(newCard, resolvedProposalID)
				if err != nil {
					errors = append(errors, batchError{Index: i, Error: err.Error()})
					continue
				}

				if card.Ref != "" {
					resolvedRefs[card.Ref] = newCard.ID
				}

				for _, link := range newCard.Links {
					if link.Relation == "indexes" {
						doStructureAdd(store, link.Target, newCard.ID)
					}
				}

				createdCards = append(createdCards, batchCreatedCard{
					ID:    newCard.ID,
					Type:  string(newCard.Type),
					Title: newCard.Title,
				})
			}

			out := cmd.OutOrStdout()
			result := batchResult{
				Created: createdCards,
				Errors:  errors,
			}

			if isJSONOutput(cmd) {
				data, _ := json.Marshal(result)
				fmt.Fprint(out, string(data))
			} else {
				if len(createdCards) > 0 {
					fmt.Fprintf(out, "✓ Created %d card(s):\n", len(createdCards))
					for _, c := range createdCards {
						fmt.Fprintf(out, "  %s %s - %s\n", c.Type, c.ID, c.Title)
					}
				}
				if len(errors) > 0 {
					fmt.Fprintf(out, "✗ %d error(s):\n", len(errors))
					for _, e := range errors {
						fmt.Fprintf(out, "  [%d] %s\n", e.Index, e.Error)
					}
				}
			}

			if len(errors) > 0 {
				return fmt.Errorf("batch completed with %d error(s)", len(errors))
			}
			return nil
		},
	}

	return cmd
}

func validateBatchManifest(m *batchManifest) error {
	refs := map[string]bool{}
	for i, card := range m.Cards {
		if card.Type == "" {
			return fmt.Errorf("card %d: --type is required", i)
		}
		ct := core.CardType(card.Type)
		if !ct.Valid() {
			return fmt.Errorf("card %d: invalid type: %s", i, card.Type)
		}
		if card.Title == "" {
			return fmt.Errorf("card %d: --title is required", i)
		}
		if card.Ref != "" {
			if refs[card.Ref] {
				return fmt.Errorf("card %d: duplicate ref %q", i, card.Ref)
			}
			refs[card.Ref] = true
		}
		for _, link := range card.Links {
			if strings.HasPrefix(link, "@") {
				refName := strings.SplitN(link[1:], ":", 2)[0]
				if !refs[refName] && !isForwardRef(m.Cards, i, refName) {
					return fmt.Errorf("card %d: ref %q not defined in batch", i, refName)
				}
			}
		}
	}
	return nil
}

func isForwardRef(cards []batchCard, currentIdx int, refName string) bool {
	for j := currentIdx + 1; j < len(cards); j++ {
		if cards[j].Ref == refName {
			return true
		}
	}
	return false
}

func resolveRef(link string, resolvedRefs map[string]string) string {
	if !strings.HasPrefix(link, "@") {
		return link
	}
	rest := link[1:]
	parts := strings.SplitN(rest, ":", 2)
	refName := parts[0]
	if actualID, ok := resolvedRefs[refName]; ok {
		if len(parts) == 2 {
			return actualID + ":" + parts[1]
		}
		return actualID
	}
	return link
}

func doStructureAdd(store *core.CardStore, structureID, cardID string) {
	_ = store.UpdateCardWithLock(structureID, func(card *core.Card) error {
		if card.Type != core.CardTypeStructure {
			return nil
		}
		card.AddLink(cardID, "indexes")
		refreshedBody, _ := refreshStructureEntriesBody(store, card)
		card.Body = refreshedBody
		return nil
	})
}