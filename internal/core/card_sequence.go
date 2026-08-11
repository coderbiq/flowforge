package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	proposalCardSequenceWidth = 3
	sequenceLockWait          = 30 * time.Second
)

// NextCardID allocates a human-readable ID for a card belonging to a proposal.
// Legacy timestamp IDs remain available for cards without proposal ownership.
func (s *CardStore) NextCardID(cardType CardType, proposalID string) (string, error) {
	if proposalID == "" {
		return GenerateCardID(cardType, ""), nil
	}

	sequence, err := s.nextProposalCardSequence(proposalID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s-%0*d", cardType.Prefix(), proposalID, proposalCardSequenceWidth, sequence), nil
}

// NextTaskID allocates a task ID using the same proposal-wide sequence as all
// other proposal cards. The task kind remains part of the stable ID shape.
func (s *CardStore) NextTaskID(proposalID, taskType string) (string, error) {
	if proposalID == "" {
		return GenerateTaskID("", taskType), nil
	}
	if taskType == "" {
		taskType = "i"
	}

	sequence, err := s.nextProposalCardSequence(proposalID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("TASK-%s-%s-%0*d", proposalID, taskType, proposalCardSequenceWidth, sequence), nil
}

func (s *CardStore) nextProposalCardSequence(proposalID string) (sequence int, err error) {
	cardsDir := s.ProposalCardsDir(proposalID)
	if err := os.MkdirAll(cardsDir, 0755); err != nil {
		return 0, fmt.Errorf("creating proposal cards directory: %w", err)
	}

	release, err := acquireSequenceLock(filepath.Join(cardsDir, ".flowforge-sequence.lock"))
	if err != nil {
		return 0, fmt.Errorf("locking proposal card sequence: %w", err)
	}
	defer func() {
		releaseErr := release()
		if err == nil && releaseErr != nil {
			err = fmt.Errorf("releasing proposal card sequence lock: %w", releaseErr)
		}
	}()

	sequencePath := filepath.Join(cardsDir, ".flowforge-card-sequence")
	next, readErr := readSequenceCounter(sequencePath)
	if readErr != nil {
		return 0, readErr
	}
	if next == 0 {
		cards, listErr := s.ListCardsFromFiles(cardsDir)
		if listErr != nil {
			return 0, fmt.Errorf("listing proposal cards for sequence: %w", listErr)
		}
		maxSequence := 0
		for _, card := range cards {
			sequence, ok := proposalCardSequence(card.ID, proposalID)
			if ok && sequence > maxSequence {
				maxSequence = sequence
			}
		}
		next = maxSequence + 1
	}
	if err := writeSequenceCounter(sequencePath, next+1); err != nil {
		return 0, err
	}
	return next, nil
}

func readSequenceCounter(path string) (int, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("reading card sequence counter: %w", err)
	}
	value := strings.TrimSpace(string(data))
	if value == "" {
		return 0, fmt.Errorf("card sequence counter is empty: %s", path)
	}
	next, err := strconv.Atoi(value)
	if err != nil || next < 1 {
		return 0, fmt.Errorf("invalid card sequence counter %q: %s", value, path)
	}
	return next, nil
}

func writeSequenceCounter(path string, next int) error {
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(strconv.Itoa(next)+"\n"), 0600); err != nil {
		return fmt.Errorf("writing card sequence counter: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		removeErr := os.Remove(tmpPath)
		if removeErr != nil && !os.IsNotExist(removeErr) {
			return fmt.Errorf("installing card sequence counter: %w; removing temporary counter: %v", err, removeErr)
		}
		return fmt.Errorf("installing card sequence counter: %w", err)
	}
	return nil
}

func proposalCardSequence(cardID, proposalID string) (int, bool) {
	parts := strings.Split(cardID, "-")
	if len(parts) < 3 || parts[1] != proposalID {
		return 0, false
	}
	sequencePart := parts[2]
	if parts[0] == "TASK" {
		if len(parts) < 4 {
			return 0, false
		}
		sequencePart = parts[3]
	}
	if len(sequencePart) < proposalCardSequenceWidth {
		return 0, false
	}
	for _, r := range sequencePart {
		if r < '0' || r > '9' {
			return 0, false
		}
	}
	sequence, err := strconv.Atoi(sequencePart)
	if err != nil || sequence < 1 {
		return 0, false
	}
	return sequence, true
}

func acquireSequenceLock(path string) (func() error, error) {
	deadline := time.Now().Add(sequenceLockWait)
	for {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err == nil {
			if closeErr := file.Close(); closeErr != nil {
				removeErr := os.Remove(path)
				if removeErr != nil && !os.IsNotExist(removeErr) {
					return nil, fmt.Errorf("closing sequence lock: %w; removing lock: %v", closeErr, removeErr)
				}
				return nil, fmt.Errorf("closing sequence lock: %w", closeErr)
			}
			return func() error { return os.Remove(path) }, nil
		}
		if !os.IsExist(err) {
			return nil, fmt.Errorf("creating sequence lock: %w", err)
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out waiting for sequence lock %s", path)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
