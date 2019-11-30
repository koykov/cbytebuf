#ifndef CBYTEBUF_EXPORT_H
#define CBYTEBUF_EXPORT_H

/**
 * @file External function to use in CGO wrapper.
 */

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

#define ERR_OK 0
#define ERR_BAD_ALLOC 1

typedef void* uintptr;
typedef uint8_t byte;
typedef uint error;

void cbb_init(error *err, uintptr *addr, const int *cap);

void cbb_grow(error *err, uintptr *addr, const int *cap);

void cbb_release(error *err, uintptr *addr);

#ifdef __cplusplus
}
#endif

#endif //CBYTEBUF_EXPORT_H
