package types

type KGGraph struct {
	Nodes []KGNode `json:"nodes"`
	Edges []KGEdge `json:"edges"`
}

type KGNode struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Name     string `json:"name"`
	NodeType string `json:"nodeType,omitempty"`
}

type KGEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Label  string `json:"label,omitempty"`
}
