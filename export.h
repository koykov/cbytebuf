#ifndef CBYTEBUF_EXPORT_H
#define CBYTEBUF_EXPORT_H

/**
 * @file External function to use in CGO wrapper.
 */

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stdbool.h>
#include "types.h"

typedef void* CByteBuf;

CByteBuf cbb_new();

void cbb_write(CByteBuf *cbb_ptr, byte *data, int *data_l);

void cbb_bytes(CByteBuf *cbb_ptr, uintptr *addr, int *buf_l, int *buf_c);

void cbb_release(CByteBuf *cbb_ptr);

#ifdef __cplusplus
}
#endif

#endif //CBYTEBUF_EXPORT_H
