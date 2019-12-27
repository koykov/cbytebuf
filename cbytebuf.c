#include <stdlib.h>
#include <string.h>
#include "cbytebuf.h"

void cbb_init(error *err, uintptr *addr, const int *cap) {
    *addr = (byte*) malloc(*cap);
    *err = *addr == NULL ? ERR_BAD_ALLOC : ERR_OK;
}

void *cbb_init_np(const int cap) {
    void *addr = (byte*) malloc(cap);
    return addr;
}

void cbb_grow(error *err, uintptr *addr, const int *cap) {
    *addr = (byte*) realloc(*addr, *cap);
    *err = *addr == NULL ? ERR_BAD_ALLOC : ERR_OK;
}

void *cbb_grow_np(void *addr, const int cap) {
    void *addr_n = (byte*) realloc(addr, cap);
    return addr_n;
}

void cbb_release(error *err, uintptr *addr) {
    if (*addr != NULL) {
        free(*addr);
        *addr = NULL;
    }
    *err = ERR_OK;
}

void cbb_release_np(void *addr) {
    if (addr != NULL) {
        free(addr);
    }
}
