// SPDX-FileCopyrightText: 2026 3mdeb <contact@3mdeb.com>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zarhus/benchctl/platform"
)

func consoleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "console",
		Short: "Attach to the host serial console (detaches on exit)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return withPlatform(func(p platform.Platform) error { return p.Console() })
		},
	}
}
