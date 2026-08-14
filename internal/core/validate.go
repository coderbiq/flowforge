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

	validateComplexAnalysisBody(card, result)

	return result
}

var analysisModePattern = regexp.MustCompile(`(?mi)^\s*<!--\s*analysis-mode:\s*complex\s*-->\s*$`)
var analysisWorkIDPattern = regexp.MustCompile(`(?mi)^\s*<!--\s*analysis-work-id:\s*([A-Za-z0-9][A-Za-z0-9._:-]*)\s*-->\s*$`)
var findingIDPattern = regexp.MustCompile(`\bFIND-[A-Za-z0-9][A-Za-z0-9-]*\b`)

var complexFeatureSections = []string{
	"Objective",
	"Current Understanding",
	"Evidence",
	"Working Design",
	"Rejected or Revised Assumptions",
	"Open Questions",
	"Next Investigation",
}

var complexFindingSections = []string{
	"Evidence",
	"Source",
	"Impact",
	"Open Questions",
}

func validateComplexAnalysisBody(card *Card, result *ValidationResult) {
	if card == nil || !analysisModePattern.MatchString(card.Body) {
		return
	}

	var sections []string
	switch card.Type {
	case CardTypeFeature:
		sections = complexFeatureSections
	case CardTypeFinding:
		sections = complexFindingSections
		match := analysisWorkIDPattern.FindStringSubmatch(card.Body)
		if len(match) != 2 {
			result.AddError("body.analysis-work-id", "complex analysis FIND requires <!-- analysis-work-id: <stable-id> -->")
		}
	default:
		result.AddError("body.analysis-mode", "complex analysis mode is only valid for FEATURE and FIND cards")
		return
	}

	for _, name := range sections {
		content, ok := markdownSection(card.Body, name)
		if !ok {
			result.AddError("body.analysis."+analysisFieldName(name), fmt.Sprintf("missing required section: ## %s", name))
			continue
		}
		if isAnalysisPlaceholder(content) {
			result.AddError("body.analysis."+analysisFieldName(name), fmt.Sprintf("section ## %s must contain a recoverable value; use None when intentionally empty", name))
		}
	}
}

func markdownSection(body, name string) (string, bool) {
	heading := regexp.MustCompile(`(?m)^##\s+` + regexp.QuoteMeta(name) + `\s*$`)
	match := heading.FindStringIndex(body)
	if match == nil {
		return "", false
	}
	contentStart := match[1]
	nextHeading := regexp.MustCompile(`(?m)^#{1,2}\s+`).FindStringIndex(body[contentStart:])
	contentEnd := len(body)
	if nextHeading != nil {
		contentEnd = contentStart + nextHeading[0]
	}
	return strings.TrimSpace(body[contentStart:contentEnd]), true
}

func isAnalysisPlaceholder(content string) bool {
	cleaned := strings.TrimSpace(content)
	cleaned = strings.TrimSpace(strings.TrimLeft(cleaned, "-* "))
	cleaned = strings.Trim(cleaned, "`[]()")
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		return true
	}
	switch strings.ToLower(cleaned) {
	case "tbd", "todo", "pending", "待补充", "待定":
		return true
	default:
		return false
	}
}

func analysisFieldName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "-")
}

func validateCardID(id string, cardType CardType, result *ValidationResult) {
	parts := strings.Split(id, "-")
	if len(parts) < 2 {
		result.AddError("id", "must have at least 2 parts separated by -")
		return
	}

	prefix := parts[0]
	expectedPrefix := currentCardPrefix(cardType)
	if expectedPrefix != "" && prefix != expectedPrefix {
		result.AddError("id", fmt.Sprintf("prefix mismatch: expected %s for type %s, got %s", expectedPrefix, cardType, prefix))
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

	validateComplexAnalysisReferences(card, store, result)

	return result
}

func validateComplexAnalysisReferences(card *Card, store *CardStore, result *ValidationResult) {
	if card == nil || store == nil || card.Type != CardTypeFeature || !analysisModePattern.MatchString(card.Body) {
		return
	}
	evidence, ok := markdownSection(card.Body, "Evidence")
	if !ok {
		return
	}
	seen := map[string]bool{}
	for _, findingID := range findingIDPattern.FindAllString(evidence, -1) {
		if seen[findingID] {
			continue
		}
		seen[findingID] = true
		finding, err := store.ReadCard(findingID)
		if err != nil {
			result.AddError("body.analysis.evidence", fmt.Sprintf("referenced finding not found: %s", findingID))
			continue
		}
		if finding.Type != CardTypeFinding {
			result.AddError("body.analysis.evidence", fmt.Sprintf("referenced evidence must be a FIND card: %s", findingID))
		}
	}
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
