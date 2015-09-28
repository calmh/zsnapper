zsnapper
========

Zsnapper automatically creates ZFS snapshots on a specified schedule, while
also removing old snapshots as required.

Builds
------

https://github.com/calmh/zsnapper/releases

Installation
------------

- Untar the distribution.
- Copy and modify `etc/zsnapper.yml.sample`.
- Run `zsnapper -c path-to-config-file -v` and observe output.

For Solaris installations, a sample SMF manifest is in etc/site-zsnapper.xml. It
assumes an installation under /opt/local/zsnapper but is easily modified for
other locations.

Documentation
-------------

https://github.com/calmh/zsnapper/blob/master/doc/zsnapper.rst

Building
--------

- Install Go
- `go get github.com/constabulary/gb/...`
- `gb build`

The bundled script `build.sh` creates distribution packages, assuming the above
requirements are met.

License
-------

MIT
