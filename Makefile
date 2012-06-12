# Requires a go dev setup with GOPATH set and code.google.com/p/tcgl/redis
# installed (go get it) in the GOPATH area.

# override GOARCH as appropriate
ifndef GOARCH
GOARCH := 386
endif

# change in ./debian/DEBIAN/control as well
VERSION := 1.0-1

all: queue-ok strip

queue-ok: queue-ok.go
	env GOARCH=386 go build queue-ok.go

strip:
	strip queue-ok

deb: all
	sed -i -e 's/^Version: .*$$/Version: ${VERSION}/' ./debian/DEBIAN/control
	mkdir -p ./debian/usr/bin
	cp queue-ok ./debian/usr/bin/
	dpkg-deb --build debian
	mv debian.deb queue-ok-${VERSION}.deb

clean-deb:
	rm -f ./debian/usr/bin/queue-ok
	test ! -d ./debian/usr/bin || { cd ./debian; rmdir -p usr/bin; }
	rm -f queue-ok-*.deb

clean: clean-deb
	rm -f queue-ok


