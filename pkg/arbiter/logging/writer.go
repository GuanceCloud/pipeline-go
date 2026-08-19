package logging

// Writer publishes log data for the workspace bound to an Arbiter execution.
// Implementations own destination selection and authentication; scripts only
// provide the data and an optional index name.
type Writer interface {
	Push(index string, data any) (map[string]any, error)
}
