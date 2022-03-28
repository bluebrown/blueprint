package types

type Project struct {
	Name string `json:"name"`
}

type Data struct {
	Project Project
	Values  map[string]any
}

type Hook struct {
	Name   string
	Script string
}

type Exclude struct {
	Pattern   string
	Condition string
}

type BlueprintMeta struct {
	Input     []string
	Exclude   []Exclude
	Raw       []string
	PreHooks  []Hook
	PostHooks []Hook
}
