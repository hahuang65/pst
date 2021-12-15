package main

import (
  "encoding/json"
  "fmt"
  "time"
)

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
  CreatedAt       time.Time   `json:"created"`
  Visibility      Visibility  `json:"visibility"`
  Sha             string      `json:"sha"`
  Files           []PasteFile `json:"files"`
  User            struct {
    CanonicalName string      `json:"canonical_name"`
    Name          string      `json:"name"`
  }
}

func createPaste(name string, visibility Visibility, contents string) {
  data := Paste{
    Visibility: visibility,
    Files: []PasteFile{{ Filename: name, Contents: contents, }},
  }

  var resp Paste

  respString := request("POST", ApiUrl + "/pastes", data)
  checkError(json.Unmarshal([]byte(respString), &resp))

  fmt.Printf("https://paste.sr.ht/%s/%s\n", resp.User.CanonicalName, resp.Sha)
}
