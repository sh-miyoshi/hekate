package output

// Format ...
type Format interface {
	ToText() (string, error)
	ToJSON() (string, error)
}
