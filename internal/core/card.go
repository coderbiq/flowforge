package core

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type CardType string

const (
	CardTypeRequirement CardType = "requirement"
	CardTypeDecision    CardType = "decision"
	CardTypeDesign      CardType = "design"
	CardTypeTask        CardType = "task"
	CardTypeLog         CardType = "log"
	CardTypeConvention  CardType = "convention"
	CardTypeFinding     CardType = "finding"
	CardTypeModule      CardType = "module"
	CardTypeStructure   CardType = "structure"
	CardTypeProposal    CardType = "proposal"
	CardTypeFeature     CardType = "feature"
)

func (t CardType) Valid() bool {
	switch t {
	case CardTypeRequirement, CardTypeDecision, CardTypeDesign, CardTypeTask,
		CardTypeLog, CardTypeConvention, CardTypeFinding, CardTypeModule, CardTypeStructure, CardTypeProposal,
		CardTypeFeature:
		return true
	}
	return false
}

func (t CardType) Prefix() string {
	switch t {
	case CardTypeRequirement:
		return "REQ"
	case CardTypeDecision:
		return "DEC"
	case CardTypeDesign:
		return "DES"
	case CardTypeTask:
		return "TASK"
	case CardTypeLog:
		return "LOG"
	case CardTypeConvention:
		return "CONV"
	case CardTypeFinding:
		return "FIND"
	case CardTypeModule:
		return "MOD"
	case CardTypeStructure:
		return "STR"
	case CardTypeProposal:
		return "PROP"
	case CardTypeFeature:
		return "FEAT"
	}
	return ""
}

func CardTypeFromPrefix(prefix string) CardType {
	switch prefix {
	case "REQ":
		return CardTypeRequirement
	case "DEC":
		return CardTypeDecision
	case "DES":
		return CardTypeDesign
	case "TASK":
		return CardTypeTask
	case "LOG":
		return CardTypeLog
	case "CONV":
		return CardTypeConvention
	case "FIND":
		return CardTypeFinding
	case "MOD":
		return CardTypeModule
	case "STR":
		return CardTypeStructure
	case "PROP":
		return CardTypeProposal
	case "FEAT":
		return CardTypeFeature
	}
	return ""
}

type CardStatus string

const (
	CardStatusDraft      CardStatus = "draft"
	CardStatusActive     CardStatus = "active"
	CardStatusAccepted   CardStatus = "accepted"
	CardStatusDeprecated CardStatus = "deprecated"
	CardStatusSuperseded CardStatus = "superseded"
	CardStatusBacklog    CardStatus = "backlog"
	CardStatusNotReady   CardStatus = "not_ready"
	CardStatusReady      CardStatus = "ready"
	CardStatusInProgress CardStatus = "in_progress"
	CardStatusDone       CardStatus = "done"
	CardStatusBlocked    CardStatus = "blocked"
	CardStatusCancelled  CardStatus = "cancelled"
	CardStatusDesigned   CardStatus = "designed"
	CardStatusPlanned    CardStatus = "planned"
	CardStatusCompleted  CardStatus = "completed"
)

func (s CardStatus) Valid() bool {
	switch s {
	case CardStatusDraft, CardStatusActive, CardStatusAccepted, CardStatusDeprecated, CardStatusSuperseded:
		return true
	case CardStatusBacklog, CardStatusNotReady, CardStatusReady, CardStatusInProgress, CardStatusDone, CardStatusBlocked, CardStatusCancelled:
		return true
	case CardStatusDesigned, CardStatusPlanned, CardStatusCompleted:
		return true
	}
	return false
}

type Importance string

const (
	ImportanceMust   Importance = "must"
	ImportanceShould Importance = "should"
	ImportanceMay    Importance = "may"
)

func (i Importance) Valid() bool {
	switch i {
	case ImportanceMust, ImportanceShould, ImportanceMay, "":
		return true
	}
	return false
}

type Link struct {
	Target   string `yaml:"target" json:"target"`
	Relation string `yaml:"relation" json:"relation"`
}

type CardSearchResult struct {
	Card        *Card
	MatchReason string
}

type LibrarySuggestion struct {
	Card              *Card
	Score             int
	MatchReason       string
	SuggestedRelation string
}

