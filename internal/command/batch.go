package command

import (
	"encoding/json"
	"fmt"
	"io"
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

func batchPhaseError(phase string, index int, card batchCard, locator, err string) batchError {
	ref := ""
	if card.Ref != "" {
		ref = fmt.Sprintf(", ref %q", card.Ref)
	}
	return batchError{Index: index, Error: fmt.Sprintf("%s: manifest[%d]%s%s: %s", phase, index, ref, locator, err)}
}

func newCardCreateBatchCmd() *cobra.Command {
	var manifestInline string

	cmd := &cobra.Command{
		Use:   "batch [<file|-]",
		Short: "Create multiple cards from a YAML manifest",
		Long: `Create multiple cards from a YAML manifest file, stdin, or inline.

With file or stdin:

  flowforge card batch manifest.yaml
  flowforge card batch -

With --manifest for inline YAML (use \n for newlines):

  flowforge card batch --manifest "cards:\n  - type: feature\n    title: Feature Card"

Manifest format:
  proposal: "CR26062001"        # optional, auto-resolves if omitted
  cards:
    - ref: "feature-core"        # optional cross-reference name
      type: feature
      title: "Architecture Feature"
      status: draft
      body: Multi-line body with \n for newlines
      links:
        - "FIND-xxx:references"

    - type: decision
      title: "Naming Rules"
      status: draft
      body: Rules content with \n newlines
      links:
        - "FIND-xxx:references"
        - "@feature-core:references" # cross-reference within this batch

Use @ref in links to reference another card in the same batch.
The batch command only creates current-v3 cards; Proposal STR metadata is managed internally.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var data []byte

			if manifestInline != "" {
				unescaped := unescapeBody(manifestInline)
				data = []byte(unescaped)
			} else if len(args) == 1 && args[0] == "-" {
				var err error
				data, err = io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("reading batch manifest from stdin: %w", err)
				}
			} else if len(args) == 1 {
				var err error
				data, err = os.ReadFile(args[0])
				if err != nil {
					return fmt.Errorf("reading batch file: %w", err)
				}
			} else {
				return fmt.Errorf("requires either --manifest or a file/stdin argument")
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

			// Phase 1 — create all cards first, resolve refs only for
			// non-batch internal links.  Skip link pre-validation here so
			// that @ref targets that point to other batch cards can be
			// checked after all cards exist.
			rawRefs := map[string]string{}   // ref -> ID (populated as each card is created)
			pendingCards := map[int]string{} // manifest index -> created ID (for Phase 2 link attachment)
			var createdCards []batchCreatedCard
			var errors []batchError

			for i, card := range manifest.Cards {
				ct := core.CardType(card.Type)
				newCard := core.NewCard(ct, card.Title)
				if card.Status != "" {
					cs := core.CardStatus(card.Status)
					if !cs.Valid() {
						errors = append(errors, batchPhaseError("phase 1", i, card, "", fmt.Sprintf("invalid status: %s", card.Status)))
						continue
					}
					newCard.Status = cs
				}
				newCard.Body = card.Body
				newCard.Tags = card.Tags
				newCard.Domain = card.Domain

				resolvedProposalID, err := resolveDefaultProposalID(manifest.Proposal, ct)
				if err != nil {
					errors = append(errors, batchPhaseError("phase 1", i, card, "", err.Error()))
					continue
				}
				newCard.ID, err = store.NextCardID(ct, resolvedProposalID)
				if err != nil {
					errors = append(errors, batchPhaseError("phase 1", i, card, fmt.Sprintf(", ID %s", newCard.ID), err.Error()))
					continue
				}

				addProposalOwnershipLink(newCard, resolvedProposalID)

				// Phase 2 — validate and resolve links AFTER all cards exist
				// (deferred; see Phase 2 below)
				pendingCards[i] = newCard.ID

				upsertLinksSection(store, newCard)

				_, err = store.CreateCard(newCard, resolvedProposalID)
				if err != nil {
					errors = append(errors, batchPhaseError("phase 1", i, card, fmt.Sprintf(", ID %s", newCard.ID), err.Error()))
					delete(pendingCards, i)
					continue
				}

				if card.Ref != "" {
					rawRefs[card.Ref] = newCard.ID
				}

				createdCards = append(createdCards, batchCreatedCard{
					ID:    newCard.ID,
					Type:  string(newCard.Type),
					Title: newCard.Title,
				})
			}

			// Phase 2 — resolve @ref links, validate targets, add links, do structure add
			linkErrors := map[int][]string{}
			var allLinkTargets []string
			targetOwners := map[string]int{}
			for idx, id := range pendingCards {
				mcard := &manifest.Cards[idx]

				card, err := store.ReadCard(id)
				if err != nil {
					linkErrors[idx] = append(linkErrors[idx], fmt.Sprintf("phase 2: manifest[%d], ID %s: cannot read card: %v", idx, id, err))
					continue
				}

				for _, link := range mcard.Links {
					resolvedLink := resolveRef(link, rawRefs)
					parts := strings.SplitN(resolvedLink, ":", 2)
					target := parts[0]
					relation := "references"
					if len(parts) == 2 {
						relation = parts[1]
					}
					if relation == "indexes" {
						linkErrors[idx] = append(linkErrors[idx], fmt.Sprintf("phase 2: manifest[%d], ID %s, target %s: indexes relation is reserved for Proposal control-plane metadata", idx, id, target))
						continue
					}

					// After Phase 1 all non-@ref targets should already exist.
					// @ref targets that still contain '@' mean the ref was not resolved.
					if strings.HasPrefix(target, "@") {
						linkErrors[idx] = append(linkErrors[idx], fmt.Sprintf("phase 2: manifest[%d], ID %s, ref %q: unresolved @ref", idx, id, link))
						continue
					}

					if _, terr := store.ReadCard(target); terr != nil {
						linkErrors[idx] = append(linkErrors[idx], fmt.Sprintf("phase 2: manifest[%d], ID %s, target %s: link target not found: %v", idx, id, target, terr))
						continue
					}
					card.AddLink(target, relation)
					allLinkTargets = append(allLinkTargets, target)
					if _, exists := targetOwners[target]; !exists {
						targetOwners[target] = idx
					}
				}

				if len(card.Links) == 0 && manifest.Proposal == "" {
					linkErrors[idx] = append(linkErrors[idx], fmt.Sprintf("phase 2: manifest[%d], ID %s: card requires at least one outbound link; add --links or set proposal", idx, id))
				}

				upsertLinksSection(store, card)

				if err := store.UpdateCardWithLock(id, func(uc *core.Card) error {
					uc.Links = card.Links
					uc.Body = card.Body
					return nil
				}); err != nil {
					linkErrors[idx] = append(linkErrors[idx], fmt.Sprintf("phase 2: manifest[%d], ID %s: updating links: %v", idx, id, err))
				}

			}

			if err := refreshTargetCardsNavigation(store, allLinkTargets); err != nil {
				idx := 0
				for target, owner := range targetOwners {
					if strings.Contains(err.Error(), target) {
						idx = owner
						break
					}
				}
				linkErrors[idx] = append(linkErrors[idx], fmt.Sprintf("phase 2: manifest[%d]: navigation: %v", idx, err))
			}

			// Merge Phase 2 link errors into main errors list after navigation refresh,
			// keeping the report ordered by manifest index.
			for i := 0; i < len(manifest.Cards); i++ {
				if errs, ok := linkErrors[i]; ok {
					for _, e := range errs {
						errors = append(errors, batchError{Index: i, Error: e})
					}
				}
			}

			out := cmd.OutOrStdout()
			result := batchResult{
				Created: createdCards,
				Errors:  errors,
			}

			if isJSONOutput(cmd) {
				data, err := json.Marshal(result)
				if err != nil {
					return fmt.Errorf("encoding batch result: %w", err)
				}
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

	cmd.Flags().StringVar(&manifestInline, "manifest", "", "Inline YAML manifest (use \\n for newlines, no shell redirect needed)")

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
