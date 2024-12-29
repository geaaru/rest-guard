/*
Copyright © 2021-2022 Funtoo Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package guard_test

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

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

			server.RouteToHandler("GET", "/body2",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/body2"),
					ghttp.RespondWithPtr(&statusCode, &returnedResp),
				),
			)

			server.RouteToHandler("POST", "/body3",
				ghttp.CombineHandlers(
					ghttp.VerifyBody([]byte("{ \"request\": \"test\" }")),
					ghttp.VerifyRequest("POST", "/body3"),
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

				fBody := func(t *specs.RestTicket) (bool, io.ReadCloser, error) {
					s := "{ \"request\": \"test\" }"
					return true,
						ioutil.NopCloser(bytes.NewReader([]byte(s))), nil
				}
				t.RequestBodyCb = fBody
				req, errReq := guard.CreateRequest(t, "GET", "/body1")
				req.Header.Add("X-Rest-Guard-Version", specs.RGuardVersion)

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

				fBody := func(t *specs.RestTicket) (bool, io.ReadCloser, error) {
					s := "{ \"request\": \"test\" }"
					return true,
						ioutil.NopCloser(bytes.NewReader([]byte(s))), nil
				}
				t.RequestBodyCb = fBody
				req, errReq := guard.CreateRequest(t, "GET", "/body1")
				req.Header.Add("X-Rest-Guard-Version", specs.RGuardVersion)

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

				req, errReq = guard.CreateRequest(t, "GET", "/body3")
				req.Method = "POST"
				req.Header.Add("X-Rest-Guard-Version", specs.RGuardVersion)

				errDo = guard.Do(t)
				Expect(errDo).Should(BeNil())
				Expect(t.Retries).Should(Equal(1))
			})

		})

		Context("Failure 2", func() {
			BeforeEach(func() {
				statusCode = 401
				returnedResp = "KO"
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
				req, errReq := guard.CreateRequest(t, "GET", "/body2")
				errDo := guard.Do(t)

				var byteValue []byte
				if t.Response != nil && t.Response.Body != nil {
					byteValue, err = ioutil.ReadAll(t.Response.Body)
				}

				Expect(guard).ShouldNot(BeNil())
				Expect(req).ShouldNot(BeNil())
				Expect(guard.GetUserAgent()).Should(Equal(
					fmt.Sprintf("RestGuard v%s", specs.RGuardVersion)))
				Expect(err).Should(BeNil())
				Expect(errAdd1).Should(BeNil())
				Expect(errReq).Should(BeNil())
				Expect(errDo).Should(Equal(errors.New("Received invalid response")))
				Expect(t.Response).ShouldNot(BeNil())
				Expect(t.Response.StatusCode).Should(Equal(401))
				Expect(string(byteValue)).Should(Equal("KO"))

				Expect(t).ShouldNot(BeNil())
				Expect(t.Id).ShouldNot(Equal(""))
				Expect(t.Retries).Should(Equal(1))
			})

		})
		Context("Failure 3 - Custom validator error", func() {
			BeforeEach(func() {
				statusCode = 401
				returnedResp = "KO"
			})

			validatorCb := func(t *specs.RestTicket) (bool, error) {
				var err error = nil

				ans := false
				if t.Response != nil &&
					(t.Response.StatusCode == 200 || t.Response.StatusCode == 201) {
					ans = true
				} else {
					err = errors.New("Custom error msg")
				}
				return ans, err
			}

			c := specs.NewConfig()
			guard, err := g.NewRestGuard(c)
			service := specs.NewRestService("local-tester")
			service.RespValidatorCb = validatorCb
			service.Retries = 1

			It("Execute call to correct valid endpoint", func() {
				fmt.Println("Client using server address ", node)
				guard.AddService(service.GetName(), service)

				errAdd1 := guard.AddRestNode(service.GetName(), node)

				t := service.GetTicket()
				defer t.Rip()
				req, errReq := guard.CreateRequest(t, "GET", "/body2")
				errDo := guard.Do(t)

				var byteValue []byte
				if t.Response != nil && t.Response.Body != nil {
					byteValue, err = ioutil.ReadAll(t.Response.Body)
				}

				Expect(guard).ShouldNot(BeNil())
				Expect(req).ShouldNot(BeNil())
				Expect(guard.GetUserAgent()).Should(Equal(
					fmt.Sprintf("RestGuard v%s", specs.RGuardVersion)))
				Expect(err).Should(BeNil())
				Expect(errAdd1).Should(BeNil())
				Expect(errReq).Should(BeNil())
				Expect(errDo).Should(Equal(errors.New("Custom error msg")))
				Expect(t.Response).ShouldNot(BeNil())
				Expect(t.Response.StatusCode).Should(Equal(401))
				Expect(string(byteValue)).Should(Equal("KO"))

				Expect(t).ShouldNot(BeNil())
				Expect(t.Id).ShouldNot(Equal(""))
				Expect(t.Retries).Should(Equal(1))
			})

		})

		// https://github.com/macaroni-os/anise-portage-converter/releases/download/v0.16.2/anise-portage-converter-v0.16.2-source.tar.gz

		Context("Download body OK", func() {
			body := "abcdefghijlmnopqrstuvwxyz"
			BeforeEach(func() {
				statusCode = 200
				returnedResp = body
			})

			validatorCb := func(t *specs.RestTicket) (bool, error) {
				var err error = nil

				ans := false
				if t.Response != nil &&
					(t.Response.StatusCode == 200 || t.Response.StatusCode == 201) {
					ans = true
				} else {
					err = errors.New("Custom error msg")
				}
				return ans, err
			}

			testFile := "/tmp/test-rest-guard-download"
			c := specs.NewConfig()
			guard, err := g.NewRestGuard(c)
			service := specs.NewRestService("local-tester")
			service.RespValidatorCb = validatorCb
			service.Retries = 1

			It("Execute call to correct valid endpoint", func() {
				fmt.Println("Client using server address ", node)
				guard.AddService(service.GetName(), service)

				errAdd1 := guard.AddRestNode(service.GetName(), node)

				t := service.GetTicket()
				defer t.Rip()
				req, errReq := guard.CreateRequest(t, "GET", "/body2")
				artefact, errDo := guard.DoDownload(t, testFile)
				defer os.Remove(testFile)

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
				Expect(t).ShouldNot(BeNil())
				Expect(t.Id).ShouldNot(Equal(""))
				Expect(t.Retries).Should(Equal(0))
				Expect(artefact).ShouldNot(BeNil())
				Expect(artefact.Path).Should(Equal(testFile))
				Expect(artefact.Size).Should(Equal(int64(len([]byte(body)))))
				Expect(artefact.Md5).Should(Equal(
					fmt.Sprintf("%x", md5.Sum([]byte(body)))))
			})

		})

	})
})
