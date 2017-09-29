package model

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
)

type Execution struct {
	DeviceId string `json:"id"`
	Cmd      string `json:"cmd"`
}

type Function struct {
	Name     string                 `json:"name"`
	Pin      int                    `json:"pin"`
	DeviceId string                 `json:"deviceid"`
	Id       string                 `json:"id"`
	Cmd      Command                `json:"cmd"`
	Data     map[string]interface{} `json:"data"`
	Owner    string                 `json:"owner"`
}
