package update

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/minio/selfupdate"
)

const releasesBaseURL = "https://github.com/coderbiq/flowforge/releases"

func manifestURL(version string) string {
	if version == "latest" {
		return releasesBaseURL + "/latest/download/manifest.json"
	}
	return releasesBaseURL + "/download/" + version + "/manifest.json"
}

func currentPlatform() string {
	return fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
}

type UpgradeResult struct {
	OldVersion string
	NewVersion string
}

func Upgrade(currentVersion string) (*UpgradeResult, error) {
	murl := manifestURL("latest")
	manifest, err := FetchManifest(murl)
	if err != nil {
		return nil, fmt.Errorf("upgrade: %w", err)
	}

	return UpgradeToVersion(manifest, currentVersion, manifest.Version)
}

func UpgradeToVersion(manifest *Manifest, currentVersion, targetVersion string) (*UpgradeResult, error) {
	if targetVersion != manifest.Version {
		murl := manifestURL(targetVersion)
		var err error
		manifest, err = FetchManifest(murl)
		if err != nil {
			return nil, fmt.Errorf("upgrade to %s: %w", targetVersion, err)
		}
	}

	if CompareVersions(manifest.Version, currentVersion) <= 0 && targetVersion != currentVersion {
		return nil, fmt.Errorf("target version %s is not newer than current version %s", manifest.Version, currentVersion)
	}

	platform := currentPlatform()
	artifact, ok := manifest.ArtifactByPlatform(platform)
	if !ok {
		return nil, fmt.Errorf("no artifact found for platform %s in manifest %s", platform, manifest.Version)
	}

	client := &http.Client{Timeout: 5 * time.Minute}

	binResp, err := client.Get(artifact.URL)
	if err != nil {
		return nil, fmt.Errorf("downloading artifact: %w", err)
	}
	defer binResp.Body.Close()

	if binResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("artifact download returned status %d", binResp.StatusCode)
	}

	bin, err := io.ReadAll(binResp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading artifact bytes: %w", err)
	}

	if artifact.SignatureURL != "" {
		if err := verifyArtifactSignature(bin, artifact.SignatureURL, client); err != nil {
			return nil, fmt.Errorf("signature verification: %w", err)
		}
	}

	if err := verifySHA256(bin, artifact.SHA256); err != nil {
		return nil, fmt.Errorf("checksum verification: %w", err)
	}

	if err := selfupdate.Apply(bytes.NewReader(bin), selfupdate.Options{}); err != nil {
		return nil, fmt.Errorf("applying update: %w", err)
	}

	return &UpgradeResult{
		OldVersion: currentVersion,
		NewVersion: manifest.Version,
	}, nil
}

func DryRunUpgrade(currentVersion string) (*Manifest, error) {
	murl := manifestURL("latest")
	manifest, err := FetchManifest(murl)
	if err != nil {
		return nil, fmt.Errorf("dry-run: %w", err)
	}

	if CompareVersions(manifest.Version, currentVersion) <= 0 {
		return manifest, nil
	}

	return manifest, nil
}

func verifyArtifactSignature(bin []byte, sigURL string, client *http.Client) error {
	sigResp, err := client.Get(sigURL)
	if err != nil {
		return fmt.Errorf("fetching signature: %w", err)
	}
	defer sigResp.Body.Close()

	if sigResp.StatusCode != http.StatusOK {
		return fmt.Errorf("signature fetch returned status %d", sigResp.StatusCode)
	}

	sig, err := io.ReadAll(io.LimitReader(sigResp.Body, 1024))
	if err != nil {
		return fmt.Errorf("reading signature: %w", err)
	}

	if !VerifySignature(bin, sig) {
		return fmt.Errorf("Ed25519 signature verification failed")
	}

	return nil
}

func verifySHA256(data []byte, expected string) error {
	actual := sha256Sum(data)
	if actual != expected {
		return fmt.Errorf("SHA256 mismatch: expected %s, got %s", expected, actual)
	}
	return nil
}

func sha256Sum(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
