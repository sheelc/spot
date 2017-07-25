package alfred

import (
	"encoding/json"
	"fmt"
)

type AlfredIcon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

type AlfredItem struct {
	Uid          string     `json:"uid,omitempty"`
	Type         string     `json:"type,omitempty"`
	Title        string     `json:"title,omitempty"`
	Subtitle     string     `json:"subtitle,omitempty"`
	Arg          string     `json:"arg,omitempty"`
	Autocomplete string     `json:"autocomplete,omitempty"`
	Valid        *bool      `json:"valid,omitempty"`
	Icon         AlfredIcon `json:"icon,omitempty"`
}

type AlfredItems struct {
	Items []AlfredItem `json:"items"`
}

func PrintMenu(items AlfredItems) error {
	itemsBytes, err := json.Marshal(items)

	if err != nil {
		return err
	}

	fmt.Println(string(itemsBytes))
	return nil
}
