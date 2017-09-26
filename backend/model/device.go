package model

type State = int

const (
	StatePendingHello State = iota
	StatePendingOwner
	StateConnected
)

type Command = string

const (
	CmdNop                Command = "NOP"
	CmdDigitalRead                = "DR"
	CmdDigitalWrite               = "DW"
	CmdAnalogRead                 = "AR"
	CmdAnalogWrite                = "AW"
	CmdIntervalAnalogRead         = "IAR"
	CmdSetServo                   = "SERVO"
	CmdIRSend                     = "IRSEND"
	CmdCap                        = "CAP"
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

// Represents a capability of a pin
type Cap struct {
	Cmd  string `json:"cmd"`
	Pin  int    `json:"pin"`
	Name string `json:"name"`
}

type Device struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Caps     []Cap  `json:"caps"`
	State    State  `json:"state"`
	LastSeen int64  `json:"lastseen"`
}
