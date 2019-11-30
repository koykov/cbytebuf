#ifndef CBYTEBUF_EXPORT_H
#define CBYTEBUF_EXPORT_H

/**
 * @file External function to use in CGO wrapper.
 */

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

typedef void* uintptr;
typedef uint8_t byte;

void cbb_init(uintptr *addr, const int *cap);

void cbb_grow(uintptr *addr, const int *cap);

void cbb_release(uintptr *addr);

#ifdef __cplusplus
}
#endif

#endif //CBYTEBUF_EXPORT_H
