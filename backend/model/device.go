package model

type State = int

const (
	StatePendingHello State = iota
	StatePendingOwner
	StateConnected
)

type Response = string

const (
	RespHello Response = "HELLO"
	RespOwner          = "OWNER"
	RespName           = "NAME"
	RespBye            = "BYE"
)

type Value = string

const (
	ValHigh Value = "HIGH"
	ValLow        = "LOW"
)

// Not to be stored in database
type Device struct {
	Id    string `json:"id"`
	Owner string `json:"owner"`

	Name      string     `json:"name"`
	Confirmed bool       `json:"confirmed"`
	Functions []Function `json:"functions"`
	State     State      `json:"state"`
	LastSeen  int64      `json:"lastseen"`
}
