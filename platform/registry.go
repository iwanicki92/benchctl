// SPDX-FileCopyrightText: 2026 3mdeb <contact@3mdeb.com>
//
// SPDX-License-Identifier: Apache-2.0

package platform

import "sort"

// registry holds the registered platform drivers by name.
var registry = map[string]Spec{}

// defaultName is the driver used when no platform is selected. It is empty
// until a binary's main sets it.
var defaultName string

// Register adds a driver. Drivers call this from their init().
func Register(spec Spec) {
	registry[spec.Name] = spec
}

// Get returns the driver registered under name.
func Get(name string) (Spec, bool) {
	spec, ok := registry[name]
	return spec, ok
}

// SetDefault selects the driver used when no platform is chosen by flag or
// environment. A binary's main calls it for the platform it is built around.
func SetDefault(name string) {
	defaultName = name
}

// Default returns the driver name set by SetDefault, or empty when none is set.
func Default() string {
	return defaultName
}

// Names returns the registered driver names in sorted order.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
