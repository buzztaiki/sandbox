#include <stdio.h>
#include "zlib.h"

extern void zuser() {
    printf("%s\n", zlibVersion());
}
