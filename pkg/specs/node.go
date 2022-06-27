/*
	Copyright © 2021 Funtoo Macaroni OS Linux
	See AUTHORS and LICENSE for the license details and contributors.
*/
package specs

import (
	"strings"
)

func NewRestNode(name, burl string, ssl bool) *RestNode {
	if strings.HasSuffix(burl, "/") {
		burl = burl[0 : len(burl)-1]
	}
	return &RestNode{
		Name:    name,
		BaseUrl: burl,
		Ssl:     ssl,
	}
}

func (n *RestNode) GetUrlPrefix() string {
	ans := ""
	if n.Ssl {
		ans = "https://"
	} else {
		ans = "http://"
	}

	ans += n.BaseUrl

	return ans
}
