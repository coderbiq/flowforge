package core

import (
	"fmt"
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
		} else if !isValidRelation(link.Relation) {
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

func isValidRelation(relation string) bool {
	validRelations := map[string]bool{
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

	filename := strings.TrimSuffix(filePath[strings.LastIndex(filePath, "/")+1:], ".md")
	expectedFilename := strings.TrimSuffix(GenerateFilename(card.ID, card.Title), ".md")
	if filename != expectedFilename {
		result.AddError("filename", fmt.Sprintf("mismatch: expected %s, got %s", expectedFilename, filename))
	}

	return result
}
