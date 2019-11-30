#include <cstdlib>
#include <cstring>
#include <cmath>
#include "bytebuf.h"

ByteBuf::ByteBuf() {
    this->buf = nullptr;
    this->len = 0;
    this->cap = 0;
}

ByteBuf::~ByteBuf() {
    this->release();
}

void ByteBuf::write(const byte *data, const int data_l) {
    if (this->buf == nullptr) {
        int cap_n = data_l * 2;
        this->buf = (byte*) realloc(this->buf, cap_n);
        this->cap = cap_n;
        memcpy(this->buf, data, data_l);
        this->len = data_l;
    } else {
        if (this->len + data_l > this->cap) {
            int cap_n = std::max(this->cap, this->len + data_l) * 2;
            this->buf = (byte*) realloc(this->buf, cap_n);
            this->cap = cap_n;
        }
        memcpy(&this->buf[this->len], data, data_l);
        this->len += data_l;
    }
}

byte *ByteBuf::get_buf() {
    return this->buf;
}

int ByteBuf::get_len() {
    return this->len;
}

int ByteBuf::get_cap() {
    return this->cap;
}

void ByteBuf::release() {
    if (this->buf != nullptr) {
        free(this->buf);
    }
}
