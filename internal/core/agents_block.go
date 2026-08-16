package core

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	agentsBlockStart = "<!-- FLOWFORGE:START -->"
	agentsBlockEnd   = "<!-- FLOWFORGE:END -->"
)

func StripBlockMarkers(content []byte) []byte {
	startIdx := bytes.Index(content, []byte(agentsBlockStart))
	endIdx := bytes.Index(content, []byte(agentsBlockEnd))

	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return content
	}

	start := startIdx + len(agentsBlockStart)
	if start < len(content) && content[start] == '\n' {
		start++
	}
	end := endIdx
	if end > 0 && content[end-1] == '\n' {
		end--
	}
	return content[start:end]
}

func ApplyAgentsBlock(targetPath string, newContent []byte) error {
	_, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		return createWithBlock(targetPath, newContent)
	}
	if err != nil {
		return fmt.Errorf("checking agents.md: %w", err)
	}

	existing, err := os.ReadFile(targetPath)
	if err != nil {
		return fmt.Errorf("reading agents.md: %w", err)
	}

	startIdx := bytes.Index(existing, []byte(agentsBlockStart))
	endIdx := bytes.Index(existing, []byte(agentsBlockEnd))

	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return appendBlock(targetPath, existing, newContent)
	}

	return replaceBlock(targetPath, existing, startIdx, endIdx, newContent)
}

func ApplyMarkedBlock(targetPath, startMarker, endMarker string, newContent []byte) error {
	existing, err := os.ReadFile(targetPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}
		existing = nil
	} else if err != nil {
		return fmt.Errorf("reading marked block target: %w", err)
	}
	start := bytes.Index(existing, []byte(startMarker))
	end := bytes.Index(existing, []byte(endMarker))
	var buf bytes.Buffer
	if start >= 0 && end > start {
		buf.Write(existing[:start])
	} else {
		buf.Write(existing)
		if len(existing) > 0 && existing[len(existing)-1] != '\n' {
			buf.WriteByte('\n')
		}
		if len(existing) > 0 {
			buf.WriteByte('\n')
		}
	}
	buf.WriteString(startMarker + "\n")
	buf.Write(newContent)
	if len(newContent) > 0 && newContent[len(newContent)-1] != '\n' {
		buf.WriteByte('\n')
	}
	buf.WriteString(endMarker + "\n")
	if start >= 0 && end > start {
		after := end + len(endMarker)
		if after < len(existing) && existing[after] == '\n' {
			after++
		}
		buf.Write(existing[after:])
	}
	return os.WriteFile(targetPath, buf.Bytes(), 0644)
}

// ExtractMarkedBlock returns the exact content between a pair of line markers.
func ExtractMarkedBlock(content []byte, startMarker, endMarker string) ([]byte, bool, error) {
	start := bytes.Index(content, []byte(startMarker))
	end := bytes.Index(content, []byte(endMarker))
	if start < 0 && end < 0 {
		return nil, false, nil
	}
	if start < 0 || end <= start {
		return nil, false, fmt.Errorf("invalid managed block markers")
	}
	start += len(startMarker)
	if start < len(content) && content[start] == '\n' {
		start++
	}
	blockEnd := end
	if blockEnd > start && content[blockEnd-1] == '\n' {
		blockEnd--
	}
	result := append([]byte(nil), content[start:blockEnd]...)
	if len(result) > 0 {
		result = append(result, '\n')
	}
	return result, true, nil
}

func RemoveMarkedBlock(targetPath, startMarker, endMarker string) error {
	existing, err := os.ReadFile(targetPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading marked block target: %w", err)
	}
	start := bytes.Index(existing, []byte(startMarker))
	end := bytes.Index(existing, []byte(endMarker))
	if start < 0 || end <= start {
		return nil
	}
	after := end + len(endMarker)
	if after < len(existing) && existing[after] == '\n' {
		after++
	}
	result := append([]byte{}, existing[:start]...)
	result = append(result, existing[after:]...)
	return os.WriteFile(targetPath, result, 0644)
}

