package core

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) AddError(field, message string) {
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) String() string {
	if !r.HasErrors() {
		return "valid"
	}
	var msgs []string
	for _, e := range r.Errors {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "; ")
}

func ValidateCard(card *Card) *ValidationResult {
	result := &ValidationResult{}

	if card.ID == "" {
		result.AddError("id", "required")
	} else {
		validateCardID(card.ID, card.Type, result)
	}

	if card.Title == "" {
		result.AddError("title", "required")
	}

	if !card.Type.Valid() {
		result.AddError("type", fmt.Sprintf("invalid type: %s", card.Type))
	}

	if !card.Status.Valid() {
		result.AddError("status", fmt.Sprintf("invalid status: %s", card.Status))
	}

	if !card.Importance.Valid() {
		result.AddError("importance", fmt.Sprintf("invalid importance: %s", card.Importance))
	}

	if card.Created.IsZero() {
		result.AddError("created", "required")
	}

	if card.Updated.IsZero() {
		result.AddError("updated", "required")
	}

	for i, link := range card.Links {
		if link.Target == "" {
			result.AddError(fmt.Sprintf("links[%d].target", i), "required")
		}
		if link.Relation == "" {
			result.AddError(fmt.Sprintf("links[%d].relation", i), "required")
		} else if !IsValidRelation(link.Relation) {
			result.AddError(fmt.Sprintf("links[%d].relation", i), fmt.Sprintf("invalid relation: %s", link.Relation))
		}
	}

	return result
}

func validateCardID(id string, cardType CardType, result *ValidationResult) {
	parts := strings.Split(id, "-")
	if len(parts) < 2 {
		result.AddError("id", "must have at least 2 parts separated by -")
		return
	}

	prefix := parts[0]
	expectedPrefix := cardType.Prefix()
	if expectedPrefix != "" && prefix != expectedPrefix {
		result.AddError("id", fmt.Sprintf("prefix mismatch: expected %s for type %s, got %s", expectedPrefix, cardType, prefix))
	}

	if cardType == CardTypeTask {
		if len(parts) < 3 {
			result.AddError("id", "task ID must have at least 3 parts: TASK-{proposalTs}-{type}-{taskTs}")
		} else if len(parts) >= 3 {
			taskType := parts[2]
			if !isValidTaskType(taskType) && len(taskType) == 1 {
				result.AddError("id", fmt.Sprintf("invalid task type letter: %s (expected a/i/t/d/f/r/c)", taskType))
			}
		}
	}
}

func IsValidRelation(relation string) bool {
	validRelations := map[string]bool{
		"belongs_to":  true,
		"references":  true,
		"extends":     true,
		"refines":     true,
		"contradicts": true,
		"supersedes":  true,
		"supports":    true,
		"questions":   true,
		"related":     true,
		"implements":  true,
		"satisfies":   true,
		"requires":    true,
		"blocks":      true,
		"produced":    true,
		"indexes":     true,
		"decomposes":  true,
		"analyzes":    true,
		"designs":     true,
		"constrains":  true,
		"records":     true,
		"discovers":   true,
	}
	return validRelations[relation]
}

func isValidRelation(relation string) bool {
	return IsValidRelation(relation)
}

func isValidTaskType(taskType string) bool {
	validTypes := map[string]bool{
		"a": true,
		"i": true,
		"t": true,
		"d": true,
		"f": true,
		"r": true,
		"c": true,
	}
	return validTypes[taskType]
}

func ValidateCardFile(filePath string) *ValidationResult {
	card, err := ParseCardFile(filePath)
	if err != nil {
		result := &ValidationResult{}
		result.AddError("file", fmt.Sprintf("failed to parse: %v", err))
		return result
	}

	result := ValidateCard(card)

	filename := strings.TrimSuffix(filepath.Base(filePath), ".md")
	if !filenameMatchesCardID(filename, card.ID) {
		result.AddError("filename", fmt.Sprintf("mismatch: filename must start with card id %s, got %s", card.ID, filename))
	}

	return result
}

