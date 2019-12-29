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
 * Change capacity of the array.
 *
 * This function allows to reduce array's capacity as well.
 * @param addr  Address of reallocated array.
 * @param cap_o Old capacity.
 * @param cap_n New capacity of the array. May be less than old capacity.
 */
uint64_t cbb_grow(uint64_t addr, int cap_o, int cap_n);

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
