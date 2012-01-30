# $GOROOT/bin must be in PATH

# override GOARCH as appropriate
ifndef GOARCH
GOARCH := amd64
endif

GOOS := linux
GOPATH := $(shell pwd)/dependencies
export GOPATH GOOS GOARCH

IPATH := ${GOPATH}/pkg/${GOOS}_${GOARCH}
LPATH := ${GOPATH}/pkg/${GOOS}_${GOARCH}

ifeq ($(GOARCH), amd64)
A:=6
else
A:=8
endif
GC:=${A}g
LD:=${A}l

all: depend resque-ok strip

resque-ok: resque-ok.A
	$(LD) -L ${LPATH} -o resque-ok resque-ok.${A}

resque-ok.A: resque-ok.go
	$(GC) -I ${IPATH} resque-ok.go

pop: populate.go
	$(GC) -I ${IPATH} populate.go
	$(LD) -L ${LPATH} -o populate populate.${A}

strip:
	strip resque-ok

clean-all: clean clean-depend clean-deb

clean:
	rm -f *.6 *.8 resque-ok populate

depend:
	mkdir -p ${GOPATH}
	goinstall tideland-rdc.googlecode.com/hg

clean-depend:
	test ! -d ${GOPATH} || rm -rf ${GOPATH}/pkg
	test ! -d ${GOPATH} || rm -rf ${GOPATH}/src
	rm -f ${GOPATH}/goinstall.log
	test ! -d ${GOPATH} || rmdir ${GOPATH}

deb: all
	mkdir -p ./debian/usr/bin
	cp resque-ok ./debian/usr/bin/
	dpkg-deb --build debian
	mv debian.deb resque-ok-0.2-1.deb

clean-deb:
	rm -f ./debian/usr/bin/resque-ok
	test ! -d ./debian/usr/bin || { cd ./debian; rmdir -p ./usr/bin; }
	rm -f resque-ok-0.2-1.deb

