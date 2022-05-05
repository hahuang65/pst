package paste

import (
	"encoding/json"
	"fmt"
	"time"

	"git.sr.ht/~hwrd/pst/internal/util"
)

const ApiUrl = "https://paste.sr.ht/api"

type Visibility string

const (
	Unlisted Visibility = "unlisted"
	Public              = "public"
	Private             = "private"
)

// This is the structure for paste creation
type pasteFile struct {
	Filename string `json:"filename,omitempty"`
	BlobID   string `json:"blob_id"`
	Contents string `json:"contents"`
}

// This is the structure for paste loading from the API
type PasteBlob struct {
	Filename string // No JSON field, have to set this manually
	ID       string `json:"sha"`
	Contents string `json:"contents"`
}

type Paste struct {
	CreatedAt  time.Time   `json:"created"`
	Visibility Visibility  `json:"visibility"`
	Sha        string      `json:"sha"`
	Files      []pasteFile `json:"files"`
	User       struct {
		CanonicalName string `json:"canonical_name"`
		Name          string `json:"name"`
	}
}

type pasteList struct {
	Pastes []Paste `json:"results"`
}

func Create(name string, visibility Visibility, contents string) {
	data := Paste{
		Visibility: visibility,
		Files:      []pasteFile{{Filename: name, Contents: contents}},
	}

	var resp Paste

	respString := util.Request("POST", ApiUrl+"/pastes", data)
	util.CheckError(json.Unmarshal([]byte(respString), &resp))

	fmt.Printf("https://paste.sr.ht/%s/%s\n", resp.User.CanonicalName, resp.Sha)
}

func (p Paste) URL() string {
	return fmt.Sprintf("https://paste.sr.ht/%s/%s", p.User.CanonicalName, p.Sha)
}

func (p Paste) Delete() {
	util.Request("DELETE", ApiUrl+"/pastes/"+p.Sha, nil)
}

func List() []Paste {
	var resp pasteList

	respString := util.Request("GET", ApiUrl+"/pastes", nil)
	util.CheckError(json.Unmarshal([]byte(respString), &resp))

	return resp.Pastes
}

func (p Paste) LoadFiles() []PasteBlob {
	files := make([]PasteBlob, len(p.Files))

	for i, f := range p.Files {
		var blob PasteBlob
		jsonString := util.Request("GET", ApiUrl+"/blobs/"+f.BlobID, nil)
		util.CheckError(json.Unmarshal([]byte(jsonString), &blob))

		blob.Filename = f.Filename
		files[i] = blob
	}

	return files
}
