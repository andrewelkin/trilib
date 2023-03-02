package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config parsing tests", func() {
	It("It should correctly initialize comma separated string slice", func() {
		c := &Config{}
		c.cfg = make(map[string]interface{})
		inputSliceString := "USD,USDC,DAI,ESD"
		expectedSlice := []string{"USD", "USDC", "DAI", "ESD"}
		c.cfg["example"] = inputSliceString
		outputSlice := c.GetStringList("example")
		Expect(outputSlice).Should(Equal(expectedSlice))
	})
})
