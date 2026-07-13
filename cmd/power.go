package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zarhus/benchctl/platform"
)

func powerCmd() *cobra.Command {
	power := &cobra.Command{
		Use:   "power",
		Short: "Host power control",
	}

	on := &cobra.Command{
		Use:   "on",
		Short: "Power on and wait until the host is on",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return withPlatform(func(p platform.Platform) error { return p.SetPower(platform.PowerOn) })
		},
	}

	off := &cobra.Command{
		Use:   "off",
		Short: "Power off and wait until the host is off",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return withPlatform(func(p platform.Platform) error { return p.SetPower(platform.PowerOff) })
		},
	}

	status := &cobra.Command{
		Use:   "status",
		Short: "Print the current power state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return withPlatform(func(p platform.Platform) error {
				status, err := p.PowerState()
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), status)
				return nil
			})
		},
	}

	var hard bool
	reset := &cobra.Command{
		Use:   "reset",
		Short: "Power-cycle the host off and back on",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return withPlatform(func(p platform.Platform) error {
				if hard {
					return p.HardReset()
				}
				return p.PowerReset()
			})
		},
	}
	reset.Flags().BoolVar(&hard, "hard", false, "hard-reset through the platform's low-level path; use after an in-OS reboot leaves the host wedged")

	power.AddCommand(on, off, status, reset)
	return power
}
