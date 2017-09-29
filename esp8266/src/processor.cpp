#include "Processor.h"
#include "util.h"


void Processor::process(uint8_t *payload, size_t length) {
  String cmd = (char*)payload;
  String op = splitSpaceTrim(cmd, 0);
  if (op == CMD_DIGITAL_READ) {
    String pinStr = splitSpaceTrim(cmd, 1);
    int pin = atoi(pinStr.c_str());
    int value = digitalRead(pin);
    messenger->send(op + " " + pinStr + String(value));
    return;
  }

  if (op == CMD_DIGITAL_WRITE) {
    String pinStr = splitSpaceTrim(cmd, 1);
    int pin = atoi(pinStr.c_str());
    String val = splitSpaceTrim(cmd, 2);
    pinMode(pin, OUTPUT);
    digitalWrite(pin, val == VAL_HIGH ? HIGH : LOW);
    return;
  }

  if (op == CMD_ANALOG_READ) {
    String pinStr = splitSpaceTrim(cmd, 1);
    int pin = atoi(pinStr.c_str());
    pinMode(pin, INPUT);
    int value = analogRead(pin);
    messenger->send(op + " " + pinStr + String(value));
    return;
  }

  if (op == CMD_ANALOG_WRITE) {
    String pinStr = splitSpaceTrim(cmd, 1);
    int pin = atoi(pinStr.c_str());
    String valStr = splitSpaceTrim(cmd, 2);
    int val = atoi(valStr.c_str());
    pinMode(pin, OUTPUT);
    analogWrite(pin, val);
    return;
  }

  if (op == CMD_INTERVAL_ANALOG_READ) {
    // TODO
    return;
  }

  if (op == CMD_SET_SERVO) {
    // TODO
    return;
  }

  if (op == CMD_IRSEND) {
    // TODO
    return;
  }
}
