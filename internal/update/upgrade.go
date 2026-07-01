package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
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

	if artifact.Size > 0 && int64(len(bin)) != artifact.Size {
		return nil, fmt.Errorf("size mismatch: expected %d bytes, got %d", artifact.Size, len(bin))
	}

	bin, err = extractBinary(bin, platform)
	if err != nil {
		return nil, fmt.Errorf("extracting binary from artifact: %w", err)
	}

	if err := selfupdate.Apply(bytes.NewReader(bin), selfupdate.Options{}); err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			return nil, fmt.Errorf("applying update failed and rollback also failed (binary may be corrupted: %v); recovery: %w", rerr, err)
		}
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

func extractBinary(data []byte, platform string) ([]byte, error) {
	binName := "flowforge"
	if strings.Contains(platform, "windows") {
		binName = "flowforge.exe"
	}

	if len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b {
		return extractFromTarGz(data, binName)
	}

	if len(data) >= 4 && data[0] == 0x50 && data[1] == 0x4b {
		return extractFromZip(data, binName)
	}

	return data, nil
}

func extractFromTarGz(data []byte, binName string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("opening gzip: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading tar: %w", err)
		}
		if hdr.Typeflag == tar.TypeReg && (hdr.Name == binName || filepath.Base(hdr.Name) == binName) {
			return io.ReadAll(tr)
		}
	}
	return nil, fmt.Errorf("binary %s not found in tar.gz archive", binName)
}

func extractFromZip(data []byte, binName string) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("opening zip: %w", err)
	}
	for _, f := range zr.File {
		if f.Name == binName || filepath.Base(f.Name) == binName {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("opening zip entry: %w", err)
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("binary %s not found in zip archive", binName)
}
