.PHONY: laz bin

UNAMES := $(shell uname -s)
CUR=$(shell pwd)

ext = .exe
f = -ldflags -H=windowsgui

ifeq ($(UNAMES), Darwin)
    ext = 
    f = 
endif

all:laz res bin

laz:
	res2go -path "./src/laz" -outpath "./src" -outres false

dep:export GOPATH=$(CUR)/gopath
dep:
	go get -u github.com/xwb1989/sqlparser
	go get -u github.com/ying32/govcl
	
res:
	cd ./src/ && rsrc -ico="godbgen.ico" -o ./godbgen.syso

bin:export GOPATH=$(CUR)/gopath
bin:
	cd ./src && go build $(f) -o ../godbgen$(ext)
ifeq ($(UNAMES), Darwin)
	cp ./godbgen$(ext) ./godbgen.app/Contents/MacOS
	cp ./liblcl.dylib ./godbgen.app/Contents/MacOS
endif

debug:export GOPATH=$(CUR)/gopath
debug:
	go build -o godbgen$(ext) ./src/ && ./godbgen$(ext)

clean:
	rm -rf godbgen$(ext)
	rm -rf ./src/laz/godbgen$(ext)
	rm -rf ./src/laz/lib
	rm -rf ./src/laz/backup
	rm -rf ./godbgen.app/Contents/MacOS/godbgen$(ext)
	rm -rf ./godbgen.app/Contents/MacOS/liblcl.dylib