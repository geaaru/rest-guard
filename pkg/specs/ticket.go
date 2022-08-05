/*
	Copyright © 2021-2022 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

import (
	"io"
	"net/http"
)

func (t *RestTicket) GetId() string               { return t.Id }
func (t *RestTicket) GetRetries() int             { return t.Retries }
func (t *RestTicket) GetService() *RestService    { return t.Service }
func (t *RestTicket) GetNode() *RestNode          { return t.Node }
func (t *RestTicket) GetRequest() *http.Request   { return t.Request }
func (t *RestTicket) GetResponse() *http.Response { return t.Response }
func (t *RestTicket) GetRequestBodyCb() func(t *RestTicket) (bool, io.ReadCloser, error) {
	return t.RequestBodyCb
}

func (t *RestTicket) Rip() {
	if t.Response != nil {
		t.Response.Body.Close()
	}
}

func (t *RestTicket) AddFail(n *RestNode) {
	ispresent := t.FailedNodes.HasNode(n)
	if !ispresent {
		t.FailedNodes = append(t.FailedNodes, n)
	}
}
