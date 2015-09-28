#!/bin/bash
set -euo pipefail

hash pandoc 2>/dev/null && pandoc -s -f rst -t man < doc/zsnapper.rst > man/man1/zsnapper.1

rm -rf build bin
mkdir -p build/zsnapper

for os in linux freebsd solaris ; do
	export GOARCH=amd64
	export GOOS="$os"
	gb build -ldflags -w
	[ -f "bin/zsnapper-$os-amd64" ] && mv "bin/zsnapper-$os-amd64" "bin/zsnapper"

	rm -rf build/zsnapper/bin
	mv bin build/zsnapper/bin
	cp -r etc man build/zsnapper
	cp README.md build/zsnapper/README.txt
	tar -C build -zcvf "zsnapper-$os-amd64.tar.gz" zsnapper
done
