#include <stdio.h>
#if __has_include("./out/CCallgo_export.h")
#include "./out/CCallgo_export.h"
#endif

#if __has_include("./out/_cgo_export.h")
#endif

#include "./out/_cgo_export.h"
void ACFunction() {
    printf("ACFunction()\n");
    int a = 1;
    int res = AGoFunction(&a);
    printf("AGoFunction() returned: a:%d, %d\n", a, res);
}

int main()
{
    ACFunction();
    return 0;
}

/* Through stadic/dynamic linking, we can call Go functions from C.
// dynamic linking
go build -o ./out/CCallgo_export.so -buildmode=c-shared CCallgo_export.go
gcc CCallgo.c ./out/CCallgo_export.so -o main.out 
// static linking without source code 
go build -o out/CCallgo_export.a -buildmode=c-archive CCallgo_export.go
gcc CCallgo.c ./out/CCallgo_export.a -o main.out 
// can't use cgo's outputs directly, because it lacks the go runtime.
go tool cgo -objdir=./out CCallgo_export.go
*/ 