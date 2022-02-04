package paste

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/hako/durafmt"
)

type ListItem struct {
	paste paste
}

func (i ListItem) filenames() []string {
	filenames := make([]string, len(i.paste.Files))
	for i, v := range i.paste.Files {
		filenames[i] = v.Filename
	}

	return filenames
}

func (i ListItem) Title() string {
	s := ""

	switch i.paste.Visibility {
	case Public:
		s += " "
	case Private:
		s += " "
	case Unlisted:
		s += " "
	}

	return s + i.Name()
}

func (i ListItem) Name() string {
	return i.paste.Sha
}

func (i ListItem) Description() string {
	return fmt.Sprintf("Files: %s (Created %s ago)",
		strings.Join(i.filenames(), ", "),
		durafmt.ParseShort(time.Since(i.paste.CreatedAt)))
}

func (i ListItem) FilterValue() string {
	// Allow filtering by the SHA or any file names in the paste
	return i.paste.Sha + ", " + strings.Join(i.filenames(), ", ")
}

func ListItems(ps []paste) []list.Item {
	list_items := []list.Item{}

	for _, p := range ps {
		list_items = append(list_items, ListItem{paste: p})
	}

	return list_items
}
