#!/bin/sh

make -C Aurora
-ldflags="-s -w -buildid=" -trimpath

./build-flatpak.sh