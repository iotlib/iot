#include <WebSocketsClient.h>
#include <Hash.h>

#include "config.h"

class Messenger {
private:
  WebSocketsClient *ws;


public:
  Messenger() {
    ws = new WebSocketsClient();
  }

  ~Messenger() {
    delete ws;
  }

  void send(String message);
  void begin(void (*onConnected)(), void (*onDisconnected)(), void (*onMessage)(uint8_t * payload, size_t length));
  void loop();

};
