#!/bin/bash

echo "Generating binary for linuxAMD64..."
env GOOS=linux GOARCH=amd64 go build
zip linuxAMD64.zip stanford-dl
echo "Done."

echo "Generating binary for linuxARM64..."
env GOOS=linux GOARCH=arm64 go build
zip linuxARM64.zip stanford-dl
echo "Done."

echo "Generating binary for linuxARM32..."
env GOOS=linux GOARCH=arm go build
zip linuxARM32.zip stanford-dl
echo "Done."

echo "Generating binary for linux386..."
env GOOS=linux GOARCH=386 go build
zip linux386.zip stanford-dl
echo "Done."

echo "Generating binary for windows64..."
env GOOS=windows GOARCH=amd64 go build
zip windows64.zip stanford-dl.exe
echo "Done."

echo "Generating binary for windows32..."
env GOOS=windows GOARCH=386 go build
zip windows32.zip stanford-dl.exe
echo "Done."

# Generate system binary
rm -f stanford-dl.exe stanford-dl
go build

mkdir -p binaries
mv *.zip binaries
