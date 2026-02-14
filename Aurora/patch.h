#ifndef _PATCH_H
#define _PATCH_H

// stop telling hytale what we're doing ...
#define DISABLE_TELEMETRY 1

void changeServers();
int needsArgumentModify(const char* program);
int modifyArgument(const char* program, char* arg);


#endif