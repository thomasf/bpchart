#define _GNU_SOURCE
#include "libomron/libomron/omron.h"
#include <stdio.h>
#include <stdlib.h>		/* atoi */
#pragma GCC diagnostic ignored "-Wformat-zero-length"


const char* m_open();
const char* m_close();

char* m_read(int bank, int idx);
int m_count(int bank);
