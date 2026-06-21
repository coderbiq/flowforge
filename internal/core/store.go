package core

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CardSyncService interface {
	ReadCard(id string) (*Card, error)
	FindCardPath(id string) (string, error)
	ListCards(dir string) ([]*Card, error)
	ListCardsByType(cardType CardType) ([]*Card, error)
	GetDependents(cardID string) ([]*Card, error)
	SyncCard(card *Card) error
	DeleteCard(cardID string) error
}

type CardStore struct {
	wikiRoot    string
	syncService CardSyncService
}

func NewCardStore(wikiRoot string) *CardStore {
	return &CardStore{wikiRoot: wikiRoot}
}

func NewCardStoreWithSync(wikiRoot string, syncService CardSyncService) *CardStore {
	return &CardStore{wikiRoot: wikiRoot, syncService: syncService}
}

func (s *CardStore) hasSync() bool {
	return s.syncService != nil
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
	case CardTypeProposal:
		dirName = "structures"
	default:
		dirName = "misc"
	}
	return filepath.Join(s.LibraryDir(), dirName)
}

func (s *CardStore) ProposalCardDir() string {
	return filepath.Join(s.wikiRoot, "03-proposal")
}

func (s *CardStore) ProposalDir(proposalID string) string {
	dir := s.findProposalDirIn(proposalID, s.ActiveDir())
	if dir == "" {
		return filepath.Join(s.ActiveDir(), proposalID)
	}
	return dir
}

func (s *CardStore) FindProposalDir(proposalID string) string {
	for _, baseDir := range []string{s.ActiveDir(), s.CompletedDir()} {
		if dir := s.findProposalDirIn(proposalID, baseDir); dir != "" {
			return dir
		}
	}
	return filepath.Join(s.ActiveDir(), proposalID)
}

func (s *CardStore) findProposalDirIn(proposalID string, baseDir string) string {
	dir := filepath.Join(baseDir, proposalID)
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		return dir
	}
	entries, _ := os.ReadDir(baseDir)
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), proposalID+"_") {
			return filepath.Join(baseDir, entry.Name())
		}
	}
	return ""
}

func (s *CardStore) ProposalCardsDir(proposalID string) string {
	return filepath.Join(s.ProposalDir(proposalID), "90-cards")
}

func (s *CardStore) ProposalRootCardPath(proposalID string) string {
	// Proposal root card 现在放在 03-proposal/ 下
	// 需要先找到 dir_name 来构造文件名
	metaCard, err := s.findProposalMetaCard(proposalID)
	if err == nil && metaCard.DirName != "" {
		return filepath.Join(s.ProposalCardDir(), metaCard.DirName+".md")
	}
	// 向后兼容：workspace 下的旧位置
	return filepath.Join(s.ProposalDir(proposalID), "ROOT-"+proposalID+".md")
}

func (s *CardStore) ProposalRequirementIndexPath(proposalID string) string {
	return filepath.Join(s.ProposalDir(proposalID), "STR-"+proposalID+"-REQ.md")
}

// findProposalMetaCard 在 03-proposal/ 或 workspace 下查找 proposal 卡片
func (s *CardStore) findProposalMetaCard(proposalID string) (*Card, error) {
	// 先在 03-proposal/ 下前缀匹配
	entries, err := os.ReadDir(s.ProposalCardDir())
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			if strings.HasPrefix(entry.Name(), proposalID+"_") {
				path := filepath.Join(s.ProposalCardDir(), entry.Name())
				return ParseCardFile(path)
			}
		}
	}
	// 向后兼容：workspace 下的旧 ROOT-{id}.md
	oldPath := filepath.Join(s.ProposalDir(proposalID), "ROOT-"+proposalID+".md")
	if card, err := ParseCardFile(oldPath); err == nil {
		return card, nil
	}
	return nil, fmt.Errorf("proposal meta card not found for %s", proposalID)
}

