package bp

type Data struct {
	Project Project        `json:"project,omitempty"`
	Values  map[string]any `json:"values,omitempty"`
}

type Project struct {
	Name string `json:"name"`
}

type Hook struct {
	Name   string `json:"name,omitempty"`
	Script string `json:"script,omitempty"`
}

type Exclude struct {
	Pattern   string `json:"pattern,omitempty"`
	Condition string `json:"condition,omitempty"`
}

type BlueprintMeta struct {
	Input     []string  `json:"input,omitempty"`
	Exclude   []Exclude `json:"exclude,omitempty"`
	Raw       []string  `json:"raw,omitempty"`
	PreHooks  []Hook    `json:"preHooks,omitempty"`
	PostHooks []Hook    `json:"postHooks,omitempty"`
}
