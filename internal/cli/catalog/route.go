package catalog

// Route describes one user-facing backend endpoint and its resource-first CLI path.
type Route struct {
	Feature    string   `json:"feature" yaml:"feature"`
	Method     string   `json:"method" yaml:"method"`
	Endpoint   string   `json:"endpoint" yaml:"endpoint"`
	Command    []string `json:"command" yaml:"command"`
	Args       []string `json:"args,omitempty" yaml:"args,omitempty"`
	Controller string   `json:"controller" yaml:"controller"`
}
