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

type PasteFile struct {
	Filename string `json:"filename,omitempty"`
	BlobID   string `json:"blob_id"`
	Contents string `json:"contents"`
}

type Paste struct {
	CreatedAt  time.Time   `json:"created"`
	Visibility Visibility  `json:"visibility"`
	Sha        string      `json:"sha"`
	Files      []PasteFile `json:"files"`
	User       struct {
		CanonicalName string `json:"canonical_name"`
		Name          string `json:"name"`
	}
}

func Create(name string, visibility Visibility, contents string) {
	data := Paste{
		Visibility: visibility,
		Files:      []PasteFile{{Filename: name, Contents: contents}},
	}

	var resp Paste

	respString := util.Request("POST", ApiUrl+"/pastes", data)
	util.CheckError(json.Unmarshal([]byte(respString), &resp))

	fmt.Printf("https://paste.sr.ht/%s/%s\n", resp.User.CanonicalName, resp.Sha)
}
