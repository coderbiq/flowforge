package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CardStore struct {
	wikiRoot string
}

func NewCardStore(wikiRoot string) *CardStore {
	return &CardStore{wikiRoot: wikiRoot}
}

func (s *CardStore) WorkspaceDir() string {
	return filepath.Join(s.wikiRoot, "01-workspace")
}

func (s *CardStore) ActiveDir() string {
	return filepath.Join(s.WorkspaceDir(), "01-active")
}

func (s *CardStore) IntakeDir() string {
	return filepath.Join(s.WorkspaceDir(), "02-intake")
}

func (s *CardStore) CompletedDir() string {
	return filepath.Join(s.WorkspaceDir(), "03-completed")
}

func (s *CardStore) LibraryDir() string {
	return filepath.Join(s.wikiRoot, "02-library")
}

func (s *CardStore) LibraryTypeDir(cardType CardType) string {
	dirName := ""
	switch cardType {
	case CardTypeRequirement:
		dirName = "10-requirements"
	case CardTypeDecision:
		dirName = "20-decisions"
	case CardTypeDesign:
		dirName = "30-designs"
	case CardTypeTask:
		dirName = "40-tasks"
	case CardTypeLog:
		dirName = "50-logs"
	case CardTypeConvention:
		dirName = "60-conventions"
	case CardTypeFinding:
		dirName = "70-findings"
	case CardTypeModule:
		dirName = "80-modules"
	case CardTypeStructure:
		dirName = "structures"
	default:
		dirName = "misc"
	}
	return filepath.Join(s.LibraryDir(), dirName)
}

func (s *CardStore) ProposalDir(proposalID string) string {
	return filepath.Join(s.ActiveDir(), proposalID)
}

func (s *CardStore) ProposalCardsDir(proposalID string) string {
	return filepath.Join(s.ProposalDir(proposalID), "90-cards")
}

func (s *CardStore) ProposalRootCardPath(proposalID string) string {
	return filepath.Join(s.ProposalDir(proposalID), "ROOT-"+proposalID+".md")
}

func (s *CardStore) ProposalRequirementIndexPath(proposalID string) string {
	return filepath.Join(s.ProposalDir(proposalID), "STR-"+proposalID+"-REQ.md")
}

func (s *CardStore) CreateProposal(proposalID, title string) (string, string, error) {
	proposalDir := s.ProposalDir(proposalID)

	if _, err := os.Stat(proposalDir); err == nil {
		return "", "", fmt.Errorf("proposal %s already exists", proposalID)
	}

	dirs := []string{
		proposalDir,
		filepath.Join(proposalDir, "90-cards"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", "", fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	rootPath := s.ProposalRootCardPath(proposalID)
	indexPath := s.ProposalRequirementIndexPath(proposalID)

	rootCard := NewCard(CardTypeStructure, title)
	rootCard.ID = "ROOT-" + proposalID
	rootCard.Status = CardStatusActive
	rootCard.Source = proposalID
	rootCard.Body = fmt.Sprintf("# %s\n\nProposal root card.\n\n## Summary\n\nStable entry for this proposal.\n", title)
	rootCard.AddLink("STR-"+proposalID+"-REQ", "references")
	if err := rootCard.Save(rootPath); err != nil {
		return "", "", fmt.Errorf("writing root card: %w", err)
	}

	indexCard := NewCard(CardTypeStructure, title+" Requirements")
	indexCard.ID = "STR-" + proposalID + "-REQ"
	indexCard.Status = CardStatusActive
	indexCard.Source = proposalID
	indexCard.Body = fmt.Sprintf("# %s Requirements\n\nTop-level requirement index for %s.\n\n## Open Questions\n\n- None\n", title, title)
	if err := indexCard.Save(indexPath); err != nil {
		return "", "", fmt.Errorf("writing requirement index card: %w", err)
	}

	return rootPath, indexPath, nil
}

func (s *CardStore) CreateCard(card *Card, proposalID string) (string, error) {
	var targetDir string

	if proposalID != "" {
		targetDir = s.ProposalCardsDir(proposalID)
		card.Source = proposalID
	} else if card.Type == CardTypeRequirement {
		targetDir = s.IntakeDir()
	} else {
		targetDir = s.LibraryTypeDir(card.Type)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("creating directory: %w", err)
	}

	filename := GenerateFilename(card.ID, card.Title)
	filePath := filepath.Join(targetDir, filename)

	existingCards, _ := s.ListCards(targetDir)
	for i, existing := range existingCards {
		if existing.FilePath == filePath {
			filename = GenerateFilename(card.ID, card.Title)
			filename = strings.TrimSuffix(filename, ".md")
			filename = fmt.Sprintf("%s-%d.md", filename, i+2)
			filePath = filepath.Join(targetDir, filename)
			break
		}
	}

	if err := card.Save(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func (s *CardStore) ReadCard(cardID string) (*Card, error) {
	filePath, err := s.FindCardPath(cardID)
	if err != nil {
		return nil, err
	}
	return ParseCardFile(filePath)
}

func (s *CardStore) UpdateCard(card *Card) error {
	if card.FilePath == "" {
		existingPath, err := s.FindCardPath(card.ID)
		if err != nil {
			return err
		}
		card.FilePath = existingPath
	}

	oldPath := card.FilePath
	newFilename := GenerateFilename(card.ID, card.Title)
	newPath := filepath.Join(filepath.Dir(oldPath), newFilename)

	if err := card.Save(newPath); err != nil {
		return err
	}

	if oldPath != newPath {
		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("removing old file: %w", err)
		}
		card.FilePath = newPath
	}

	return nil
}

func (s *CardStore) DeleteCard(cardID string) error {
	card, err := s.ReadCard(cardID)
	if err != nil {
		return err
	}

	if card.Status != CardStatusDraft {
		return fmt.Errorf("only draft cards can be deleted (current status: %s)", card.Status)
	}

	if err := os.Remove(card.FilePath); err != nil {
		return fmt.Errorf("deleting card: %w", err)
	}

	return nil
}

func (s *CardStore) FindCardPath(cardID string) (string, error) {
	searchDirs := []string{
		s.ActiveDir(),
		s.IntakeDir(),
		s.LibraryDir(),
	}

	for _, dir := range searchDirs {
		path, err := s.findCardInDir(cardID, dir)
		if err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("card not found: %s", cardID)
}

func (s *CardStore) findCardInDir(cardID string, dir string) (string, error) {
	var found string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		filename := filepath.Base(path)
		id, _, parseErr := ParseFilename(filename)
		if parseErr != nil {
			id = strings.TrimSuffix(filename, ".md")
		}

		if id == cardID {
			found = path
			return filepath.SkipAll
		}

		return nil
	})

	if found != "" {
		return found, nil
	}
	if err != nil && err != filepath.SkipAll {
		return "", err
	}
	return "", fmt.Errorf("not found")
}

func (s *CardStore) ListCards(dir string) ([]*Card, error) {
	var cards []*Card

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		filename := filepath.Base(path)
		if strings.HasPrefix(filename, "00-") || strings.HasPrefix(filename, "01-") ||
			strings.HasPrefix(filename, "02-") || strings.HasPrefix(filename, "03-") {
			return nil
		}

		card, err := ParseCardFile(path)
		if err != nil {
			return nil
		}

		cards = append(cards, card)
		return nil
	})

	return cards, err
}

