package models

type Role struct {
	Privs  Privileges
	Color  string
	RoleID string
	Name   string
}