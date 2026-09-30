package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"flowforge/internal/subagent"
)

const agentModelStateName = "agent-model-state.json"

type modelFieldValue struct {
	Present bool   `json:"present"`
	Value   string `json:"value"`
}
type modelFieldState struct {
	LastGenerated *modelFieldValue `json:"last_generated"`
	RetainedLocal *modelFieldValue `json:"retained_local,omitempty"`
}
type agentModelRecord struct {
	Model  *modelFieldState `json:"model"`
	Effort *modelFieldState `json:"effort"`
}
type agentModelState struct {
	Version int                          `json:"version"`
	Paths   map[string]*agentModelRecord `json:"paths"`
}
type nativeModelFields struct{ Model, Effort modelFieldValue }

func (v *modelFieldValue) UnmarshalJSON(data []byte) error {
	var wire struct {
		Present *bool   `json:"present"`
		Value   *string `json:"value"`
	}
	if err := strictStateJSON(data, &wire); err != nil {
		return err
	}
	if wire.Present == nil || wire.Value == nil {
		return fmt.Errorf("field requires present and value")
	}
	if !*wire.Present && *wire.Value != "" {
		return fmt.Errorf("absent field must have empty value")
	}
	v.Present, v.Value = *wire.Present, *wire.Value
	return nil
}
func strictStateJSON(data []byte, target any) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(target); err != nil {
		return err
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("unexpected trailing JSON")
	}
	return nil
}
func agentModelStatePath(root string) string {
	return filepath.Join(root, ".flowforge", agentModelStateName)
}
func readAgentModelState(root string) (*agentModelState, error) {
	path := agentModelStatePath(root)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &agentModelState{Version: 1, Paths: map[string]*agentModelRecord{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	var state agentModelState
	if err := strictStateJSON(data, &state); err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	if state.Version != 1 || state.Paths == nil {
		return nil, fmt.Errorf("reading %s: unsupported version or missing paths", path)
	}
	for _, key := range slices.Sorted(maps.Keys(state.Paths)) {
		record := state.Paths[key]
		if key == "" || filepath.IsAbs(key) || strings.Contains(key, "\\") || filepath.ToSlash(filepath.Clean(key)) != key || key == ".." || strings.HasPrefix(key, "../") {
			return nil, fmt.Errorf("reading %s: invalid record path %q", path, key)
		}
		if record == nil || record.Model == nil || record.Effort == nil || record.Model.LastGenerated == nil || record.Effort.LastGenerated == nil {
			return nil, fmt.Errorf("reading %s: incomplete record %q", path, key)
		}
	}
	return &state, nil
}
func serializeAgentModelState(state *agentModelState) ([]byte, error) {
	data, err := json.MarshalIndent(state, "", "  ")
	return append(data, '\n'), err
}
func readNativeModelFields(path, host string) (nativeModelFields, bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nativeModelFields{}, false, nil
	}
	if err != nil {
		return nativeModelFields{}, false, fmt.Errorf("reading %s: %w", path, err)
	}
	fields, err := parseNativeModelFields(data, host)
	if err != nil {
		return fields, true, fmt.Errorf("reading %s: %w", path, err)
	}
	return fields, true, nil
}

// Actual fields matching the last deployment are generated; mismatches capture
// local edits, including deletions, before explicit configuration can hide them.
func captureModelState(previous *agentModelRecord, actual nativeModelFields) *agentModelRecord {
	capture := func(previous *modelFieldState, actual modelFieldValue) *modelFieldState {
		next := &modelFieldState{}
		if previous != nil && *previous.LastGenerated == actual {
			next.RetainedLocal = previous.RetainedLocal
		} else if previous != nil || (actual.Present && actual.Value != "") {
			value := actual
			next.RetainedLocal = &value
		}
		return next
	}
	var model, effort *modelFieldState
	if previous != nil {
		model, effort = previous.Model, previous.Effort
	}
	return &agentModelRecord{Model: capture(model, actual.Model), Effort: capture(effort, actual.Effort)}
}

// Prune only paths that the existing host cleanup will actually remove. No
// snapshot key is ever used to discover or operate on a filesystem path.
func pruneCleanedModelState(root string, selected []hostTarget, defs []*subagent.Definition, state *agentModelState) error {
	enabled := map[string]bool{}
	for _, host := range selected {
		enabled[host.key] = true
	}
	for _, host := range allHostTargets() {
		if enabled[host.key] {
			continue
		}
		for _, def := range defs {
			rel := filepath.Join(host.relDir, def.Name+host.ext)
			if _, err := os.Stat(filepath.Join(root, rel)); err == nil {
				delete(state.Paths, filepath.ToSlash(rel))
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("checking %s: %w", filepath.Join(root, rel), err)
			}
		}
	}
	return nil
}
func writeAgentModelState(root string, data []byte) error {
	path := agentModelStatePath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".agent-model-state-*")
	if err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	name := temp.Name()
	defer os.Remove(name)
	if _, err := temp.Write(data); err != nil {
		_ = temp.Close()
		return fmt.Errorf("writing %s: %w", path, err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	if err := os.Rename(name, path); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