func (s *CardStore) CreateProposal(proposalID, title string) (string, string, error) {
	slug := ToSlug(title)
	dirName := proposalID + "_" + slug
	proposalDir := filepath.Join(s.ActiveDir(), dirName)

	if _, err := os.Stat(proposalDir); err == nil {
		return "", "", fmt.Errorf("proposal %s already exists", proposalID)
	}

	dirs := []string{
		proposalDir,
		filepath.Join(proposalDir, "90-cards"),
		s.ProposalCardDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", "", fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	rootCard := NewCard(CardTypeProposal, title)
	rootCard.ID = "PROP-" + proposalID
	rootCard.Status = CardStatusActive
	rootCard.Source = proposalID
	rootCard.ProposalID = proposalID
	rootCard.DirName = dirName
	rootCard.Slug = slug
	rootCard.Body = fmt.Sprintf("# %s\n\n## Purpose\n\nStable entry for proposal %s.\n\n## Entries\n\n- [STR-%s-REQ](STR-%s-REQ.md) (structure, active) - Requirement index\n\n## Summary\n\nProposal root card.\n", title, proposalID, proposalID, proposalID)
	rootCard.AddLink("STR-"+proposalID+"-REQ", "indexes")

	rootPath := filepath.Join(s.ProposalCardDir(), dirName+".md")
	if err := rootCard.Save(rootPath); err != nil {
		return "", "", fmt.Errorf("writing proposal card: %w", err)
	}

	indexPath := s.ProposalRequirementIndexPath(proposalID)
	indexCard := NewCard(CardTypeStructure, title+" Requirements")
	indexCard.ID = "STR-" + proposalID + "-REQ"
	indexCard.Status = CardStatusActive
	indexCard.Source = proposalID
	indexCard.AddLink("PROP-"+proposalID, "belongs_to")
	indexCard.Body = fmt.Sprintf("# %s Requirements\n\n## Purpose\n\nTop-level requirement index for %s.\n\n## Entries\n\n- None\n\n## Open Questions\n\n- None\n", title, title)
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

	if s.hasSync() {
		existingCards, _ := s.syncService.ListCards(targetDir)
		for i, existing := range existingCards {
			if existing.FilePath == filePath {
				filename = GenerateFilename(card.ID, card.Title)
				filename = strings.TrimSuffix(filename, ".md")
				filename = fmt.Sprintf("%s-%d.md", filename, i+2)
				filePath = filepath.Join(targetDir, filename)
				break
			}
		}
	} else {
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
	}

	if err := card.Save(filePath); err != nil {
		return "", err
	}

	if s.hasSync() {
		card.FilePath = filePath
		if err := s.syncService.SyncCard(card); err != nil {
			return filePath, fmt.Errorf("syncing card: %w", err)
		}
	}

	return filePath, nil
}

func (s *CardStore) ReadCard(cardID string) (*Card, error) {
	if s.hasSync() {
		card, err := s.syncService.ReadCard(cardID)
		if err == nil {
			return card, nil
		}
	}

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

	if err := card.Save(card.FilePath); err != nil {
		return err
	}

	if s.hasSync() {
		if err := s.syncService.SyncCard(card); err != nil {
			return fmt.Errorf("syncing card: %w", err)
		}
	}

	return nil
}

func (s *CardStore) UpdateCardWithLock(cardID string, mutate func(*Card) error) (err error) {
	lockPath := s.cardLockPath(cardID)
	release, err := acquireCardLock(lockPath)
	if err != nil {
		return err
	}
	defer func() {
		releaseErr := release()
		if err == nil && releaseErr != nil {
			err = releaseErr
		}
	}()

	card, err := s.ReadCard(cardID)
	if err != nil {
		return fmt.Errorf("reading card %s: %w", cardID, err)
	}

	if err := mutate(card); err != nil {
		return err
	}

	if err := s.UpdateCard(card); err != nil {
		return fmt.Errorf("updating card %s: %w", cardID, err)
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

	if s.hasSync() {
		if err := s.syncService.DeleteCard(cardID); err != nil {
			return fmt.Errorf("syncing card deletion: %w", err)
		}
	}

	return nil
}

// ForceDeleteCard removes a card regardless of its draft/active status
// and cleans up all backlinks from other cards that point to it.
func (s *CardStore) ForceDeleteCard(cardID string) error {
	card, err := s.ReadCard(cardID)
	if err != nil {
		return err
	}

	// Remove backlinks from all cards that reference the deleted card.
	dependents, err := s.GetDependents(cardID)
	if err != nil {
		return fmt.Errorf("finding dependents of %s: %w", cardID, err)
	}
	for _, dep := range dependents {
		found := dep.RemoveLink(cardID, "")
		if !found {
			// Try removing with common relations.
			for _, rel := range []string{"references", "requires", "implements", "satisfies", "records", "indexes", "belongs_to", "related"} {
				if dep.RemoveLink(cardID, rel) {
					found = true
					break
				}
			}
		}
		if found {
			if err := s.UpdateCardWithLock(dep.ID, func(uc *Card) error {
				uc.Links = dep.Links
				return nil
			}); err != nil {
				return fmt.Errorf("removing backlink from %s: %w", dep.ID, err)
			}
		}
	}

	if err := os.Remove(card.FilePath); err != nil {
		return fmt.Errorf("deleting card: %w", err)
	}

	if s.hasSync() {
		if err := s.syncService.DeleteCard(cardID); err != nil {
			return fmt.Errorf("syncing card deletion: %w", err)
		}
	}

	return nil
}

func (s *CardStore) FindCardPath(cardID string) (string, error) {
	if s.hasSync() {
		path, err := s.syncService.FindCardPath(cardID)
		if err == nil {
			return path, nil
		}
	}

	searchDirs := []string{
		s.ActiveDir(),
		s.IntakeDir(),
		s.LibraryDir(),
		s.ProposalCardDir(),
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

		if strings.HasPrefix(cardID, "PROP-") {
			proposalID := strings.TrimPrefix(cardID, "PROP-")
			if strings.HasPrefix(filename, proposalID+"_") || strings.HasPrefix(id, proposalID) || id == "ROOT-"+proposalID {
				found = path
				return filepath.SkipAll
			}
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
	if s.hasSync() {
		cards, err := s.syncService.ListCards(dir)
		if err == nil && len(cards) > 0 {
			return cards, nil
		}
	}

	return s.ListCardsFromFiles(dir)
}

func (s *CardStore) ListCardsFromFiles(dir string) ([]*Card, error) {
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
	if s.hasSync() {
		cards, err := s.syncService.ListCardsByType(cardType)
		if err == nil && len(cards) > 0 {
			return cards, nil
		}
	}

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
	if s.hasSync() {
		cards, err := s.syncService.GetDependents(cardID)
		if err == nil && len(cards) > 0 {
			return cards, nil
		}
	}

	var dependents []*Card

	allDirs := []string{s.ActiveDir(), s.LibraryDir(), s.IntakeDir(), s.CompletedDir(), s.ProposalCardDir()}

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

func (s *CardStore) cardLockPath(cardID string) string {
	sum := sha1.Sum([]byte(s.wikiRoot + "\x00" + cardID))
	lockDir := filepath.Join(os.TempDir(), "flowforge-card-locks")
	return filepath.Join(lockDir, hex.EncodeToString(sum[:])+".lock")
}

func acquireCardLock(lockPath string) (func() error, error) {
	if err := os.MkdirAll(filepath.Dir(lockPath), 0755); err != nil {
		return nil, fmt.Errorf("creating lock directory: %w", err)
	}

	ownerToken, err := newLockOwnerToken()
	if err != nil {
		return nil, err
	}

	deadline := time.Now().Add(10 * time.Second)
	for {
		file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err == nil {
			if _, writeErr := file.WriteString(ownerToken); writeErr != nil {
				closeErr := file.Close()
				_ = os.Remove(lockPath)
				if closeErr != nil {
					return nil, fmt.Errorf("writing lock owner: %w (closing lock file: %v)", writeErr, closeErr)
				}
				return nil, fmt.Errorf("writing lock owner: %w", writeErr)
			}
			if closeErr := file.Close(); closeErr != nil {
				_ = os.Remove(lockPath)
				return nil, fmt.Errorf("closing lock file: %w", closeErr)
			}
			return func() error {
				data, readErr := os.ReadFile(lockPath)
				if os.IsNotExist(readErr) {
					return nil
				}
				if readErr != nil {
					return fmt.Errorf("reading lock owner: %w", readErr)
				}
				if string(data) != ownerToken {
					return nil
				}
				if removeErr := os.Remove(lockPath); removeErr != nil && !os.IsNotExist(removeErr) {
					return fmt.Errorf("removing lock file: %w", removeErr)
				}
				return nil
			}, nil
		}
		if !os.IsExist(err) {
			return nil, fmt.Errorf("creating lock file: %w", err)
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out waiting for card lock %s", lockPath)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func newLockOwnerToken() (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("creating lock owner token: %w", err)
	}
	return fmt.Sprintf("%d:%s", os.Getpid(), hex.EncodeToString(randomBytes)), nil
}
