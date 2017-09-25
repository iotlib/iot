package model


type State int

const (
	StatePendingHello State = iota
	StatePendingOwner
	StateConnected
)


type Device struct {
	Id string
	Name string
	Owner string
	State State

	// Todo functions
}
