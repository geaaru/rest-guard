/*
Copyright © 2021-2023 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package guard_test

import (
	"fmt"

	g "github.com/geaaru/rest-guard/pkg/guard"
	"github.com/geaaru/rest-guard/pkg/specs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Guard Test", func() {

	Describe("Create struct", func() {

		Context("Simple1", func() {

			c := specs.NewConfig()
			guard, err := g.NewRestGuard(c)

			service := specs.NewRestService("google")
			service.Retries = 1
			nodeFailed := specs.NewRestNode("gError", "127.0.0.1", true)
			node := specs.NewRestNode("g1", "www.google.it", true)

			guard.AddService(service.GetName(), service)

			errAdd1 := guard.AddRestNode(service.GetName(), nodeFailed)
			errAdd2 := guard.AddRestNode(service.GetName(), node)

			t := service.GetTicket()
			defer t.Rip()
			req, errReq := guard.CreateRequest(t, "GET", "/")
			errDo := guard.Do(t)
			It("Check nil", func() {
				Expect(guard).ShouldNot(BeNil())
				Expect(req).ShouldNot(BeNil())
			})
			It("Check err", func() {
				Expect(err).Should(BeNil())
				Expect(errAdd1).Should(BeNil())
				Expect(errAdd2).Should(BeNil())
				Expect(errReq).Should(BeNil())
				Expect(errDo).Should(BeNil())
			})
			It("User Agent", func() {
				Expect(guard.GetUserAgent()).Should(Equal(
					fmt.Sprintf("RestGuard v%s", specs.RGuardVersion)))
			})
			It("Check ticket", func() {
				Expect(t).ShouldNot(BeNil())
				Expect(t.Id).ShouldNot(Equal(""))
				Expect(t.Retries).Should(Equal(1))
			})
			It("Check Response ", func() {
				Expect(t.Response).ShouldNot(BeNil())
				Expect(t.Response.StatusCode).Should(Equal(200))
			})

		})

		Context("Simple2", func() {

			c := specs.NewConfig()
			guard, err := g.NewRestGuard(c)

			service := specs.NewRestService("google")
			service.Retries = 0
			nodeFailed := specs.NewRestNode("gError", "127.0.0.1", true)
			node := specs.NewRestNode("g1", "www.google.it", true)

			guard.AddService(service.GetName(), service)

			errAdd1 := guard.AddRestNode(service.GetName(), nodeFailed)
			errAdd2 := guard.AddRestNode(service.GetName(), node)

			t := service.GetTicket()
			defer t.Rip()
			req, errReq := guard.CreateRequest(t, "GET", "/")
			errDo := guard.Do(t)
			It("Check nil", func() {
				Expect(guard).ShouldNot(BeNil())
				Expect(req).ShouldNot(BeNil())
			})
			It("Check err", func() {
				Expect(err).Should(BeNil())
				Expect(errAdd1).Should(BeNil())
				Expect(errAdd2).Should(BeNil())
				Expect(errReq).Should(BeNil())
				Expect(errDo).ShouldNot(BeNil())
			})
			It("User Agent", func() {
				Expect(guard.GetUserAgent()).Should(Equal(
					fmt.Sprintf("RestGuard v%s", specs.RGuardVersion)))
			})
			It("Check ticket", func() {
				Expect(t).ShouldNot(BeNil())
				Expect(t.Id).ShouldNot(Equal(""))
				Expect(t.Retries).Should(Equal(0))
			})
			It("Check Response ", func() {
				Expect(t.Response).Should(BeNil())
			})

		})

		Context("Simple3 - Fail", func() {

			c := specs.NewConfig()
			guard, err := g.NewRestGuard(c)

			service := specs.NewRestService("google")
			service.Retries = 2
			nodeFailed := specs.NewRestNode("gError", "127.0.0.1", true)
			node := specs.NewRestNode("g1", "www.google.it", true)

			guard.AddService(service.GetName(), service)
			guard.RetryCb = func(guard *g.RestGuard, t *specs.RestTicket) (*specs.RestNode, error) {
				// Return always the broken node to test error with
				// multiple node
				return t.Node, nil
			}

			errAdd1 := guard.AddRestNode(service.GetName(), nodeFailed)
			errAdd2 := guard.AddRestNode(service.GetName(), node)

			t := service.GetTicket()
			defer t.Rip()
			req, errReq := guard.CreateRequest(t, "GET", "/")
			errDo := guard.Do(t)
			It("Check nil", func() {
				Expect(guard).ShouldNot(BeNil())
				Expect(req).ShouldNot(BeNil())
			})
			It("Check err", func() {
				Expect(err).Should(BeNil())
				Expect(errAdd1).Should(BeNil())
				Expect(errAdd2).Should(BeNil())
				Expect(errReq).Should(BeNil())
				Expect(errDo).ShouldNot(BeNil())
			})
			It("User Agent", func() {
				Expect(guard.GetUserAgent()).Should(Equal(
					fmt.Sprintf("RestGuard v%s", specs.RGuardVersion)))
			})
			It("Check ticket", func() {
				Expect(t).ShouldNot(BeNil())
				Expect(t.Id).ShouldNot(Equal(""))
				Expect(t.Retries).Should(Equal(2))
			})
			It("Check Response ", func() {
				Expect(t.Response).Should(BeNil())
			})

		})

		Context("Simple1 Custom Timeout", func() {

			c := specs.NewConfig()
			guard, err := g.NewRestGuard(c)

			service := specs.NewRestService("google")
			service.Retries = 1
			nodeFailed := specs.NewRestNode("gError", "127.0.0.1", true)
			node := specs.NewRestNode("g1", "www.google.it", true)

			guard.AddService(service.GetName(), service)

			errAdd1 := guard.AddRestNode(service.GetName(), nodeFailed)
			errAdd2 := guard.AddRestNode(service.GetName(), node)

			t := service.GetTicket()
			defer t.Rip()
			req, errReq := guard.CreateRequest(t, "GET", "/")
			errDo := guard.DoWithTimeout(t, 5)
			It("Check nil", func() {
				Expect(guard).ShouldNot(BeNil())
				Expect(req).ShouldNot(BeNil())
			})
			It("Check err", func() {
				Expect(err).Should(BeNil())
				Expect(errAdd1).Should(BeNil())
				Expect(errAdd2).Should(BeNil())
				Expect(errReq).Should(BeNil())
				Expect(errDo).Should(BeNil())
			})
			It("User Agent", func() {
				Expect(guard.GetUserAgent()).Should(Equal(
					fmt.Sprintf("RestGuard v%s", specs.RGuardVersion)))
			})
			It("Check ticket", func() {
				Expect(t).ShouldNot(BeNil())
				Expect(t.Id).ShouldNot(Equal(""))
				Expect(t.Retries).Should(Equal(1))
			})

			It("Check Response ", func() {
				Expect(t.Response).ShouldNot(BeNil())
				Expect(t.Response.StatusCode).Should(Equal(200))
			})

		})

		Context("Simple1 Custom Timeout with Not active node", func() {

			c := specs.NewConfig()
			guard, err := g.NewRestGuard(c)

			service := specs.NewRestService("google")
			service.Retries = 1
			nodeFailed := specs.NewRestNode("gError", "127.0.0.1", true)
			nodeFailed.Disable = true
			node := specs.NewRestNode("g1", "www.google.it", true)

			guard.AddService(service.GetName(), service)

			errAdd1 := guard.AddRestNode(service.GetName(), nodeFailed)
			errAdd2 := guard.AddRestNode(service.GetName(), node)

			t := service.GetTicket()
			defer t.Rip()
			req, errReq := guard.CreateRequest(t, "GET", "/")
			errDo := guard.DoWithTimeout(t, 5)
			It("Check nil", func() {
				Expect(guard).ShouldNot(BeNil())
				Expect(req).ShouldNot(BeNil())
			})
			It("Check err", func() {
				Expect(err).Should(BeNil())
				Expect(errAdd1).Should(BeNil())
				Expect(errAdd2).Should(BeNil())
				Expect(errReq).Should(BeNil())
				Expect(errDo).Should(BeNil())
			})
			It("User Agent", func() {
				Expect(guard.GetUserAgent()).Should(Equal(
					fmt.Sprintf("RestGuard v%s", specs.RGuardVersion)))
			})
			It("Check ticket", func() {
				Expect(t).ShouldNot(BeNil())
				Expect(t.Id).ShouldNot(Equal(""))
				Expect(t.Retries).Should(Equal(0))
			})

			It("Check Response ", func() {
				Expect(t.Response).ShouldNot(BeNil())
				Expect(t.Response.StatusCode).Should(Equal(200))
			})

		})

	})

})
