package customroleapi

import (
	"net/http"
)

// AllRoleGetHandler ...
//   require role: project-read
func AllRoleGetHandler(w http.ResponseWriter, r *http.Request) {
}

// RoleCreateHandler ...
//   require role: customrole-write
func RoleCreateHandler(w http.ResponseWriter, r *http.Request) {

}

// RoleDeleteHandler ...
//   require role: role-write
func RoleDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

// RoleGetHandler ...
//   require role: role-read
func RoleGetHandler(w http.ResponseWriter, r *http.Request) {

}
