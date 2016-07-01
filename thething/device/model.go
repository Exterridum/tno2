package device

import "encoding/json"

type Model struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Tags         Tags         `json:"tags"`
	CustomFields CustomFields `json:"customFields"`
	Links        Links        `json:"links"`
}

type Tags []string

type CustomFields map[string]interface{}

type Links struct {
	Product    Product    `json:"product"`
	Properties Properties `json:"properties"`
	Actions    Actions    `json:"actions"`
	Type       Type       `json:"type"`
	Help       Help       `json:"help"`
	UI         UI         `json:"ui"`
}

type Link struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

type Product struct {
	Link
}

type Properties struct {
	Link
	Resources map[string]PropertyResource `json:"resources"`
}

type PropertyResource struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Values      map[string]PropertyValue `json:"values"`
	Tags        Tags                     `json:"tags"`
}

type PropertyValue struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Unit         string       `json:"unit"`
	CustomFields CustomFields `json:"customFields"`
}

type Actions struct {
	Link
	Resources map[string]ActionResource `json:"resources"`
}

type ActionResource struct {
	Values map[string]ActionValue `json:"values"`
}

type ActionValue struct {
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type Type struct {
	Link
}

type Help struct {
	Link
}

type UI struct {
	Link
}

func (m Model) ToString() string {
	out, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return string(out)
}
