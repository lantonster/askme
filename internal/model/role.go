package model

const (
	RoleIdUser = iota + 1
	RoleIdAdminID
	RoleIdModeratorID
)

type Role struct {
	Id          int64
	CreatedAt   int64
	UpdatedAt   int64
	Name        string
	Description string
}

func (Role) TableName() string {
	return "role"
}
