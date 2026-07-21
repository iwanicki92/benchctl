// SPDX-FileCopyrightText: 2026 3mdeb <contact@3mdeb.com>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"runtime/debug"

	"github.com/zarhus/benchctl/platform"
)

// versionLine is the full --version text: the VCS stamp plus the binary's
// default platform, which tells the build flavors apart.
func versionLine() string {
	v := version()
	if def := platform.Default(); def != "" {
		v += " (" + def + " build)"
	}
	return v
}

// version returns a build identifier read from the VCS stamp Go records for any
// build or install made in a git checkout: the commit revision, suffixed with
// "-dirty" when the working tree had uncommitted changes. It falls back to "dev"
// when no stamp is present, such as under `go run`.
func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return formatVersion("", false)
	}
	var revision string
	var modified bool
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}
	return formatVersion(revision, modified)
}

// shortLen is how many leading hex characters of a commit identify a build.
// Twelve is unambiguous in practice while staying readable.
const shortLen = 12

// formatVersion renders the version string from the raw VCS stamp fields, so the
// formatting is testable without a real build. An empty revision means the build
// carries no stamp. The commit is abbreviated to its leading characters and gets
// a "-dirty" suffix when the working tree had uncommitted changes at build time.
func formatVersion(revision string, modified bool) string {
	if revision == "" {
		return "dev"
	}
	if len(revision) > shortLen {
		revision = revision[:shortLen]
	}
	if modified {
		return revision + "-dirty"
	}
	return revision
}
