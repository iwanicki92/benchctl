// SPDX-FileCopyrightText: 2026 3mdeb <contact@3mdeb.com>
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zarhus/benchctl/platform"
)

// parseFlashTarget maps an optional [host|bmc] argument to a FlashTarget,
// defaulting to host when no argument is given.
func parseFlashTarget(args []string) (platform.FlashTarget, error) {
	if len(args) == 0 {
		return platform.FlashHost, nil
	}
	switch args[0] {
	case "host":
		return platform.FlashHost, nil
	case "bmc":
		return platform.FlashBMC, nil
	default:
		return 0, fmt.Errorf("unknown flash target %q (want host or bmc)", args[0])
	}
}

func flashCmd() *cobra.Command {
	flash := &cobra.Command{
		Use:   "flash",
		Short: "Flash firmware and manage updates",
	}

	var hostForce bool
	host := &cobra.Command{
		Use:   "host <firmware>",
		Short: "Flash the host boot flash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return withPlatform(func(p platform.Platform) error {
				return p.Flash(platform.FlashHost, args[0], hostForce)
			})
		},
	}
	host.Flags().BoolVar(&hostForce, "force", false, "skip the firmware size check")

	var bmcForce bool
	bmc := &cobra.Command{
		Use:   "bmc <firmware>",
		Short: "Flash the BMC flash (where the platform supports it)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return withPlatform(func(p platform.Platform) error {
				return p.Flash(platform.FlashBMC, args[0], bmcForce)
			})
		},
	}
	bmc.Flags().BoolVar(&bmcForce, "force", false, "skip the firmware size check")

	status := &cobra.Command{
		Use:   "status [host|bmc]",
		Short: "Show the update status (default host)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := parseFlashTarget(args)
			if err != nil {
				return err
			}
			return withPlatform(func(p platform.Platform) error {
				status, err := p.FlashStatus(target)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "state: %s\n", status.State)
				return nil
			})
		},
	}

	abort := &cobra.Command{
		Use:   "abort [host|bmc]",
		Short: "Abort a stuck or stale update (default host)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := parseFlashTarget(args)
			if err != nil {
				return err
			}
			return withPlatform(func(p platform.Platform) error {
				return p.FlashAbort(target)
			})
		},
	}

	flash.AddCommand(host, bmc, status, abort)
	return flash
}
