#ifndef CBYTEBUF_EXPORT_H
#define CBYTEBUF_EXPORT_H

/**
 * @file
 * Memory manipulation functions.
 */

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

/**
 * Initialize byte array with given capacity.
 *
 * @param cap Capacity of the array.
 * @return uint64_t
 */
uint64_t cbb_init(int cap);

/**
 * Change capacity of the array using malloc().
 *
 * This function allows to reduce array's capacity as well.
 * @see http://www.cplusplus.com/reference/cstdlib/malloc
 * @see http://www.cplusplus.com/reference/cstring/memcpy
 * @see http://www.cplusplus.com/reference/cstdlib/free
 * @param addr  Address of reallocated array.
 * @param cap_o Old capacity.
 * @param cap_n New capacity of the array. May be less than old capacity.
 * @return uint64 address of first item of array in virtual memory.
 */
uint64_t cbb_grow_m(uint64_t addr, int cap_o, int cap_n);

/**
 * Change capacity of the array using realloc().
 *
 * This function allows to reduce array's capacity as well.
 * @see http://www.cplusplus.com/reference/cstdlib/realloc
 * @param addr Address of reallocated array.
 * @param cap  New capacity of the array. May be less than old capacity.
 * @return uint64 address of first item of array in virtual memory.
 */
uint64_t cbb_grow_r(uint64_t addr, int cap);

/**
 * Release buffer memory.
 *
 * @param addr  Address of the array to release.
 */
void cbb_release(uint64_t addr);

#ifdef __cplusplus
}
#endif

#endif //CBYTEBUF_EXPORT_H
