package role

// Resource ...
type Resource struct {
	value string
}

// String method returns a name of resouce
func (r Resource) String() string {
	return r.value
}

// Type ...
type Type struct {
	value string
}

// String method returns a name of role type
func (t Type) String() string {
	return t.value
}

var (
	// ResCluster ...
	ResCluster = Resource{"cluster"}
	// ResProject ...
	ResProject = Resource{"project"}

	// TypeRead ...
	TypeRead = Type{"read"}
	// TypeWrite ...
	TypeWrite = Type{"write"}
)

// Info ...
type Info struct {
	ID             string
	Name           string
	TargetResource Resource
	RoleType       Type
}
