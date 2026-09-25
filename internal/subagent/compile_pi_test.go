package subagent

import (
	"bytes"
	"strings"
	"testing"
)

// piModelFixture returns a stable definition whose zero-options PI output
// is pinned byte for byte below (ticket 02 hard acceptance: unconfigured
// projects must keep byte-identical PI artifacts — the model key is
// omitted and nothing else moves).
func piModelFixture() *Definition {
	return &Definition{
		Name:         "flowforge-fixture",
		Description:  "Fixture agent for pi model injection tests.",
		ModelProfile: ModelProfileToolCapable,
		DefaultSkill: "flowforge-fixture",
		Body:         "Fixture system prompt body.\n",
	}
}

// pinnedUnconfiguredPiOutput is the pre-ticket-02 CompilePi output for
// piModelFixture, captured from the implementation before the Model field
// landed (design d-pi-model-injection: regression zero for unconfigured
// projects). Any byte drift here fails the hard acceptance.
const pinnedUnconfiguredPiOutput = "---\nname: flowforge-fixture\ndescription: Fixture agent for pi model injection tests.\nthinking: medium\nskills:\n    - flowforge-fixture\ninheritSkills: false\n---\nFixture system prompt body.\n"

// TestCompilePiModelThreeStates pins the design's three-state contract:
// config Model → written, preserved FallbackModel only → written, neither
// → key omitted (inherit parent session model).
func TestCompilePiModelThreeStates(t *testing.T) {
	def := piModelFixture()

	tests := []struct {
		name      string
		opts      CompileOptions
		wantModel string // "" means the model key must be absent
	}{
		{"config model is written", CompileOptions{Model: "cpa/deepseek-v4.1-flash"}, "cpa/deepseek-v4.1-flash"},
		{"config model beats preserved fallback", CompileOptions{Model: "cpa/pinned/x", FallbackModel: "residue/y"}, "cpa/pinned/x"},
		{"preserved fallback is written when model absent", CompileOptions{FallbackModel: "residue/y"}, "residue/y"},
		{"unconfigured omits the model key", CompileOptions{}, ""},
	}
	for _, tt := range tests {
		got, err := CompilePiWithOptions(def, tt.opts)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		fm, _, ok := splitFrontmatter(got)
		if !ok {
			t.Fatalf("%s: frontmatter delimiters missing", tt.name)
		}
		if tt.wantModel == "" {
			if piFrontmatterHasKey(fm, "model") {
				t.Errorf("%s: model key must be omitted, got frontmatter:\n%s", tt.name, fm)
			}
			continue
		}
		wantLine := "model: " + tt.wantModel
		if !strings.Contains(string(fm), wantLine+"\n") && !strings.HasSuffix(string(fm), wantLine) {
			t.Errorf("%s: frontmatter missing %q, got:\n%s", tt.name, wantLine, fm)
		}
	}
}

// TestCompilePiUnconfiguredOutputByteIdentical is the regression-zero pin:
// zero options must reproduce the pre-change output byte for byte.
func TestCompilePiUnconfiguredOutputByteIdentical(t *testing.T) {
	got, err := CompilePi(piModelFixture())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte(pinnedUnconfiguredPiOutput)) {
		t.Errorf("unconfigured PI output drifted from the pre-change bytes\n got: %q\nwant: %q", got, pinnedUnconfiguredPiOutput)
	}
	// The zero-options path must also be unaffected by other inexpressible
	// options (EditDeny / DenyQuestion stay ignored by the PI compiler).
	withInexpressible, err := CompilePiWithOptions(piModelFixture(), CompileOptions{EditDeny: []string{"ignored"}, DenyQuestion: true})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(withInexpressible, []byte(pinnedUnconfiguredPiOutput)) {
		t.Errorf("inexpressible options must not change PI output\n got: %q\nwant: %q", withInexpressible, pinnedUnconfiguredPiOutput)
	}
}

// TestCompilePiModelFieldOrder pins the frontmatter field order contract:
// model sits between description and thinking (struct and assignment must
// stay in sync).
func TestCompilePiModelFieldOrder(t *testing.T) {
	got, err := CompilePiWithOptions(piModelFixture(), CompileOptions{Model: "cpa/order/x"})
	if err != nil {
		t.Fatal(err)
	}
	fm, _, ok := splitFrontmatter(got)
	if !ok {
		t.Fatal("frontmatter delimiters missing")
	}
	idx := func(key string) int {
		for i, line := range strings.Split(string(fm), "\n") {
			if strings.HasPrefix(line, key+":") {
				return i
			}
		}
		return -1
	}
	for _, key := range []string{"name", "description", "model", "thinking", "skills", "inheritSkills"} {
		if idx(key) == -1 {
			t.Errorf("frontmatter missing %q key:\n%s", key, fm)
		}
	}
	if d, m, th := idx("description"), idx("model"), idx("thinking"); !(d < m && m < th) {
		t.Errorf("model must sit between description and thinking (description=%d, model=%d, thinking=%d):\n%s", d, m, th, fm)
	}
}

// piFrontmatterHasKey reports whether a frontmatter block carries a key at
// line start (substring checks would false-positive on prose descriptions).
func piFrontmatterHasKey(fm []byte, key string) bool {
	for _, line := range strings.Split(string(fm), "\n") {
		if strings.HasPrefix(line, key+":") {
			return true
		}
	}
	return false
}
