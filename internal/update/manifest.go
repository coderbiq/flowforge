package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ManifestArtifact struct {
	Platform     string `json:"platform"`
	URL          string `json:"url"`
	SHA256       string `json:"sha256"`
	Size         int64  `json:"size_bytes"`
	SignatureURL string `json:"signature_url,omitempty"`
}

type Manifest struct {
	Version     string             `json:"version"`
	PublishedAt string             `json:"published_at"`
	ReleaseNotes string            `json:"release_notes,omitempty"`
	Artifacts   []ManifestArtifact `json:"artifacts"`
}

func FetchManifest(url string) (*Manifest, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching manifest from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest fetch returned status %d", resp.StatusCode)
	}

	var m Manifest
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("decoding manifest: %w", err)
	}

	if m.Version == "" {
		return nil, fmt.Errorf("manifest missing version field")
	}

	return &m, nil
}

func (m *Manifest) ArtifactByPlatform(platform string) (*ManifestArtifact, bool) {
	for i := range m.Artifacts {
		if m.Artifacts[i].Platform == platform {
			return &m.Artifacts[i], true
		}
	}
	return nil, false
}
