package core

import (
	"fmt"
	"regexp"
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
	prefix := currentCardPrefix(cardType)
	cardTs := GenerateCardTimestamp()

	if proposalTs == "" {
		return fmt.Sprintf("%s-%s", prefix, cardTs)
	}
	return fmt.Sprintf("%s-%s-%s", prefix, proposalTs, cardTs)
}

func currentCardPrefix(cardType CardType) string {
	switch cardType {
	case CardTypeDecision:
		return "DEC"
	case CardTypeConvention:
		return "CONV"
	case CardTypeFinding:
		return "FIND"
	case CardTypeModule:
		return "MOD"
	case CardTypeProposal:
		return "PROP"
	case CardTypeFeature:
		return "FEAT"
	default:
		return string(cardType)
	}
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
	if !strings.HasSuffix(filename, ".md") {
		return "", "", fmt.Errorf("invalid filename format: %s (expected {ID}_{slug}.md)", filename)
	}
	filename = strings.TrimSuffix(filename, ".md")

	parts := strings.SplitN(filename, "_", 2)
	if len(parts) != 2 || !currentV3CardIDPattern.MatchString(parts[0]) || parts[1] == "" {
		return "", "", fmt.Errorf("invalid filename format: %s (expected {ID}_{slug}.md)", filename)
	}

	return parts[0], parts[1], nil
}

var currentV3CardIDPattern = regexp.MustCompile(`^(?:DEC|CONV|FIND|MOD|PROP|FEAT)-[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`)
