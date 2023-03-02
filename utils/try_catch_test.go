package utils

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {

})

var _ = Describe("try-catch tests", func() {

	It("Checking try-catch", func() {

		thrInfo := "<Thrown exception info>"
		var caught string
		var finally bool

		TryBlock{
			Try: func() {
				fmt.Println("I tried")
				Throw(thrInfo)
			},
			Catch: func(e Exception) {
				fmt.Printf("Caught:  %v\n", e)
				caught = fmt.Sprintf("%v", e)

			},
			Finally: func() {
				fmt.Println("Finally...")
				finally = true
			},
		}.Do()

		Expect(caught).Should(Equal(thrInfo))
		Expect(finally).Should(BeTrue())

	})

})

func TestTemplate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils tests")
}
