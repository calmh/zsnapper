#!/bin/bash
set -euo pipefail

hash pandoc 2>/dev/null && pandoc -s -f rst -t man < doc/zsnapper.rst > man/man1/zsnapper.1

rm -rf build bin
mkdir -p build/zsnapper

for os in linux freebsd solaris ; do
	export GOARCH=amd64
	export GOOS="$os"
	go build -ldflags -w

	rm -rf build/zsnapper/bin
	mkdir build/zsnapper/bin
	mv zsnapper build/zsnapper/bin
	cp -r etc man build/zsnapper
	cp README.md build/zsnapper/README.txt
	tar -C build -zcf "zsnapper-$os-amd64.tar.gz" zsnapper
done
