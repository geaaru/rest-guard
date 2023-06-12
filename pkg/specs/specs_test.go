/*
Copyright © 2021-2023 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package specs_test

import (
	"github.com/geaaru/rest-guard/pkg/specs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rest Guard Specs Test", func() {

	Describe("Create struct", func() {

		Context("Ticket", func() {
			service := specs.NewRestService("service1")
			t := service.GetTicket()
			defer t.Rip()

			node := specs.NewRestNode("g1", "www.google.it", true)

			var i interface{} = node
			t.SetClosure("node", i)

			var pointer interface{} = nil
			var ok bool = false

			fClose := func(t *specs.RestTicket) {
				pointer, ok = t.GetClosure("node")
			}

			t.SetRequestCloseCb(fClose)

			t.RequestCloseCb(t)

			It("Check callback and closure", func() {
				Expect(t.RequestCloseCb).ShouldNot(BeNil())
				Expect(pointer).ShouldNot(BeNil())
				Expect(ok).Should(Equal(true))
				Expect(pointer).Should(Equal(node))
			})

		})

	})

})
