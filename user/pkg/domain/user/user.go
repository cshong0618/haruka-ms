package user

type Status string

var (
	Active      = Status("ACTIVE")
	Pending     = Status("PENDING")
	Deactivated = Status("DEACTIVATED")
)

type User struct {
	ID         string
	Name       string
	UserStatus Status
}
