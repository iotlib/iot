#include <Arduino.h>

#include <ESP8266WiFi.h>
#include <ESP8266WiFiMulti.h>

#include <WebSocketsClient.h>

#include <Hash.h>

// CONFIG

#define REMOTE_ADDR "iot.twinone.xyz"
#define REMOTE_URL "/echo"
#define REMOTE_PORT 443
#define USE_SSL true

// END CONFIG

ESP8266WiFiMulti WiFiMulti;
WebSocketsClient ws;


#define USE_SERIAL Serial

#define LED_PIN D4

void processCommand(uint8_t *payload) {
  String cmd = (char*)payload;
  USE_SERIAL.println(cmd);
  if (!cmd.startsWith("CMD ")) return;
  cmd = cmd.substring(4);

  if (cmd == "ON") digitalWrite(LED_PIN, HIGH);
  if (cmd == "OFF") digitalWrite(LED_PIN, LOW);
}

void webSocketEvent(WStype_t type, uint8_t * payload, size_t length) {
  switch(type) {
  case WStype_DISCONNECTED:
    USE_SERIAL.printf("[WSc] Disconnected!\n");
    break;

  case WStype_CONNECTED:
    USE_SERIAL.printf("[WSc] Connected to url: %s\n",  payload);
    // send message to server when Connected
    ws.sendTXT("HELLO b2345245ea");
    break;

  case WStype_TEXT:
    USE_SERIAL.printf("[WSc] get text: %s\n", payload);
    processCommand(payload);
    // send message to server
    // ws.sendTXT("message here");
    break;
  case WStype_BIN:
    //USE_SERIAL.printf("[WSc] get binary length: %u\n", length);
    //hexdump(payload, length);

    // send data to server
    // ws.sendBIN(payload, length);
    break;
  }

}

void setup() {

    pinMode(LED_PIN, OUTPUT);
    // USE_SERIAL.begin(921600);
    USE_SERIAL.begin(9600);

    //Serial.setDebugOutput(true);
    //USE_SERIAL.setDebugOutput(true);
    USE_SERIAL.println();
    USE_SERIAL.println();
    USE_SERIAL.println();
    USE_SERIAL.println();

    for(uint8_t t = 3; t > 0; t--) {
        USE_SERIAL.printf("[SETUP] BOOT WAIT %d...\n", t);
        USE_SERIAL.flush();
        delay(1000);
    }

    WiFiMulti.addAP("Orange-FC53", "95A5E324");
    WiFiMulti.addAP("R2", "e8257628246fb6");

    //WiFi.disconnect();
    while(WiFiMulti.run() != WL_CONNECTED) {
        delay(100);
        USE_SERIAL.print('.');
        USE_SERIAL.flush();
    }

    USE_SERIAL.print("Connected to ");
    USE_SERIAL.println(WiFi.SSID().c_str());

    if (USE_SSL) {
      ws.beginSSL(REMOTE_ADDR, REMOTE_PORT, REMOTE_URL);
    }
    else {
      ws.begin(REMOTE_ADDR, REMOTE_PORT, REMOTE_URL);
    }

    ws.onEvent(webSocketEvent);

}

void loop() {
    ws.loop();
    delay(200);
}
