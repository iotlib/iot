#include "messenger.h"


void Messenger::send(String msg) {
  ws->sendTXT(msg.c_str(), msg.length());
  // Sending over an encrypted WS is blocking
  // allow the ESP to do stuff
  yield();
  this->loop();
}

void Messenger::loop() {
  ws->loop();
  yield(); // just in case
}

void Messenger::begin(
  void (*onConnected)(),
  void (*onDisconnected)(),
  void (*onMessage)(uint8_t * payload, size_t length)
) {

  if (USE_SSL) {
    ws->beginSSL(REMOTE_ADDR, REMOTE_PORT, REMOTE_URL);
  }
  else {
    ws->begin(REMOTE_ADDR, REMOTE_PORT, REMOTE_URL);
  }

  ws->onEvent([onConnected, onDisconnected, onMessage](WStype_t type, uint8_t * payload, size_t length) {
    switch(type) {
    case WStype_DISCONNECTED:
      onDisconnected();
      break;
    case WStype_CONNECTED:
      onConnected();
      break;
    case WStype_TEXT:
      onMessage(payload, length);
      break;
    }
  });
}
