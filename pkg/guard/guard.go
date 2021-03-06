/*
	Copyright © 2021-2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package guard

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/geaaru/rest-guard/pkg/specs"
)

type RestGuard struct {
	Client    *http.Client `json:"-" yaml:"-"`
	UserAgent string       `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`

	Services map[string]*specs.RestService `json:"services" yaml:"services"`
	RetryCb  func(guard *RestGuard, t *specs.RestTicket) (*specs.RestNode, error)
}

func NewRestGuard(cfg *specs.RestGuardConfig) (*RestGuard, error) {
	idleConnTimeout, err := time.ParseDuration(fmt.Sprintf("%ds",
		cfg.IdleConnTimeout))
	if err != nil {
		return nil, err
	}

	reqsTimeout, err := time.ParseDuration(fmt.Sprintf("%ds",
		cfg.ReqsTimeout))
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConns:        cfg.MaxIdleConns,
		IdleConnTimeout:     idleConnTimeout,
		MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
		MaxConnsPerHost:     cfg.MaxConnsPerHost,
	}

	if cfg.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	ans := &RestGuard{
		UserAgent: cfg.UserAgent,
		RetryCb:   nil,
		Services:  make(map[string]*specs.RestService, 0),
	}

	ans.Client = &http.Client{
		Transport: transport,
		Timeout:   reqsTimeout,
	}

	return ans, nil
}

func (g *RestGuard) AddRestNode(srv string, n *specs.RestNode) error {
	_, ok := g.Services[srv]
	if !ok {
		return errors.New("Service " + srv + " not found")
	}

	g.Services[srv].AddNode(n)

	return nil
}

func (g *RestGuard) AddService(srv string, s *specs.RestService) {
	g.Services[srv] = s
}

func (g *RestGuard) GetUserAgent() string { return g.UserAgent }

func (g *RestGuard) GetService(srv string) (*specs.RestService, error) {
	s, ok := g.Services[srv]
	if !ok {
		return nil, errors.New("Service " + srv + " not found")
	}
	return s, nil
}

func (g *RestGuard) CreateRequest(t *specs.RestTicket, method, path string) (*http.Request, error) {

	if t.Service == nil {
		return nil, errors.New("The ticket is without service.")
	}

	if len(t.Service.Nodes) == 0 {
		return nil, errors.New("The service is without nodes.")
	}

	if t.Service.RespValidatorCb == nil {
		return nil, errors.New("Service without response validator")
	}

	if t.Request != nil {
		t.Response = nil
	}
	t.Path = path

	var rn *specs.RestNode
	if t.Node == nil {
		rn = t.Service.Nodes[t.Retries%len(t.Service.Nodes)]
		t.Node = rn
	} else {
		rn = t.Node
	}

	url := rn.GetUrlPrefix()
	if strings.HasPrefix(path, "/") {
		url += path
	} else {
		url += "/" + path
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if g.GetUserAgent() != "" {
		req.Header.Add("User-Agent", g.GetUserAgent())
	}

	t.Request = req

	return req, nil
}

func (g *RestGuard) Do(t *specs.RestTicket) error {
	var ans error = nil

	handleRetry := func() error {
		t.Retries++
		currReq := t.Request
		t.AddFail(t.Node)
		if g.RetryCb != nil {
			node, err := g.RetryCb(g, t)
			if err != nil {
				return err
			}
			t.Node = node
		} else {
			t.Node = nil
		}
		newReq, err := g.CreateRequest(t, currReq.Method, t.Path)
		if err != nil {
			return err
		}
		newReq.Header = currReq.Header
		newReq.Body = currReq.Body

		if t.FailedNodes.HasNode(t.Node) && t.Service.RetryIntervalMs > 0 {
			sleepms, err := time.ParseDuration(fmt.Sprintf(
				"%dms", t.Service.RetryIntervalMs))
			if err != nil {
				return err
			}
			time.Sleep(sleepms)
		}

		return nil
	}

	for t.Retries <= t.Service.Retries {
		resp, err := g.Client.Do(t.Request)
		if err != nil {
			ans = err
			err = handleRetry()
			if err != nil {
				return err
			}
		} else {
			ans = nil
			t.Response = resp
			valid := t.Service.RespValidatorCb(t)
			if !valid {
				err = handleRetry()
				if err != nil {
					return err
				}
			}
			break
		}
	}

	if ans != nil {
		t.Retries--
	}

	return ans
}
