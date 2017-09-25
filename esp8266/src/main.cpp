#include <Arduino.h>
#include <ESP8266WiFi.h>
#include <ESP8266WiFiMulti.h>
#include "Messenger.h"
#include "config.h"

ESP8266WiFiMulti WiFiMulti;
Messenger m;


#define USE_SERIAL Serial

#define LED_PIN D4

void onMessage(uint8_t *payload, size_t length) {
  String cmd = (char*)payload;
  USE_SERIAL.println(cmd);
  if (!cmd.startsWith("CMD ")) return;
  cmd = cmd.substring(4);
  cmd.trim();


  if (cmd == "ON") {
    USE_SERIAL.println("On");
    digitalWrite(LED_PIN, LOW);
  }
  if (cmd == "OFF") {
    USE_SERIAL.println("Off");
    digitalWrite(LED_PIN, HIGH);
  }
}

void onConnected() {
  Serial.println("Connected");
  m.send("HELLO " BOARD_ID);
  m.send("OWNER " ADMIN_ACCOUNT);
  #ifdef BOARD_NAME
  m.send("NAME " BOARD_NAME);
  #endif
}

void onDisconnected() {
  Serial.println("Disconnected");
}



void setup() {
    pinMode(LED_PIN, OUTPUT);
    USE_SERIAL.begin(9600);

    USE_SERIAL.println();
    USE_SERIAL.println();
    USE_SERIAL.println();
    USE_SERIAL.println();

    for(uint8_t t = 3; t > 0; t--) {
        USE_SERIAL.printf("[SETUP] BOOT WAIT %d...\n", t);
        USE_SERIAL.flush();
        delay(1000);
    }

    WiFiMulti.addAP(WIFI_SSID, WIFI_PASS);

    //WiFi.disconnect();
    while(WiFiMulti.run() != WL_CONNECTED) {
        delay(100);
        USE_SERIAL.print('.');
        USE_SERIAL.flush();
    }
    USE_SERIAL.println();

    USE_SERIAL.print("Connected to ");
    USE_SERIAL.println(WiFi.SSID().c_str());

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
