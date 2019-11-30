#include <cstdlib>
#include <cstring>
#include "cbytebuf.h"

void cbb_init(error *err, uintptr *addr, const int *cap) {
    *addr = (byte*) realloc(nullptr, *cap);
    *err = *addr == nullptr ? ERR_BAD_ALLOC : ERR_OK;
}

void cbb_grow(error *err, uintptr *addr, const int *cap) {
    *addr = (byte*) realloc(*addr, *cap);
    *err = *addr == nullptr ? ERR_BAD_ALLOC : ERR_OK;
}

void cbb_release(error *err, uintptr *addr) {
    if (*addr != nullptr) {
        free(*addr);
        *addr = nullptr;
    }
    *err = ERR_OK;
}
