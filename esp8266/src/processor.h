#ifndef __PROCESSOR_H__
#define __PROCESSOR_H__

#include <Arduino.h>
#include "messenger.h"

#include "commands.h"


class Processor {
private:
  Messenger *messenger;

public:
  Processor(Messenger *m) {
    messenger = m;
  }
  void process(uint8_t *payload, size_t length);

};

#endif /* __PROCESSOR_H__ */
