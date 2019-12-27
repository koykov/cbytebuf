#ifndef CBYTEBUF_EXPORT_H
#define CBYTEBUF_EXPORT_H

/**
 * @file
 * Memory manipulation functions and types.
 */

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

/**
 * Error codes macros.
 */
// No error.
#define ERR_OK 0
// Failed alloc/realloc call.
#define ERR_BAD_ALLOC 1

/**
 * Type declarations to correspond Go types.
 */
typedef void* uintptr;
typedef uint8_t byte;
typedef unsigned int error;

/**
 * Initialize byte array with given capacity.
 *
 * @param err   Error code. Output param.
 * @param addr  Address of allocated array. Output param.
 * @param cap   Capacity of the array.
 */
void cbb_init(error *err, uintptr *addr, const int *cap);
void *cbb_init_np(const int cap);

/**
 * Change capacity of the array.
 *
 * This function allows to reduce array's capacity as well.
 * @see http://www.cplusplus.com/reference/cstdlib/realloc
 * @param err   Error code. Output param.
 * @param addr  Address of reallocated array. Output param.
 * @param cap   New capacity if the array. May be less than old capacity.
 */
void cbb_grow(error *err, uintptr *addr, const int *cap);
void *cbb_grow_np(void *addr, const int cap);

/**
 * Release buffer memory.
 *
 * @param err   Error code. Output param.
 * @param addr  Address of the array to release. Must become NULL after release.
 */
void cbb_release(error *err, uintptr *addr);
void cbb_release_np(void *addr);

#ifdef __cplusplus
}
#endif

#endif //CBYTEBUF_EXPORT_H
