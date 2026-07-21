// SPDX-FileCopyrightText: 2026 3mdeb <contact@3mdeb.com>
//
// SPDX-License-Identifier: Apache-2.0

// Command benchctl controls test benches from a workstation or on the bench
// itself, through a platform driver selected at runtime.
package main

import (
	"os"

	"github.com/zarhus/benchctl/cmd"
	"github.com/zarhus/benchctl/platform"
)

func main() {
	platform.SetDefault("benchrack")
	os.Exit(cmd.Main())
}
