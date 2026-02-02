#ifndef _SHARED_H
#define _SHARED_H 1
#include <stdint.h>
#include <wchar.h>
#include <assert.h>

static_assert(sizeof(wchar_t) == 2, "sizeof wchar is not 2 byte; try using -fshort-wchar");

typedef struct modinfo {
	uint8_t* start;
	size_t sz;
} modinfo;


#ifdef __GNUC__
#define PACK( declaration ) declaration __attribute__((__packed__))
#elif _MSC_VER
#define PACK( declaration ) __pragma(pack(push, 1) ) declaration __pragma(pack(pop))
#endif

int get_prot(void* addr);
int change_prot(uintptr_t addr, int newProt);
modinfo get_base();
int get_rw_perms();

// linux:   48 8D ?? ?? E8 ?? ?? ?? 00 80 ?? ?? 00 0F 84
// windows: 48 8D ?? ?? ?? E8 ?? ?? ?? ?? 80 ?? ?? ?? 00 0F 84
#define ISDEBUG_PATTERN_LINUX   (mem[0x0] == 0x48 && mem[0x1] == 0x8d && mem[0x4] == 0xe8 && mem[0x8] == 0x0 && mem[0x9] == 0x80 && mem[0xc] == 0x0 && mem[0xd] == 0xf && mem[0xe] == 0x84)
#define ISDEBUG_PATTERN_WINDOWS (mem[0x0] == 0x48 && mem[0x1] == 0x8d && mem[0x5] == 0xe8 && mem[0xa] == 0x80 && mem[0xe] == 0x0 && mem[0xf] == 0xf && mem[0x10] == 0x84)

// windows: 48 8D ?? ?? ?? ?? ?? ?? E8 ?? ?? ?? ?? 80 ?? ?? ?? 00 00 00 00 0F 85 ?? ?? ?? ?? 48 
// linux:   48 8D ?? ?? ?? ?? ?? E8 ?? ?? ?? ?? 80 ?? ?? ?? ?? ?? 00 0F 85 ?? ?? ?? ?? 48
#define SETAUTH_PATTERN_WINDOWS (mem[0x0] == 0x48 && mem[0x1] == 0x8d && mem[0x8] == 0xe8 && mem[0xd] == 0x80 && mem[0x11] == 0x0 && mem[0x12] == 0x0 && mem[0x13] == 0x0 && mem[0x14] == 0x0 && mem[0x15] == 0xf && mem[0x16] == 0x85 && mem[0x1b] == 0x48)
#define SETAUTH_PATTERN_LINUX   (mem[0x0] == 0x48 && mem[0x1] == 0x8d && mem[0x7] == 0xe8 && mem[0xc] == 0x80 && mem[0x12] == 0x0 && mem[0x13] == 0xf && mem[0x14] == 0x85 && mem[0x19] == 0x48)

#ifdef __linux__
#define ISDEBUG_PATTERN_PLATFORM ISDEBUG_PATTERN_LINUX
#define SETAUTH_PATTERN_PLATFORM SETAUTH_PATTERN_LINUX
#elif _WIN32
#define ISDEBUG_PATTERN_PLATFORM ISDEBUG_PATTERN_WINDOWS
#define SETAUTH_PATTERN_PLATFORM SETAUTH_PATTERN_WINDOWS
#endif

#endif