func ValidateCardFileInStore(filePath string, store *CardStore) *ValidationResult {
	result := ValidateCardFile(filePath)

	card, err := ParseCardFile(filePath)
	if err != nil {
		return result
	}

	for i, link := range card.Links {
		if _, err := store.ReadCard(link.Target); err != nil {
			result.AddError(fmt.Sprintf("links[%d].target", i), fmt.Sprintf("target not found: %s", link.Target))
		}
	}

	for _, target := range collectWikiLinkTargets(card.Body) {
		result.AddError("body.wikilink", fmt.Sprintf("wikilink is not supported; use a standard Markdown link for %s", target))
	}

	for _, target := range collectMarkdownLinkTargets(card.Body) {
		if err := validateMarkdownLinkTarget(filePath, target); err != nil {
			result.AddError("body.link", err.Error())
		}
	}

	if requiresOutboundLink(card) && len(card.Links) == 0 {
		result.AddError("links", "at least one outbound frontmatter link is required")
	}

	if card.Type == CardTypeStructure && strings.HasPrefix(card.ID, "STR-") && strings.Contains(card.ID, "-REQ") {
		for i, link := range card.Links {
			if link.Relation != "indexes" {
				continue
			}
			targetCard, err := store.ReadCard(link.Target)
			if err != nil {
				continue
			}
			if targetCard.Type != CardTypeRequirement && targetCard.Type != CardTypeStructure {
				result.AddError(fmt.Sprintf("links[%d].target", i), fmt.Sprintf("proposal requirement index can only index requirement or structure cards, got %s", targetCard.Type))
			}
		}
	}

	return result
}

var wikiLinkPattern = regexp.MustCompile(`\[\[([^\]|#]+)(?:#[^\]|]+)?(?:\|[^\]]+)?\]\]`)
var markdownLinkPattern = regexp.MustCompile(`!?\[[^\]]*\]\(([^)]+)\)`)

func filenameMatchesCardID(filename string, cardID string) bool {
	if filename == cardID {
		return true
	}
	return strings.HasPrefix(filename, cardID+"_")
}

func requiresOutboundLink(card *Card) bool {
	if card == nil {
		return false
	}
	if card.Type == CardTypeProposal || card.ID == "STR-HOME" {
		return false
	}
	return true
}

func collectWikiLinkTargets(body string) []string {
	matches := wikiLinkPattern.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return nil
	}

	targets := make([]string, 0, len(matches))
	seen := map[string]bool{}
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		target := strings.TrimSpace(match[1])
		if target == "" || seen[target] {
			continue
		}
		seen[target] = true
		targets = append(targets, target)
	}
	return targets
}

func collectMarkdownLinkTargets(body string) []string {
	matches := markdownLinkPattern.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return nil
	}

	targets := make([]string, 0, len(matches))
	seen := map[string]bool{}
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		target := strings.TrimSpace(match[1])
		target = strings.Trim(target, "<>")
		if target == "" || shouldSkipMarkdownLinkTarget(target) || seen[target] {
			continue
		}
		seen[target] = true
		targets = append(targets, target)
	}
	return targets
}

func shouldSkipMarkdownLinkTarget(target string) bool {
	if strings.HasPrefix(target, "#") {
		return true
	}
	if strings.HasPrefix(target, "mailto:") {
		return true
	}
	schemeIdx := strings.Index(target, ":")
	if schemeIdx > 0 {
		scheme := target[:schemeIdx]
		for _, r := range scheme {
			if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '+' && r != '-' && r != '.' {
				return false
			}
		}
		return true
	}
	return false
}

func validateMarkdownLinkTarget(fromFile string, target string) error {
	pathPart := strings.SplitN(target, "#", 2)[0]
	pathPart = strings.SplitN(pathPart, "?", 2)[0]
	pathPart = strings.TrimSpace(pathPart)
	if pathPart == "" {
		return nil
	}
	if filepath.IsAbs(pathPart) {
		if _, err := os.Stat(pathPart); err != nil {
			return fmt.Errorf("target not found: %s", target)
		}
		return nil
	}
	resolved := filepath.Clean(filepath.Join(filepath.Dir(fromFile), filepath.FromSlash(pathPart)))
	if _, err := os.Stat(resolved); err != nil {
		return fmt.Errorf("target not found: %s", target)
	}
	return nil
}
