# Requires a go dev setup with GOPATH set and code.google.com/p/tcgl/redis
# installed (go get it) in the GOPATH area.

# override GOARCH as appropriate
ifndef GOARCH
GOARCH := 386
endif

DEBARCH := amd64
ifeq ($(GOARCH), 386)
	DEBARCH := i386
endif

# change in ./debian/DEBIAN/control as well
VERSION := 1.3-1

all: queue-ok strip

queue-ok: queue-ok.go
	env GOARCH=${GOARCH} go build queue-ok.go

strip:
	strip queue-ok

deb-386:
	GOARCH=386 make deb

deb-amd64:
	GOARCH=amd64 make deb

deb: all
	sed -i -e 's/^Version: .*$$/Version: ${VERSION}/' ./debian/DEBIAN/control
	sed -i -e 's/^Architecture: .*$$/Architecture: ${DEBARCH}/' ./debian/DEBIAN/control
	mkdir -p ./debian/usr/bin
	cp queue-ok ./debian/usr/bin/
	dpkg-deb --build debian
	mv debian.deb queue-ok-${VERSION}_${DEBARCH}.deb
	sed -i -e 's/^Architecture: .*$$/Architecture: XXX/' ./debian/DEBIAN/control

clean-deb:
	rm -f ./debian/usr/bin/queue-ok
	test ! -d ./debian/usr/bin || { cd ./debian; rmdir -p usr/bin; }
	rm -f queue-ok-*.deb

clean: clean-deb
	rm -f queue-ok


