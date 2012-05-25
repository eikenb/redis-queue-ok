# Requires a go dev setup with GOPATH set and code.google.com/p/tcgl/redis
# installed (go get it) in the GOPATH area.

# override GOARCH as appropriate
ifndef GOARCH
GOARCH := 386
endif

all: resque-ok strip

resque-ok: resque-ok.go
	env GOARCH=386 go build resque-ok.go

pop: populate.go
	go build populate.go

strip:
	strip resque-ok

deb: all
	mkdir -p ./debian/usr/bin
	cp resque-ok ./debian/usr/bin/
	dpkg-deb --build debian
	mv debian.deb resque-ok-0.3-1.deb

clean-deb:
	rm -f ./debian/usr/bin/resque-ok
	test ! -d ./debian/usr/bin || { cd ./debian; rmdir -p usr/bin; }
	rm -f resque-ok-0.3-1.deb

clean: clean-deb
	rm -f resque-ok populate


