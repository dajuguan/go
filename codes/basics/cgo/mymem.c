#include "mymem.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

void memory_resize(Memory* m, size_t newlen) {
    if (!m) return;

    if (newlen > m->store.cap) {
        // allocate more capacity (double strategy)
        size_t newcap = newlen * 2;
        void* newdata = realloc(m->store.data, newcap);
        if (!newdata) {
            fprintf(stderr, "OOM in memory_resize\n");
            exit(1);
        }
        m->store.data = newdata;
        m->store.cap = newcap;
    }
    m->store.len = newlen;
}

void memory_set(Memory* m, size_t offset, size_t size, const uint8_t* value) {
    if (!m) return;

    size_t end = offset + size;
    printf("resize:\n");
    if (end > m->store.len) {
        memory_resize(m, end);
    }
    memcpy((uint8_t*)m->store.data + offset, value, size);
}

void memory_set2(Memory* m, size_t offset, size_t size) {
    if (!m) return;
    uint8_t value[1] = {42};
    size_t end = offset + size;
    printf("resize:\n");
    if (end > m->store.len) {
        memory_resize(m, end);
    }
    memcpy((uint8_t*)m->store.data + offset, value, size);
}


// Allocate a new Memory struct with initial capacity
Memory* memory_new() {
    Memory* m = (Memory*)malloc(sizeof(Memory));
    m->store.data = NULL;
    m->store.len = 0;
    m->store.cap = 0;
    m->lastGasCost = 0;
    return m;
}

void memory_free(Memory* m) {
    if (!m) return;
    free(m->store.data);
    free(m);
}