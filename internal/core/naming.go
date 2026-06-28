package core

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode"
)

var lastCardTimestampNano int64

func GenerateCardTimestamp() string {
	now := time.Now().UnixNano()
	for {
		last := atomic.LoadInt64(&lastCardTimestampNano)
		if now <= last {
			now = last + 1
		}
		if atomic.CompareAndSwapInt64(&lastCardTimestampNano, last, now) {
			return strconv.FormatInt(now, 36)
		}
	}
}

func GenerateCardID(cardType CardType, proposalTs string) string {
	prefix := cardType.Prefix()
	if prefix == "" {
		prefix = string(cardType)
	}
	cardTs := GenerateCardTimestamp()

	if proposalTs == "" {
		return fmt.Sprintf("%s-%s", prefix, cardTs)
	}
	return fmt.Sprintf("%s-%s-%s", prefix, proposalTs, cardTs)
}

func GenerateTaskID(proposalTs string, taskType string) string {
	cardTs := GenerateCardTimestamp()
	if taskType == "" {
		taskType = "i"
	}
	if proposalTs == "" {
		return fmt.Sprintf("TASK-%s-%s", taskType, cardTs)
	}
	return fmt.Sprintf("TASK-%s-%s-%s", proposalTs, taskType, cardTs)
}

func GenerateSubTaskID(parentTaskID string) (string, error) {
	parts := strings.Split(parentTaskID, "-")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid parent task ID format: %s (expected TASK-{proposalTs}-{type}-{taskTs})", parentTaskID)
	}

	existingLetter := 'a' - 1
	lastPart := parts[len(parts)-1]
	if len(lastPart) == 1 && lastPart[0] >= 'a' && lastPart[0] <= 'z' {
		existingLetter = rune(lastPart[0])
		parts = parts[:len(parts)-1]
	}

	nextLetter := string([]rune{existingLetter + 1})
	baseID := strings.Join(parts, "-")
	return baseID + "-" + nextLetter, nil
}

func GenerateProposalID() string {
	return GenerateProposalIDPrefix() + "01"
}

func GenerateProposalIDPrefix() string {
	now := time.Now()
	yy := now.Year() % 100
	mm := int(now.Month())
	dd := now.Day()
	return fmt.Sprintf("CR%02d%02d%02d", yy, mm, dd)
}

func GenerateFilename(id string, title string) string {
	slug := ToSlug(title)
	if slug == "" {
		slug = "untitled"
	}
	return fmt.Sprintf("%s_%s.md", id, slug)
}

func ToSlug(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	var result strings.Builder
	prevDash := false
	var prevR rune
	runes := []rune(s)

	for i, r := range runes {
		if (unicode.IsLetter(r) || unicode.IsDigit(r)) && isVisible(r) {
			hasLowerNext := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if unicode.IsUpper(r) && result.Len() > 0 && !prevDash && unicode.IsLower(prevR) && hasLowerNext {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
			prevDash = false
			prevR = r
		} else if r == ' ' || r == '_' || r == '-' {
			if !prevDash && result.Len() > 0 {
				result.WriteRune('-')
				prevDash = true
			}
			prevR = 0
		} else {
			prevR = 0
		}
	}

	slug := result.String()
	slug = strings.Trim(slug, "-")

	runes = []rune(slug)
	if len(runes) > 50 {
		runes = runes[:50]
		lastDash := -1
		for i := len(runes) - 1; i >= 0; i-- {
			if runes[i] == '-' {
				lastDash = i
				break
			}
		}
		if lastDash > 30 {
			runes = runes[:lastDash]
		}
		slug = string(runes)
	}

	return slug
}

func isVisible(r rune) bool {
	if unicode.Is(unicode.Cf, r) {
		return false
	}
	switch r {
	case '\u115F', '\u1160', '\u3164':
		return false
	}
	return true
}

func ParseFilename(filename string) (id string, slug string, err error) {
	filename = strings.TrimSuffix(filename, ".md")

	parts := strings.SplitN(filename, "_", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid filename format: %s (expected {ID}_{slug}.md)", filename)
	}

	return parts[0], parts[1], nil
}

func ParseCardID(cardID string) (cardType CardType, proposalTs string, cardTs string, err error) {
	parts := strings.Split(cardID, "-")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid card ID format: %s", cardID)
	}

	prefix := parts[0]
	cardType = CardTypeFromPrefix(prefix)
	if cardType == "" {
		return "", "", "", fmt.Errorf("unknown card type prefix: %s", prefix)
	}

	if cardType == CardTypeTask {
		if len(parts) >= 4 {
			proposalTs = parts[1]
			cardTs = parts[3]
			if len(parts) > 4 {
				cardTs = parts[3] + "-" + parts[4]
			}
		} else if len(parts) == 3 {
			cardTs = parts[2]
		}
	} else {
		if len(parts) >= 3 {
			proposalTs = parts[1]
			cardTs = parts[2]
		} else {
			cardTs = parts[1]
		}
	}

	return cardType, proposalTs, cardTs, nil
}

func IsSubTaskID(cardID string) bool {
	parts := strings.Split(cardID, "-")
	if len(parts) < 2 {
		return false
	}
	lastPart := parts[len(parts)-1]
	return len(lastPart) == 1 && lastPart[0] >= 'a' && lastPart[0] <= 'z'
}

func GetParentTaskID(subTaskID string) (string, error) {
	if !IsSubTaskID(subTaskID) {
		return "", fmt.Errorf("not a sub-task ID: %s", subTaskID)
	}
	parts := strings.Split(subTaskID, "-")
	return strings.Join(parts[:len(parts)-1], "-"), nil
}
