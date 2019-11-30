#include <cstdlib>
#include <cstring>
#include "cbytebuf.h"

void cbb_init(uintptr *addr, const int *cap) {
    *addr = realloc(nullptr, *cap);
}

void cbb_grow(uintptr *addr, const int *cap) {
    *addr = (byte*) realloc(*addr, *cap);
}

void cbb_release(uintptr *addr) {
    if (*addr != nullptr) {
        free(*addr);
        *addr = 0;
    }
}
