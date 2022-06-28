/*
	Copyright © 2021-2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package guard_test

import (
	//"os"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	g "github.com/geaaru/rest-guard/pkg/guard"
	"github.com/geaaru/rest-guard/pkg/specs"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("HTTP Body Tests", func() {

	var (
		server       *ghttp.Server
		statusCode   int
		returnedResp string
		node         *specs.RestNode
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		fmt.Println("Using server address ", server.Addr())
		node = specs.NewRestNode("LocalServer", server.Addr(), false)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Testing body1", func() {

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/body1"),
					ghttp.VerifyHeader(http.Header{
						"X-Rest-Guard-Version": []string{specs.RGuardVersion},
					}),
					ghttp.VerifyBody([]byte("{ \"request\": \"test\" }")),
					ghttp.RespondWithPtr(&statusCode, &returnedResp),
				),
			)
		})

		Context("Simple1", func() {
			BeforeEach(func() {
				statusCode = 200
				returnedResp = "OK"
			})

			c := specs.NewConfig()
			guard, err := g.NewRestGuard(c)
			service := specs.NewRestService("local-tester")
			service.Retries = 1

			It("Execute call to correct valid endpoint", func() {
				fmt.Println("Client using server address ", node)
				guard.AddService(service.GetName(), service)
				errAdd1 := guard.AddRestNode(service.GetName(), node)

				t := service.GetTicket()
				defer t.Rip()
				req, errReq := guard.CreateRequest(t, "GET", "/body1")
				req.Header.Add("X-Rest-Guard-Version", specs.RGuardVersion)
				req.Body = ioutil.NopCloser(bytes.NewReader([]byte(
					"{ \"request\": \"test\" }",
				)))

				errDo := guard.Do(t)

				var byteValue []byte
				if errDo == nil && t.Response != nil {
					byteValue, err = ioutil.ReadAll(t.Response.Body)
				}

				Expect(guard).ShouldNot(BeNil())
				Expect(req).ShouldNot(BeNil())
				Expect(guard.GetUserAgent()).Should(Equal(
					fmt.Sprintf("RestGuard v%s", specs.RGuardVersion)))
				Expect(err).Should(BeNil())
				Expect(errAdd1).Should(BeNil())
				Expect(errReq).Should(BeNil())
				Expect(errDo).Should(BeNil())
				Expect(t.Response).ShouldNot(BeNil())
				Expect(t.Response.StatusCode).Should(Equal(200))
				Expect(string(byteValue)).Should(Equal("OK"))

				Expect(t).ShouldNot(BeNil())
				Expect(t.Id).ShouldNot(Equal(""))
				Expect(t.Retries).Should(Equal(0))
			})

		})

		Context("Failure 1", func() {
			BeforeEach(func() {
				statusCode = 201
				returnedResp = "OK"
			})

			c := specs.NewConfig()
			guard, err := g.NewRestGuard(c)
			service := specs.NewRestService("local-tester")
			service.Retries = 1

			It("Execute call to correct valid endpoint", func() {
				fmt.Println("Client using server address ", node)
				guard.AddService(service.GetName(), service)

				nodeFailed := specs.NewRestNode("failed", "127.0.0.1:10000", false)
				errAdd1 := guard.AddRestNode(service.GetName(), nodeFailed)
				errAdd2 := guard.AddRestNode(service.GetName(), node)

				t := service.GetTicket()
				defer t.Rip()
				req, errReq := guard.CreateRequest(t, "GET", "/body1")
				req.Header.Add("X-Rest-Guard-Version", specs.RGuardVersion)
				req.Body = ioutil.NopCloser(bytes.NewReader([]byte(
					"{ \"request\": \"test\" }",
				)))

				errDo := guard.Do(t)

				var byteValue []byte
				if errDo == nil && t.Response != nil {
					byteValue, err = ioutil.ReadAll(t.Response.Body)
				}

				Expect(guard).ShouldNot(BeNil())
				Expect(req).ShouldNot(BeNil())
				Expect(guard.GetUserAgent()).Should(Equal(
					fmt.Sprintf("RestGuard v%s", specs.RGuardVersion)))
				Expect(err).Should(BeNil())
				Expect(errAdd1).Should(BeNil())
				Expect(errAdd2).Should(BeNil())
				Expect(errReq).Should(BeNil())
				Expect(errDo).Should(BeNil())
				Expect(t.Response).ShouldNot(BeNil())
				Expect(t.Response.StatusCode).Should(Equal(201))
				Expect(string(byteValue)).Should(Equal("OK"))

				Expect(t).ShouldNot(BeNil())
				Expect(t.Id).ShouldNot(Equal(""))
				Expect(t.Retries).Should(Equal(1))
			})

		})
	})

})
