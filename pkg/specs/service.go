/*
	Copyright © 2021 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

import (
	"github.com/google/uuid"
)

func defaultRespCheck(t *RestTicket) bool {
	ans := false
	if t.Response != nil &&
		(t.Response.StatusCode == 200 || t.Response.StatusCode == 201) {
		ans = true
	}
	return ans
}

func NewRestService(n string) *RestService {
	return &RestService{
		Name:            n,
		Nodes:           []*RestNode{},
		Retries:         0,
		RespValidatorCb: defaultRespCheck,
		RetryIntervalMs: 10,
		Options:         make(map[string]string, 0),
	}
}

func (s *RestService) GetName() string  { return s.Name }
func (s *RestService) SetName(n string) { s.Name = n }
func (s *RestService) AddNode(n *RestNode) {
	s.Nodes = append(s.Nodes, n)
}

func (s *RestService) GetNodes() []*RestNode {
	return s.Nodes
}

func (s *RestService) HasOption(k string) bool {
	_, isPresent := s.Options[k]
	return isPresent
}

func (s *RestService) GetOption(k string) (string, error) {
	return s.Options[k]
}

func (s *RestSerivce) SetOption(k, v string) {
	s.Options[k] = v
}

func (s *RestService) GetTicket() *RestTicket {
	ans := &RestTicket{
		Id:      uuid.New().String(),
		Retries: 0,
		Node:    nil,
		Path:    "",
		Service: s,
	}

	return ans
}

func (s *RestService) Clone() (*RestService) {
	ans := &RestService{
		Id:      uuid.New().String(),
		Retries: s.Retries
		RespValidatorCb: s.defaultRespCheck,
		RetryIntervalMs: s.RetryIntervalMs,
	}

	for k, v := range s.Options {
		ans[k] = v
	}

	return ans
}
