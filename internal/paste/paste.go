package paste

import (
	"encoding/json"
	"fmt"
	"time"

	"git.sr.ht/~hwrd/pst/internal/util"
	tea "github.com/charmbracelet/bubbletea"
)

const ApiUrl = "https://paste.sr.ht/api"

type Visibility string
type ListMsg []paste

const (
	Unlisted Visibility = "unlisted"
	Public              = "public"
	Private             = "private"
)

type pasteFile struct {
	Filename string `json:"filename,omitempty"`
	BlobID   string `json:"blob_id"`
	Contents string `json:"contents"`
}

type paste struct {
	CreatedAt  time.Time   `json:"created"`
	Visibility Visibility  `json:"visibility"`
	Sha        string      `json:"sha"`
	Files      []pasteFile `json:"files"`
	User       struct {
		CanonicalName string `json:"canonical_name"`
		Name          string `json:"name"`
	}
}

type listResponse struct {
	Pastes []paste `json:"results"`
}

func Create(name string, visibility Visibility, contents string) {
	data := paste{
		Visibility: visibility,
		Files:      []pasteFile{{Filename: name, Contents: contents}},
	}

	var resp paste

	respString := util.Request("POST", ApiUrl+"/pastes", data)
	util.CheckError(json.Unmarshal([]byte(respString), &resp))

	fmt.Printf("https://paste.sr.ht/%s/%s\n", resp.User.CanonicalName, resp.Sha)
}

func List() tea.Msg {
	var resp listResponse

	respString := util.Request("GET", ApiUrl+"/pastes", nil)
	util.CheckError(json.Unmarshal([]byte(respString), &resp))

	return ListMsg(resp.Pastes)
}

func (p paste) URL() string {
	return fmt.Sprintf("https://paste.sr.ht/%s/%s", p.User.CanonicalName, p.Sha)
}

func (p paste) Delete() {
	util.Request("DELETE", ApiUrl+"/pastes/"+p.Sha, nil)
}
