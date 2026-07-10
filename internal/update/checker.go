package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const versionCheckDebounce = 1 * time.Hour

type VersionChecker struct {
	currentVersion string
	getCheck       func() (time.Time, string, error)
	setCheck       func(version string) error
	httpClient     *http.Client
}

type VersionCheckStore interface {
	GetVersionCheck() (time.Time, string, error)
	SetVersionCheck(version string) error
}

func NewVersionChecker(currentVersion string, store VersionCheckStore) *VersionChecker {
	return &VersionChecker{
		currentVersion: currentVersion,
		getCheck:       store.GetVersionCheck,
		setCheck:       store.SetVersionCheck,
		httpClient:     &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *VersionChecker) CheckAsync(notify func(string)) {
	go func() {
		msg, err := c.check()
		if err != nil || msg == "" {
			return
		}
		if notify != nil {
			notify(msg)
		}
	}()
}

func (c *VersionChecker) check() (string, error) {
	lastCheck, lastVersion, err := c.getCheck()
	if err != nil {
		return "", err
	}

	if time.Since(lastCheck) < versionCheckDebounce && lastVersion != "" {
		return "", nil
	}

	latestVersion, err := c.fetchLatestVersion()
	if err != nil {
		return "", err
	}

	if err := c.setCheck(latestVersion); err != nil {
		return "", err
	}

	if CompareVersions(latestVersion, c.currentVersion) > 0 {
		return fmt.Sprintf("新版本 %s 可用，运行 flowforge upgrade 升级", latestVersion), nil
	}

	return "", nil
}

func (c *VersionChecker) fetchLatestVersion() (string, error) {
	latestTag, err := resolveLatestTag()
	if err != nil {
		return "", err
	}

	manifestURL := fmt.Sprintf("https://github.com/coderbiq/flowforge/releases/download/%s/manifest.json", latestTag)
	resp, err := c.httpClient.Get(manifestURL)
	if err != nil {
		return "", fmt.Errorf("fetching manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("manifest fetch returned status %d", resp.StatusCode)
	}

	var manifest struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return "", fmt.Errorf("decoding manifest: %w", err)
	}

	if manifest.Version == "" {
		return "", fmt.Errorf("manifest missing version field")
	}

	return manifest.Version, nil
}

func CompareVersions(a, b string) int {
	a = strings.TrimPrefix(a, "v")
	b = strings.TrimPrefix(b, "v")

	aParts := splitVersion(a)
	bParts := splitVersion(b)

	for i := 0; i < 3; i++ {
		if aParts[i] > bParts[i] {
			return 1
		}
		if aParts[i] < bParts[i] {
			return -1
		}
	}
	return 0
}

func splitVersion(v string) [3]int {
	parts := strings.SplitN(v, ".", 3)
	var nums [3]int
	for i := 0; i < len(parts) && i < 3; i++ {
		n, _ := strconv.Atoi(parts[i])
		nums[i] = n
	}
	return nums
}
