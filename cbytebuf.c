#include <stdlib.h>
#include <string.h>
#include "cbytebuf.h"

void cbb_init(error *err, uintptr *addr, const int *cap) {
    *addr = (byte*) malloc(*cap);
    *err = *addr == NULL ? ERR_BAD_ALLOC : ERR_OK;
}

uint64_t cbb_init_np(int cap) {
    return (uint64_t) malloc(cap);
}

void cbb_grow(error *err, uintptr *addr, const int *cap) {
    *addr = (byte*) realloc(*addr, *cap);
    *err = *addr == NULL ? ERR_BAD_ALLOC : ERR_OK;
}

uint64_t cbb_grow_np(uint64_t addr, int cap) {
    return (uint64_t) realloc((void*)addr, cap);
}

uint64_t cbb_grow_np1(uint64_t addr, int cap_o, int cap_n) {
    uint64_t addr_n = (uint64_t) malloc(cap_n);
    memcpy((void*)addr_n, (void*)addr, cap_o);
    free((void*)addr);
    return addr_n;
}

void cbb_release(error *err, uintptr *addr) {
    if (*addr != NULL) {
        free(*addr);
        *addr = NULL;
    }
    *err = ERR_OK;
}

void cbb_release_np(uint64_t addr) {
    if ((void*)addr != NULL) {
        free((void*)addr);
    }
}
