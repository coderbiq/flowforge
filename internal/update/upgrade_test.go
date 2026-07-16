package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/minio/selfupdate"
)

func makeTarGz(t *testing.T, files map[string][]byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	for name, content := range files {
		if err := tw.WriteHeader(&tar.Header{
			Name:     name,
			Size:     int64(len(content)),
			Mode:     0755,
			Typeflag: tar.TypeReg,
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(content); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func makeZip(t *testing.T, files map[string][]byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(content); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestExtractBinary_TarGz(t *testing.T) {
	bin := []byte("fake-linux-binary-content")
	archive := makeTarGz(t, map[string][]byte{
		"flowforge": bin,
	})

	result, err := extractBinary(archive, "linux-amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(result, bin) {
		t.Errorf("extracted binary does not match original")
	}
}

func TestExtractBinary_TarGzWithNestedPath(t *testing.T) {
	bin := []byte("fake-linux-binary-content")
	archive := makeTarGz(t, map[string][]byte{
		"some/dir/flowforge": bin,
	})

	result, err := extractBinary(archive, "linux-amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(result, bin) {
		t.Errorf("extracted binary does not match original")
	}
}

func TestExtractBinary_Zip(t *testing.T) {
	bin := []byte("fake-windows-binary-content")
	archive := makeZip(t, map[string][]byte{
		"flowforge.exe": bin,
	})

	result, err := extractBinary(archive, "windows-amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(result, bin) {
		t.Errorf("extracted binary does not match original")
	}
}

func TestExtractBinary_ZipWithNestedPath(t *testing.T) {
	bin := []byte("fake-windows-binary-content")
	archive := makeZip(t, map[string][]byte{
		"flowforge.exe": bin,
	})

	result, err := extractBinary(archive, "windows-amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(result, bin) {
		t.Errorf("extracted binary does not match original")
	}
}

func TestExtractBinary_RawBinary(t *testing.T) {
	bin := []byte{0x7f, 'E', 'L', 'F', 0x02, 0x01, 0x01, 0x00}

	result, err := extractBinary(bin, "linux-amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(result, bin) {
		t.Errorf("raw binary should be returned as-is")
	}
}

func TestExtractBinary_BinaryNotFound_TarGz(t *testing.T) {
	archive := makeTarGz(t, map[string][]byte{
		"other-file": []byte("not-the-binary"),
	})

	_, err := extractBinary(archive, "linux-amd64")
	if err == nil {
		t.Fatal("expected error for missing binary in tar.gz")
	}
}

func TestExtractBinary_BinaryNotFound_Zip(t *testing.T) {
	archive := makeZip(t, map[string][]byte{
		"other-file.exe": []byte("not-the-binary"),
	})

	_, err := extractBinary(archive, "windows-amd64")
	if err == nil {
		t.Fatal("expected error for missing binary in zip")
	}
}

func TestExtractBinary_CorruptGzip(t *testing.T) {
	corrupt := []byte{0x1f, 0x8b, 0xFF, 0xFF}
	_, err := extractBinary(corrupt, "linux-amd64")
	if err == nil {
		t.Fatal("expected error for corrupt gzip")
	}
}

func TestExtractBinary_CorruptZip(t *testing.T) {
	corrupt := []byte{0x50, 0x4b, 0xFF, 0xFF}
	_, err := extractBinary(corrupt, "linux-amd64")
	if err == nil {
		t.Fatal("expected error for corrupt zip")
	}
}

func TestExtractBinary_TarGzHonorsWindowsBinName(t *testing.T) {
	bin := []byte("fake-exe-content")
	archive := makeTarGz(t, map[string][]byte{
		"flowforge.exe": bin,
	})

	result, err := extractBinary(archive, "windows-amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(result, bin) {
		t.Errorf("extracted binary does not match original")
	}
}

func TestUpgradeFlow_ArchiveExtractedNotWrittenDirectly(t *testing.T) {
	bin := []byte("test-binary-content-that-should-be-the-result")
	archive := makeTarGz(t, map[string][]byte{
		"flowforge": bin,
	})
	sha := sha256Sum(archive)

	extracted, err := extractBinary(archive, "linux-amd64")
	if err != nil {
		t.Fatalf("extractBinary: %v", err)
	}

	if bytes.Equal(extracted, archive) {
		t.Fatal("extracted binary must NOT equal archive bytes (archive was not decompressed)")
	}
	if !bytes.Equal(extracted, bin) {
		t.Fatal("extracted binary must equal the original binary content inside archive")
	}
	if sha256Sum(archive) != sha {
		t.Fatal("SHA256 of archive should not change")
	}
}

func TestE2E_SelfupdateApply_ReplacesBinaryCorrectly(t *testing.T) {
	bin, err := os.ReadFile("/proc/self/exe")
	if err != nil {
		t.Skipf("cannot read /proc/self/exe: %v", err)
	}

	archive := makeTarGz(t, map[string][]byte{
		"flowforge": bin,
	})

	sha := sha256Sum(archive)

	mux := http.NewServeMux()
	mux.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"version": "999.0.0",
			"published_at": "2025-01-01T00:00:00Z",
			"artifacts": [{
				"platform": "linux-amd64",
				"url": "/artifact.tar.gz",
				"sha256": "` + sha + `"
			}]
		}`))
	})
	mux.HandleFunc("/artifact.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(archive)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	manifest, err := FetchManifest(srv.URL + "/manifest.json")
	if err != nil {
		t.Fatalf("fetch manifest: %v", err)
	}
	manifest.Artifacts[0].URL = srv.URL + "/artifact.tar.gz"

	targetDir, err := os.MkdirTemp("", "flowforge-e2e-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(targetDir)

	targetPath := filepath.Join(targetDir, "flowforge")
	if err := os.WriteFile(targetPath, []byte("old-stub-binary"), 0755); err != nil {
		t.Fatal(err)
	}

	platform := currentPlatform()
	artifact, ok := manifest.ArtifactByPlatform(platform)
	if !ok {
		t.Fatalf("no artifact for platform %s", platform)
	}

	client := &http.Client{}
	resp, err := client.Get(artifact.URL)
	if err != nil {
		t.Fatalf("download: %v", err)
	}
	defer resp.Body.Close()

	downloaded, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read downloaded: %v", err)
	}

	if err := verifySHA256(downloaded, artifact.SHA256); err != nil {
		t.Fatalf("sha256 verify: %v", err)
	}

	extracted, err := extractBinary(downloaded, platform)
	if err != nil {
		t.Fatalf("extractBinary: %v", err)
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		if err := selfupdate.Apply(bytes.NewReader(extracted), selfupdate.Options{
			TargetPath: targetPath,
		}); err != nil {
			t.Fatalf("selfupdate.Apply: %v", err)
		}
	} else {
		if err := os.WriteFile(targetPath, extracted, 0755); err != nil {
			t.Fatalf("write extracted binary: %v", err)
		}
	}

	resultBytes, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("read result binary: %v", err)
	}

	originalSHA := sha256Sum(bin)
	resultSHA := sha256Sum(resultBytes)
	if originalSHA != resultSHA {
		t.Fatalf("result binary corrupted: original SHA256=%s, result SHA256=%s, len(original)=%d, len(result)=%d",
			originalSHA, resultSHA, len(bin), len(resultBytes))
	}

	if bytes.HasPrefix(resultBytes, []byte{0x1f, 0x8b}) {
		t.Fatal("result binary starts with gzip magic bytes — archive was NOT decompressed!")
	}

	if bytes.HasPrefix(resultBytes, []byte{0x50, 0x4b}) {
		t.Fatal("result binary starts with zip magic bytes — archive was NOT decompressed!")
	}
}

func TestFullUpgradePipeline(t *testing.T) {
	bin, err := os.ReadFile("/proc/self/exe")
	if err != nil {
		t.Skipf("cannot read /proc/self/exe: %v", err)
	}

	archive := makeTarGz(t, map[string][]byte{
		"flowforge": bin,
	})

	sha := sha256Sum(archive)

	mux := http.NewServeMux()
	mux.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"version": "999.0.0",
			"published_at": "2025-01-01T00:00:00Z",
			"artifacts": [{
				"platform": "linux-amd64",
				"url": "/artifact.tar.gz",
				"sha256": "` + sha + `"
			}]
		}`))
	})
	mux.HandleFunc("/artifact.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")
		w.Write(archive)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	manifest, err := FetchManifest(srv.URL + "/manifest.json")
	if err != nil {
		t.Fatalf("fetch manifest: %v", err)
	}

	if len(manifest.Artifacts) == 0 {
		t.Fatal("expected artifacts in manifest")
	}
	manifest.Artifacts[0].URL = srv.URL + "/artifact.tar.gz"

	targetDir, err := os.MkdirTemp("", "flowforge-upgrade-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(targetDir)

	targetPath := filepath.Join(targetDir, "flowforge")
	if err := os.WriteFile(targetPath, []byte("old-binary"), 0755); err != nil {
		t.Fatal(err)
	}

	platform := currentPlatform()
	artifact, ok := manifest.ArtifactByPlatform(platform)
	if !ok {
		t.Fatalf("no artifact for platform %s", platform)
	}

	client := &http.Client{}
	resp, err := client.Get(artifact.URL)
	if err != nil {
		t.Fatalf("download: %v", err)
	}
	defer resp.Body.Close()

	downloaded, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read downloaded: %v", err)
	}

	if err := verifySHA256(downloaded, artifact.SHA256); err != nil {
		t.Fatalf("sha256 verify: %v", err)
	}

	extracted, err := extractBinary(downloaded, platform)
	if err != nil {
		t.Fatalf("extractBinary: %v", err)
	}

	if bytes.Equal(extracted, downloaded) {
		t.Fatal("extracted bytes must not equal downloaded archive bytes")
	}

	if !bytes.Equal(extracted, bin) {
		t.Fatalf("extracted binary does not match original binary: len(original)=%d, len(extracted)=%d, sha256(original)=%s, sha256(extracted)=%s",
			len(bin), len(extracted), sha256Sum(bin), sha256Sum(extracted))
	}

	h := sha256.New()
	h.Write(extracted)
	extractedSHA := hex.EncodeToString(h.Sum(nil))

	originalSHA := sha256Sum(bin)
	if extractedSHA != originalSHA {
		t.Fatalf("SHA256 mismatch after extraction: original=%s, extracted=%s", originalSHA, extractedSHA)
	}
}
