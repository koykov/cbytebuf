#ifndef CBYTEBUF_BYTEBUF_H
#define CBYTEBUF_BYTEBUF_H

#include <string>
#include "types.h"

class ByteBuf {
public:
    explicit ByteBuf();
    ~ByteBuf();

    void write(const byte *data, int data_l);
    byte *get_buf();
    int get_len();
    int get_cap();
    void release();

private:
    byte *buf;
    int len;
    int cap;
};

#endif //CBYTEBUF_BYTEBUF_H