func (s *CardStore) ListCardsByType(cardType CardType) ([]*Card, error) {
	var allCards []*Card

	libraryCards, _ := s.ListCards(s.LibraryTypeDir(cardType))
	allCards = append(allCards, libraryCards...)

	activeCards, _ := s.ListCards(s.ActiveDir())
	for _, card := range activeCards {
		if card.Type == cardType {
			allCards = append(allCards, card)
		}
	}

	intakeCards, _ := s.ListCards(s.IntakeDir())
	for _, card := range intakeCards {
		if card.Type == cardType {
			allCards = append(allCards, card)
		}
	}

	completedCards, _ := s.ListCards(s.CompletedDir())
	for _, card := range completedCards {
		if card.Type == cardType {
			allCards = append(allCards, card)
		}
	}

	return allCards, nil
}

func (s *CardStore) ListCardsByStatus(cardType CardType, status CardStatus) ([]*Card, error) {
	cards, err := s.ListCardsByType(cardType)
	if err != nil {
		return nil, err
	}

	var filtered []*Card
	for _, card := range cards {
		if card.Status == status {
			filtered = append(filtered, card)
		}
	}

	return filtered, nil
}

func (s *CardStore) GetDependents(cardID string) ([]*Card, error) {
	var dependents []*Card

	allDirs := []string{s.ActiveDir(), s.LibraryDir(), s.IntakeDir(), s.CompletedDir()}

	for _, dir := range allDirs {
		cards, _ := s.ListCards(dir)
		for _, card := range cards {
			for _, link := range card.Links {
				if link.Target == cardID {
					dependents = append(dependents, card)
					break
				}
			}
		}
	}

	return dependents, nil
}

func (s *CardStore) GetRelated(cardID string, relation string, depth int) ([]*Card, error) {
	if depth <= 0 {
		depth = 1
	}

	var related []*Card
	visited := map[string]bool{cardID: true}

	currentLevel := []string{cardID}

	for d := 0; d < depth; d++ {
		var nextLevel []string

		for _, id := range currentLevel {
			card, err := s.ReadCard(id)
			if err != nil {
				continue
			}

			for _, link := range card.Links {
				if relation != "" && link.Relation != relation {
					continue
				}
				if visited[link.Target] {
					continue
				}

				targetCard, err := s.ReadCard(link.Target)
				if err != nil {
					continue
				}

				related = append(related, targetCard)
				visited[link.Target] = true
				nextLevel = append(nextLevel, link.Target)
			}
		}

		currentLevel = nextLevel
	}

	return related, nil
}
