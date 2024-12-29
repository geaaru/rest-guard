/*
Copyright © 2024-2025 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package guard

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/geaaru/rest-guard/pkg/specs"
)

type ArtefactWriter struct {
	fd     *os.File
	hasher hash.Hash
	path   string

	count int64
}

func NewArtefactWriter(file string) (*ArtefactWriter, error) {
	fd, err := os.Create(file)
	if err != nil {
		return nil, fmt.Errorf("error on create file %s: %s",
			file, err.Error())
	}

	return &ArtefactWriter{
		fd:     fd,
		path:   file,
		hasher: md5.New(),
		count:  0,
	}, nil
}

func (a *ArtefactWriter) Write(p []byte) (int, error) {
	// Write incoming bytes to file
	n, err := a.fd.Write(p)
	if err != nil {
		return n, err
	}

	// Increment byte counter
	a.count += int64(len(p))

	// Update md5
	_, err = a.hasher.Write(p)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (a *ArtefactWriter) Close() error {
	return a.fd.Close()
}

func (a *ArtefactWriter) MD5() string {
	return fmt.Sprintf("%x", a.hasher.Sum(nil))
}

func (a *ArtefactWriter) GetPath() string { return a.path }
func (a *ArtefactWriter) GetCount() int64 { return a.count }

func (g *RestGuard) DoDownload(t *specs.RestTicket, artefactPath string) (*specs.RestArtefact, error) {
	artefactWriter, err := NewArtefactWriter(artefactPath)
	if err != nil {
		return nil, err
	}
	defer artefactWriter.Close()

	err = g.doClient(g.Client, t)
	if err != nil {
		return nil, err
	}

	if t.Response == nil {
		return nil, fmt.Errorf("invalid response received")
	}

	if t.Response.StatusCode != 200 {
		return nil, fmt.Errorf("received response code %d", t.Response.StatusCode)
	}

	// Read response and write file
	_, err = io.Copy(artefactWriter, t.Response.Body)
	if err != nil {
		return nil, fmt.Errorf("error on writing file %s: %s",
			artefactPath, err.Error())
	}

	ans := &specs.RestArtefact{
		Path: artefactPath,
		Size: artefactWriter.GetCount(),
		Md5:  artefactWriter.MD5(),
	}

	return ans, nil
}
