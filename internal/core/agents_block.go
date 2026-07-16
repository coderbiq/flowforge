package core

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
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
