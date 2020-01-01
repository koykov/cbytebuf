#include <stdlib.h>
#include <string.h>
#include "cbytebuf.h"

uint64_t cbb_init(int cap) {
    return (uint64_t) malloc(cap);
}

uint64_t cbb_grow_m(uint64_t addr, int cap_o, int cap_n) {
    uint64_t addr_n = (uint64_t) malloc(cap_n);
    memcpy((void*)addr_n, (void*)addr, cap_o);
    free((void*)addr);
    return addr_n;
}

uint64_t cbb_grow_r(uint64_t addr, int cap) {
    return (uint64_t) realloc((void*)addr, cap);
}

void cbb_release(uint64_t addr) {
    if ((void*)addr != NULL) {
        free((void*)addr);
    }
}
