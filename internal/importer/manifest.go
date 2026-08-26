package importer

import (
	"crypto/sha256"
	"fmt"
	"sort"
)

type Manifest struct {
	ImportID string   `json:"import_id"`
	Files    []string `json:"files"`
	Digest   string   `json:"digest"`
	Records  int      `json:"records"`
}

func BuildManifest(id string, files []string, records int) Manifest {
	ordered := append([]string(nil), files...)
	sort.Strings(ordered)
	h := sha256.New()
	for _, file := range ordered {
		h.Write([]byte(file))
		h.Write([]byte{0})
	}
	return Manifest{ImportID: id, Files: ordered, Digest: fmt.Sprintf("sha256:%x", h.Sum(nil)), Records: records}
}
