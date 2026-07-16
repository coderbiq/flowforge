package command

import (
	"fmt"
	"strings"

	"flowforge/internal/core"
)

type parsedLinkArg struct {
	target   string
	relation string
}

func parseLinkArg(linkStr string) (target string, relation string) {
	parts := strings.SplitN(strings.TrimSpace(linkStr), ":", 2)
	target = strings.TrimSpace(parts[0])
	relation = "related"
	if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
		relation = strings.TrimSpace(parts[1])
	}
	return target, relation
}

func parseLinkArgs(values []string) ([]parsedLinkArg, error) {
	var parsed []parsedLinkArg
	for _, value := range values {
		for _, raw := range strings.Split(value, ",") {
			raw = strings.TrimSpace(raw)
			if raw == "" {
				continue
			}
			target, relation := parseLinkArg(raw)
			if target == "" {
				return nil, fmt.Errorf("invalid link %q: target is required", raw)
			}
			if !core.IsValidRelation(relation) {
				return nil, fmt.Errorf("invalid link %q: relation %q is not supported", raw, relation)
			}
			parsed = append(parsed, parsedLinkArg{target: target, relation: relation})
		}
	}
	return parsed, nil
}
