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
	RespCap            = "CAP"
	RespBye            = "BYE"
)

type Value = string

const (
	ValHigh Value = "HIGH"
	ValLow        = "LOW"
)

// Not to be stored in database
type Device struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Owner     string     `json:"owner"`
	Caps      []Function `json:"caps"`
	State     State      `json:"state"`
	LastSeen  int64      `json:"lastseen"`
	Confirmed bool       `json:"confirmed"`
}
