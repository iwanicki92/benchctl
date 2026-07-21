// SPDX-FileCopyrightText: 2026 3mdeb <contact@3mdeb.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package rte is a client for the RTE controller's REST API, which exposes the
// bench's GPIO pins by logical id at http://<host>:8000/api/v1/gpio/{id}.
//
// The RTE has two kinds of pins. Ids 1-12 are open-collector: they can only
// pull "low" (asserted) or release to "high-z". Ids 0 and 13-19 are push-pull:
// they drive a clean "high" or "low". Set takes the pin's semantic state and
// maps it to the wire encoding the API expects.
package rte

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// Pin is a GPIO pin's reported state.
type Pin struct {
	ID          int    `json:"id"`
	State       uint   `json:"state"`
	Direction   string `json:"direction"`
	Description string `json:"description"`
}

// Asserted reports whether the pin is driven to its active level: low for an
// open-collector pin, high for a push-pull pin. The RTE carries this in the low
// bit of State.
func (p Pin) Asserted() bool {
	return p.State%2 == 1
}

// Client talks to one RTE's REST API.
type Client struct {
	base string
	hc   *http.Client
}

// New returns a client for the RTE reachable at host. net.JoinHostPort brackets
// an IPv6 literal, so a host such as fd00::1 yields a valid URL authority.
func New(host string) *Client {
	return &Client{
		base: fmt.Sprintf("http://%s/api/v1", net.JoinHostPort(host, "8000")),
		hc:   &http.Client{Timeout: 30 * time.Second},
	}
}

// isOpenCollector reports whether id is an open-collector pin (ids 1-12), as
// opposed to a push-pull pin (ids 0 and 13-19).
func isOpenCollector(id int) bool {
	return id >= 1 && id <= 12
}

// encodeState maps a semantic state to the API's integer encoding, rejecting a
// state that does not apply to the pin's type.
func encodeState(id int, state string) (uint, error) {
	if isOpenCollector(id) {
		switch state {
		case "low":
			return 1, nil
		case "high-z":
			return 0, nil
		default:
			return 0, fmt.Errorf("gpio %d is open-collector: state must be low or high-z, got %q", id, state)
		}
	}
	switch state {
	case "high":
		return 1, nil
	case "low":
		return 0, nil
	default:
		return 0, fmt.Errorf("gpio %d is push-pull: state must be high or low, got %q", id, state)
	}
}

// Get returns the current state of pin id.
func (c *Client) Get(id int) (Pin, error) {
	resp, err := c.hc.Get(fmt.Sprintf("%s/gpio/%d", c.base, id))
	if err != nil {
		return Pin{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Pin{}, fmt.Errorf("get gpio %d: %s", id, apiError(resp))
	}
	var pin Pin
	if err := json.NewDecoder(resp.Body).Decode(&pin); err != nil {
		return Pin{}, fmt.Errorf("parse gpio %d: %w", id, err)
	}
	return pin, nil
}

// Set drives pin id to state as an output. For holdSecs > 0 the RTE reverts the
// pin to its opposite state after that many seconds, which pulses a momentary
// line such as a power or reset button.
func (c *Client) Set(id int, state string, holdSecs int) error {
	encoded, err := encodeState(id, state)
	if err != nil {
		return err
	}
	body, err := json.Marshal(struct {
		State     uint   `json:"state"`
		Direction string `json:"direction"`
		Time      uint   `json:"time"`
	}{State: encoded, Direction: "out", Time: uint(holdSecs)})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/gpio/%d", c.base, id), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set gpio %d %s: %s", id, state, apiError(resp))
	}
	return nil
}

// apiError renders a failed response for an error message, preferring the JSON
// error field the API returns and falling back to the status.
func apiError(resp *http.Response) string {
	var parsed struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err == nil && parsed.Error != "" {
		return parsed.Error
	}
	return strings.TrimSpace(resp.Status)
}
