#include <iostream>
#include "export.h"
#include "bytebuf.h"

CByteBuf cbb_new() {
    auto *cbb = new ByteBuf();
    return (void*) cbb;
}

void cbb_write(CByteBuf *cbb_ptr, byte *data, int *data_l) {
    auto *cbb = (ByteBuf*) cbb_ptr;
    cbb->write(data, *data_l);
}

void cbb_bytes(CByteBuf *cbb_ptr, uintptr *addr, int *buf_l, int *buf_c) {
    auto *cbb = (ByteBuf*) cbb_ptr;
    auto buf = cbb->get_buf();
    auto len = cbb->get_len();
    auto cap = cbb->get_cap();
    *addr = (uintptr) buf;
    *buf_l = len;
    *buf_c = cap;
}

void cbb_release(CByteBuf *cbb_ptr) {
    auto *cbb = (ByteBuf*) cbb_ptr;
    cbb->release();
}
