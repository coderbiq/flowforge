package update

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const flowforgePublicKeyHex = "8d5cba55df8b8ad1c5dc33f4e365ff1db5737f3b02a4e43baa0d39e95302bb7f"

var flowforgePublicKey ed25519.PublicKey

func init() {
	key, err := hex.DecodeString(flowforgePublicKeyHex)
	if err != nil || len(key) != ed25519.PublicKeySize {
		log.Printf("flowforge: invalid public key constant, signature verification will fail")
		flowforgePublicKey = make(ed25519.PublicKey, ed25519.PublicKeySize)
	} else {
		flowforgePublicKey = ed25519.PublicKey(key)
	}
}

func VerifySignature(message, sig []byte) bool {
	return ed25519.Verify(flowforgePublicKey, message, sig)
}

func VerifyArtifact(artifact *ManifestArtifact) error {
	if artifact == nil {
		return fmt.Errorf("artifact is nil")
	}

	client := &http.Client{Timeout: 30 * time.Second}

	sigResp, err := client.Get(artifact.SignatureURL)
	if err != nil {
		return fmt.Errorf("fetching signature from %s: %w", artifact.SignatureURL, err)
	}
	defer sigResp.Body.Close()

	if sigResp.StatusCode != http.StatusOK {
		return fmt.Errorf("signature fetch returned status %d", sigResp.StatusCode)
	}

	sig, err := io.ReadAll(io.LimitReader(sigResp.Body, 1024))
	if err != nil {
		return fmt.Errorf("reading signature: %w", err)
	}

	binResp, err := client.Get(artifact.URL)
	if err != nil {
		return fmt.Errorf("fetching artifact from %s: %w", artifact.URL, err)
	}
	defer binResp.Body.Close()

	if binResp.StatusCode != http.StatusOK {
		return fmt.Errorf("artifact fetch returned status %d", binResp.StatusCode)
	}

	bin, err := io.ReadAll(binResp.Body)
	if err != nil {
		return fmt.Errorf("reading artifact: %w", err)
	}

	if !VerifySignature(bin, sig) {
		return fmt.Errorf("Ed25519 signature verification failed for %s", artifact.Platform)
	}

	expectedSHA256 := artifact.SHA256
	hasher := sha256.New()
	hasher.Write(bin)
	actualSHA256 := hex.EncodeToString(hasher.Sum(nil))

	if actualSHA256 != expectedSHA256 {
		return fmt.Errorf("SHA256 mismatch: expected %s, got %s", expectedSHA256, actualSHA256)
	}

	return nil
}
