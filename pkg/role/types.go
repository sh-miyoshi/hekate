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
	// ResRole ...
	ResRole = Resource{"role"}
	// ResUser ...
	ResUser = Resource{"user"}
	// ResClient ...
	ResClient = Resource{"client"}
	// ResCustomRole ...
	ResCustomRole = Resource{"customrole"}

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
