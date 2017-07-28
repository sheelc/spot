package alfred

import (
	"encoding/json"
	"fmt"
)

type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

type Mod struct {
	Valid    bool   `json:"valid,omitempty"`
	Arg      string `json:"arg,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}

type Mods struct {
	Ctrl Mod `json:"ctrl,omitempty"`
	Alt  Mod `json:"alt,omitempty"`
	Cmd  Mod `json:"cmd,omitempty"`
}

type Item struct {
	Uid          string `json:"uid,omitempty"`
	Type         string `json:"type,omitempty"`
	Title        string `json:"title,omitempty"`
	Subtitle     string `json:"subtitle,omitempty"`
	Arg          string `json:"arg,omitempty"`
	Autocomplete string `json:"autocomplete,omitempty"`
	Valid        *bool  `json:"valid,omitempty"`
	Icon         Icon   `json:"icon,omitempty"`
	Mods         Mods   `json:"mods,omitempty"`
}

type Items struct {
	Items []Item `json:"items"`
}

func PrintMenu(items Items) error {
	itemsBytes, err := json.Marshal(items)

	if err != nil {
		return err
	}

	fmt.Println(string(itemsBytes))
	return nil
}
