#include <stdlib.h>

int main() {
    free((int*)0xffffffff);
}
