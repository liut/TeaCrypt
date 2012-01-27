include $(GOROOT)/src/Make.inc

all:
	cd tea && gomake

install:
	cd tea && gomake install

clean:
	cd tea && gomake clean

test:
	cd tea && gomake test

