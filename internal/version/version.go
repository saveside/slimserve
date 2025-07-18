package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

// Build-time variables injected via ldflags
var (
	Version    = "dev"
	CommitHash = "unknown"
	BuildDate  = "unknown"
	BuildUser  = "unknown"
)

type Info struct {
	Version    string `json:"version"`
	CommitHash string `json:"commit_hash"`
	BuildDate  string `json:"build_date"`
	BuildUser  string `json:"build_user"`
	GoVersion  string `json:"go_version"`
	Platform   string `json:"platform"`
	Arch       string `json:"arch"`
}

var (
	cachedInfo  *Info
	cachedShort string
)

func init() {
	cachedInfo = &Info{
		Version:    Version,
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
		BuildUser:  BuildUser,
		GoVersion:  runtime.Version(),
		Platform:   runtime.GOOS,
		Arch:       runtime.GOARCH,
	}

	if cachedInfo.CommitHash != "unknown" && len(cachedInfo.CommitHash) > 7 {
		cachedShort = fmt.Sprintf("%s (%s)", cachedInfo.Version, cachedInfo.CommitHash[:7])
	} else {
		cachedShort = cachedInfo.Version
	}
}

// Get returns the current version information (cached)
func Get() Info {
	return *cachedInfo
}

// String returns a human-readable version string
func (i Info) String() string {
	if i.CommitHash != "unknown" && len(i.CommitHash) > 7 {
		return fmt.Sprintf("%s (%s)", i.Version, i.CommitHash[:7])
	}
	return i.Version
}

// JSON returns the version information as JSON
func (i Info) JSON() ([]byte, error) {
	return json.MarshalIndent(i, "", "  ")
}

// GetShort returns a short version string for display (cached)
func GetShort() string {
	return cachedShort
}

// GetBuildTime returns the build time as a time.Time if parseable
func GetBuildTime() (time.Time, error) {
	if BuildDate == "unknown" {
		return time.Time{}, fmt.Errorf("build date unknown")
	}

	// Try RFC3339 format first (ISO 8601)
	if t, err := time.Parse(time.RFC3339, BuildDate); err == nil {
		return t, nil
	}

	// Try common build date formats
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, BuildDate); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse build date: %s", BuildDate)
}
