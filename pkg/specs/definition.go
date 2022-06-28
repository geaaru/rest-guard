/*
	Copyright © 2021 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

import (
	"net/http"
)

type RestTicket struct {
	Id          string
	Request     *http.Request
	Response    *http.Response
	Path        string
	Retries     int
	Service     *RestService
	Node        *RestNode
	FailedNodes RestNodes
}

type RestNode struct {
	Name    string `json:"name" yaml:"name"`
	BaseUrl string `json:"base_url" yaml:"base_url"`
	Ssl     bool   `json:"ssl,omitempty" yaml:"ssl,omitempty"`
}

type RestNodes []*RestNode

type RestService struct {
	Name            string      `json:"name" yaml:"name"`
	Nodes           []*RestNode `json:"nodes" yaml:"nodes"`
	Retries         int         `json:"retries,omitempty" yaml:"retries,omitempty"`
	RetryIntervalMs int         `json:"retry_interval_ms,omitempty" yaml:"retry_interval_ms,omitempty"`

	RespValidatorCb func(t *RestTicket) bool
}

type RestGuardConfig struct {
	UserAgent string `json:"user_agent,omitempty" yaml:"user_agent,omitempty"`

	ReqsTimeout         int  `json:"reqs_timeout,omitempty" yaml:"reqs_timeout,omitempty"`
	MaxIdleConns        int  `json:"max_idle_conns,omitempty" yaml:"max_idle_conns,omitempty"`
	IdleConnTimeout     int  `json:"idle_conn_timeout,omitempty" yaml:"idle_conn_timeout,omitempty"`
	MaxConnsPerHost     int  `json:"max_conns4host,omitempty" yaml:"max_conns4host,omitempty"`
	MaxIdleConnsPerHost int  `json:"max_idleconns4host,omitempty" yaml:"max_idleconns4host,omitempty"`
	DisableCompression  bool `json:"disable_compression,omitempty" yaml:"disable_compression,omitempty"`
	InsecureSkipVerify  bool `json:"insecure_skip_verify,omitempty" yaml:"insecure_skip_verify,omitempty"`
}
