/*
gcc -c mylib.c -o mylib.o
clang -shared -o mylib.so mylib.o
ar rcs mylib.a mylib.o
*/

#include <stdio.h>
#include <stdlib.h>
int square(int x) {
    return x + x;
}

int cmalloc() {
    // Allocate 128 bytes
    void* ptr = malloc(128);
    if (ptr == NULL) {
        printf("malloc failed!\n");
        return 5;
    }
    printf("memory_new allocated at %p\n", ptr);

    // Optionally: initialize memory
    for (int i = 0; i < 128; i++) {
        ((char*)ptr)[i] = 0;
    }

    // Free the memory
    free(ptr);
    return 6;
}