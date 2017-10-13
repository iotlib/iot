package model

import "gopkg.in/mgo.v2/bson"

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
	Id       bson.ObjectId          `json:"id" bson:"_id,omitempty"`
	Name     string                 `json:"name"`
	DeviceId string                 `json:"deviceid"`
	Owner    string                 `json:"owner"`
	Pin      int                    `json:"pin"`
	Cmd      Command                `json:"cmd"`
	Data     map[string]interface{} `json:"data"`
}