// RemoveMarkedBlockContent removes only the requested marker pair and keeps
// every byte outside it. It is useful to validate a change before writing it.
func RemoveMarkedBlockContent(existing []byte, startMarker, endMarker string) ([]byte, bool, error) {
	start := bytes.Index(existing, []byte(startMarker))
	end := bytes.Index(existing, []byte(endMarker))
	if start < 0 && end < 0 {
		return append([]byte(nil), existing...), false, nil
	}
	if start < 0 || end <= start {
		return nil, false, fmt.Errorf("invalid marked block markers")
	}
	after := end + len(endMarker)
	if after < len(existing) && existing[after] == '\n' {
		after++
	}
	result := append([]byte{}, existing[:start]...)
	result = append(result, existing[after:]...)
	return result, true, nil
}

func createWithBlock(path string, content []byte) error {
	var buf bytes.Buffer
	buf.WriteString(agentsBlockStart)
	buf.WriteString("\n")
	buf.Write(content)
	if len(content) > 0 && content[len(content)-1] != '\n' {
		buf.WriteString("\n")
	}
	buf.WriteString(agentsBlockEnd)
	buf.WriteString("\n")

	if err := os.MkdirAll(".", 0755); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0644)
}

func appendBlock(path string, existing, newContent []byte) error {
	var buf bytes.Buffer
	buf.Write(existing)
	if len(existing) > 0 && existing[len(existing)-1] != '\n' {
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
	buf.WriteString(agentsBlockStart)
	buf.WriteString("\n")
	buf.Write(newContent)
	if len(newContent) > 0 && newContent[len(newContent)-1] != '\n' {
		buf.WriteString("\n")
	}
	buf.WriteString(agentsBlockEnd)
	buf.WriteString("\n")

	return os.WriteFile(path, buf.Bytes(), 0644)
}

func replaceBlock(path string, existing []byte, startIdx, endIdx int, newContent []byte) error {
	var buf bytes.Buffer
	buf.Write(existing[:startIdx])
	buf.WriteString(agentsBlockStart)
	buf.WriteString("\n")
	buf.Write(newContent)
	if len(newContent) > 0 && newContent[len(newContent)-1] != '\n' {
		buf.WriteString("\n")
	}
	buf.WriteString(agentsBlockEnd)
	endOfEndLine := endIdx + len(agentsBlockEnd)
	if endOfEndLine < len(existing) && existing[endOfEndLine] == '\n' {
		endOfEndLine++
	}
	buf.Write(existing[endOfEndLine:])

	return os.WriteFile(path, buf.Bytes(), 0644)
}

func RemoveAgentsBlock(targetPath string) error {
	existing, err := os.ReadFile(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading agents.md: %w", err)
	}

	startIdx := bytes.Index(existing, []byte(agentsBlockStart))
	if startIdx == -1 {
		return nil
	}

	endIdx := bytes.Index(existing, []byte(agentsBlockEnd))
	if endIdx == -1 || endIdx <= startIdx {
		return nil
	}

	var buf bytes.Buffer
	beforeStart := existing[:startIdx]
	if len(bytes.TrimRight(beforeStart, "\n")) > 0 {
		buf.Write(beforeStart)
	}
	afterEnd := endIdx + len(agentsBlockEnd)
	if afterEnd < len(existing) && existing[afterEnd] == '\n' {
		afterEnd++
	}
	buf.Write(existing[afterEnd:])

	return os.WriteFile(targetPath, buf.Bytes(), 0644)
}

func HashBlockContent(content []byte) string {
	lines := bufio.NewScanner(bytes.NewReader(content))
	var buf strings.Builder
	inBlock := false

	for lines.Scan() {
		line := lines.Text()
		if line == agentsBlockStart {
			inBlock = true
			continue
		}
		if line == agentsBlockEnd {
			inBlock = false
			continue
		}
		if inBlock {
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}

	return sha256Hex([]byte(buf.String()))
}
