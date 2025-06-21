package utils

type Status string

const (
	StatusActive    Status = "active"
	StatusCancelled Status = "cancelled"
	StatusDeleted   Status = "deleted"
)
