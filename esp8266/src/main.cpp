#include <Arduino.h>
#include <ESP8266WiFi.h>
#include <ESP8266WiFiMulti.h>
#include "Messenger.h"
#include "Processor.h"
#include "config.h"
#include "commands.h"

ESP8266WiFiMulti WiFiMulti;
Messenger m;
Processor proc(&m);

#define LED_PIN D4

void onMessage(uint8_t *payload, size_t length) {
  //proc.process(payload, length);
}

void onConnected() {
  Serial.println("Connected");
  m.send("HELLO " BOARD_ID);
  m.send("OWNER " ADMIN_ACCOUNT);
  #ifdef BOARD_NAME
  m.send("NAME " BOARD_NAME);
  #endif
  // TODO store this in EEPROM
  m.send(RESP_CAP " " + String(D3) + " " CMD_DIGITAL_WRITE + " Light Bulb");
}

void onDisconnected() {
  Serial.println("Disconnected");
}



void setup() {
    pinMode(LED_PIN, OUTPUT);
    Serial.begin(9600);

    Serial.println();
    Serial.println();
    Serial.println();
    Serial.println();

    for(uint8_t t = 3; t > 0; t--) {
        Serial.printf("[SETUP] BOOT WAIT %d...\n", t);
        Serial.flush();
        delay(200);
    }

    WiFiMulti.addAP(WIFI_SSID, WIFI_PASS);

    //WiFi.disconnect();
    while(WiFiMulti.run() != WL_CONNECTED) {
        delay(100);
        Serial.print('.');
        Serial.flush();
    }
    Serial.println();

    Serial.print("Connected to ");
    Serial.println(WiFi.SSID().c_str());

    m.begin(onConnected,
            onDisconnected,
            onMessage
    );

}

void loop() {
    m.loop();
    // blink for 50ms
    digitalWrite(LED_PIN, LOW);
    delay(50);
    digitalWrite(LED_PIN, HIGH);
    delay(100);
}
