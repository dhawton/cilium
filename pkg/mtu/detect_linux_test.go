// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package mtu

import (
	"fmt"
	"os/exec"

	"github.com/cilium/cilium/pkg/testutils"

	. "gopkg.in/check.v1"
)

func (m *MTUSuite) TestAutoDetect(c *C) {
	testutils.PrivilegedCheck(c)

	mtu, err := autoDetect()
	if err != nil {
		fmt.Printf("MTU auto detection failed: %s, retrying...\n", err)
		mtu, err = autoDetect()
	} else {
		fmt.Printf("MTU auto detection worked\n")
	}
	if err != nil {
		// Execute "ip route show all"
		showCmd := exec.Command("ip", "route", "show", "all")
		showOutput, err := showCmd.Output()
		c.Assert(err, IsNil)
		fmt.Println("ip route show all output:")
		fmt.Println(string(showOutput))

		// Execute "ip rule list"
		showCmd = exec.Command("ip", "rule", "list")
		showOutput, err = showCmd.Output()
		c.Assert(err, IsNil)
		fmt.Println("ip route route list output:")
		fmt.Println(string(showOutput))

		// Execute "ip route show table local"
		showCmd = exec.Command("ip", "route", "show", "table", "local")
		showOutput, err = showCmd.Output()
		c.Assert(err, IsNil)
		fmt.Println("ip route show table local output:")
		fmt.Println(string(showOutput))

		// Execute "ip link"
		showCmd = exec.Command("ip", "link")
		showOutput, err = showCmd.Output()
		c.Assert(err, IsNil)
		fmt.Println("ip link output:")
		fmt.Println(string(showOutput))

		// Execute "ip route get 1.1.1.1"
		getCmd := exec.Command("ip", "route", "get", "1.1.1.1")
		getOutput, err := getCmd.Output()
		c.Assert(err, IsNil)
		fmt.Println("ip route get 1.1.1.1 output:")
		fmt.Println(string(getOutput))
		mtu, err = autoDetect()
		if err != nil {
			fmt.Printf("MTU auto detection failed: %s, retrying...\n", err)
			mtu, err = autoDetect()
		} else {
			fmt.Printf("MTU auto detection worked\n")
		}
	}
	c.Assert(err, IsNil)
	c.Assert(mtu, Not(Equals), 0)
}
