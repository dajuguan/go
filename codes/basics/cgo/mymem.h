// mymem.h
/*
clang -c -fPIC mymem.c -o mymem.o
clang -shared -o libmymem.so mymem.o
ar rcs mymem.a mymem.o
*/
#include <stdint.h>
#include <stddef.h>

typedef struct {
    uint8_t   *data;
    size_t  len;
    size_t  cap;
} SliceHeader;

typedef struct {
    SliceHeader store;
    uint64_t    lastGasCost;
} Memory;

void memory_resize(Memory* m, size_t newlen);
void memory_set(Memory* m, size_t offset, size_t size, const uint8_t* value);
void memory_set2(Memory* m, size_t offset, size_t size);
Memory* memory_new();
void memory_free(Memory* m);