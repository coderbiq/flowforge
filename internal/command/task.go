package command

import "strings"

func proposalTimestamp(proposalID string) string {
	if proposalID == "" {
		return ""
	}
	parts := strings.Split(proposalID, "-")
	if len(parts) >= 2 {
		return parts[1]
	}
	return proposalID
}