type Card struct {
	ID          string      `yaml:"id" json:"id"`
	Title       string      `yaml:"title" json:"title"`
	Type        CardType    `yaml:"type" json:"type"`
	Status      CardStatus  `yaml:"status" json:"status"`
	Importance  Importance  `yaml:"importance,omitempty" json:"importance,omitempty"`
	Tags        []string    `yaml:"tags,omitempty" json:"tags,omitempty"`
	Links       []Link      `yaml:"links,omitempty" json:"links,omitempty"`
	Created     time.Time   `yaml:"created" json:"created"`
	Updated     time.Time   `yaml:"updated" json:"updated"`
	Source      string      `yaml:"source,omitempty" json:"source,omitempty"`
	Domain      string      `yaml:"domain,omitempty" json:"domain,omitempty"`
	ProposalID  string      `yaml:"proposal_id,omitempty" json:"proposalId,omitempty"`
	DirName     string      `yaml:"dir_name,omitempty" json:"dirName,omitempty"`
	Slug        string      `yaml:"slug,omitempty" json:"slug,omitempty"`
	Project     string      `yaml:"project,omitempty" json:"project,omitempty"`
	Role        string      `yaml:"role,omitempty" json:"role,omitempty"`
	Body        string      `yaml:"-" json:"body"`
	FilePath    string      `yaml:"-" json:"filePath,omitempty"`
}

func NewCard(cardType CardType, title string) *Card {
	now := time.Now()
	return &Card{
		Type:       cardType,
		Title:      title,
		Status:     CardStatusDraft,
		Importance: ImportanceShould,
		Created:    now,
		Updated:    now,
	}
}

func ParseCardFile(filePath string) (*Card, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading card file: %w", err)
	}
	return ParseCard(data, filePath)
}

func ParseCard(data []byte, filePath string) (*Card, error) {
	content := string(data)

	frontmatterStart := strings.Index(content, "---\n")
	if frontmatterStart != 0 {
		return nil, fmt.Errorf("card must start with ---")
	}

	frontmatterEnd := strings.Index(content[3:], "\n---")
	if frontmatterEnd < 0 {
		return nil, fmt.Errorf("missing closing ---")
	}
	frontmatterEnd += 3

	frontmatterStr := content[4:frontmatterEnd]
	bodyStart := frontmatterEnd + 4
	if bodyStart < len(content) && content[bodyStart-1] == '\n' {
		bodyStart--
	}
	body := ""
	if bodyStart < len(content) {
		body = strings.TrimSpace(content[bodyStart:])
	}

	var card Card
	if err := yaml.Unmarshal([]byte(frontmatterStr), &card); err != nil {
		return nil, fmt.Errorf("parsing frontmatter: %w", err)
	}

	card.Body = body
	card.FilePath = filePath
	return &card, nil
}

func (c *Card) ToMarkdown() ([]byte, error) {
	frontmatter := struct {
		ID         string     `yaml:"id"`
		Title      string     `yaml:"title"`
		Type       CardType   `yaml:"type"`
		Status     CardStatus `yaml:"status"`
		Importance Importance `yaml:"importance,omitempty"`
		Tags       []string   `yaml:"tags,omitempty"`
		Links      []Link     `yaml:"links,omitempty"`
		Created    time.Time  `yaml:"created"`
		Updated    time.Time  `yaml:"updated"`
		Source     string     `yaml:"source,omitempty"`
		Domain     string     `yaml:"domain,omitempty"`
		ProposalID string     `yaml:"proposal_id,omitempty"`
		DirName    string     `yaml:"dir_name,omitempty"`
		Slug       string     `yaml:"slug,omitempty"`
		Project    string     `yaml:"project,omitempty"`
		Role       string     `yaml:"role,omitempty"`
	}{
		ID:         c.ID,
		Title:      c.Title,
		Type:       c.Type,
		Status:     c.Status,
		Importance: c.Importance,
		Tags:       c.Tags,
		Links:      c.Links,
		Created:    c.Created,
		Updated:    c.Updated,
		Source:     c.Source,
		Domain:     c.Domain,
		ProposalID: c.ProposalID,
		DirName:    c.DirName,
		Slug:       c.Slug,
		Project:    c.Project,
		Role:       c.Role,
	}

	yamlData, err := yaml.Marshal(frontmatter)
	if err != nil {
		return nil, fmt.Errorf("marshaling frontmatter: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(yamlData)
	buf.WriteString("---\n\n")
	if c.Body != "" {
		buf.WriteString(c.Body)
		buf.WriteString("\n")
	}

	return buf.Bytes(), nil
}

func (c *Card) Save(filePath string) error {
	data, err := c.ToMarkdown()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("writing card file: %w", err)
	}
	c.FilePath = filePath
	return nil
}

func (c *Card) AddLink(target string, relation string) {
	for _, link := range c.Links {
		if link.Target == target && link.Relation == relation {
			return
		}
	}
	c.Links = append(c.Links, Link{Target: target, Relation: relation})
	c.Updated = time.Now()
}

func (c *Card) RemoveLink(target string, relation string) bool {
	for i, link := range c.Links {
		if link.Target == target && link.Relation == relation {
			c.Links = append(c.Links[:i], c.Links[i+1:]...)
			c.Updated = time.Now()
			return true
		}
	}
	return false
}
