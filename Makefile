include $(GOROOT)/src/Make.inc

all:
	cd tea && gomake
	cd xtea && gomake

install:
	cd tea && gomake install
	cd xtea && gomake install

clean:
	cd tea && gomake clean
	cd xtea && gomake clean

test:
	cd tea && gomake test
	cd xtea && gomake test

