package config

import (
	"fmt"
	"strings"
)

type sideEffectFunc func(svc *ConfigService, oldValue, newValue string) error

type sideEffectRegistry struct {
	effects []sideEffectEntry
}

type sideEffectEntry struct {
	pattern string
	fn      sideEffectFunc
}

func newSideEffectRegistry() *sideEffectRegistry {
	return &sideEffectRegistry{}
}

func (r *sideEffectRegistry) register(pattern string, fn sideEffectFunc) {
	r.effects = append(r.effects, sideEffectEntry{pattern: pattern, fn: fn})
}

func (r *sideEffectRegistry) trigger(svc *ConfigService, key, oldValue, newValue string) error {
	for _, e := range r.effects {
		if matchPattern(e.pattern, key) {
			if err := e.fn(svc, oldValue, newValue); err != nil {
				return fmt.Errorf("side effect %q failed: %w", e.pattern, err)
			}
		}
	}
	return nil
}

func matchPattern(pattern, key string) bool {
	pp := strings.Split(pattern, ".")
	kp := strings.Split(key, ".")
	if len(pp) != len(kp) {
		return false
	}
	for i := range pp {
		if pp[i] == "*" {
			continue
		}
		if pp[i] != kp[i] {
			return false
		}
	}
	return true
}
