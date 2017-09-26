#ifndef __COMMANDS_H_
#define __COMMANDS_H_

// Does nothing, may be used for pings or for customization
#define CMD_NOP            "NOP"

// Read whether a pin is HIGH or LOW
// Format: DR pin
// Response: DR pin value
#define CMD_DIGITAL_READ   "DR"

// Set a pin to HIGH or LOW
// Format: DW HIGH|LOW
#define CMD_DIGITAL_WRITE  "DW"

// Analog read will send a message of type AR to the server
// Format: AR pin
// Response: AR pin value
#define CMD_ANALOG_READ    "AR" // Sensors

// Write a PWM value
// Format: AW pin value
#define CMD_ANALOG_WRITE   "AW" // PWM

// Analog read will send a message of type IAR to the server
// every fixed time
// Format: IAR pin interval
// 0 to stop reading
// Response: IAR pin value
// Response is sent every interval ms
#define CMD_INTERVAL_ANALOG_READ "IAR"

// Set a servo
#define CMD_SET_SERVO "SERVO"
// Send an IR command
#define CMD_IRSEND "IRSEND"

// Send to the server a capability of a pin
// Response: CAP pin cmd name
// Example: CAP 15 DR Light bulb
#define RESP_CAP "CAP"


#define VAL_HIGH  "HIGH"
#define VAL_LOW   "LOW"


#endif /* __COMMANDS_H_ */
