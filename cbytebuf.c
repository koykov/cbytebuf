#include <stdlib.h>
#include <string.h>
#include "cbytebuf.h"

uint64_t cbb_init(int cap) {
    return (uint64_t) malloc(cap);
}

uint64_t cbb_grow(uint64_t addr, int cap_o, int cap_n) {
    // @see http://www.cplusplus.com/reference/cstdlib/realloc
    // @see http://www.cplusplus.com/reference/cstdlib/malloc
    // @see http://www.cplusplus.com/reference/cstring/memcpy
    // @see http://www.cplusplus.com/reference/cstdlib/free
    // Set of malloc()+memcpy()+free() is fastest than call of realloc().
    uint64_t addr_n = (uint64_t) malloc(cap_n);
    memcpy((void*)addr_n, (void*)addr, cap_o);
    free((void*)addr);
    return addr_n;
}

void cbb_release(uint64_t addr) {
    if ((void*)addr != NULL) {
        free((void*)addr);
    }
}
