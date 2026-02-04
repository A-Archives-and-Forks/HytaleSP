#!/bin/env sh
# just to script to help me not have to ctrl + r the command lol
echo This script is for quickly building and installing the flatpak
set -x
FPBUILD=
if command -v flatpak-builder; then
  FPBUILD=$(command -v flatpak-builder)
else
  FPBUILD="flatpak run org.flatpak.Builder"
fi

$FPBUILD --force-clean build --repo=repo hytaleSP.yaml --user --install
