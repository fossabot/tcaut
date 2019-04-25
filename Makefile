
MOMENT=$(shell date +'%Y%m%d-%H%M')
VERSION=$(shell git rev-parse --short HEAD)
RANDOM=$(shell awk 'BEGIN{srand();printf("%d", 65536*rand())}')

all: clean windows linux macos out

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/tcaut.exe

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/tcautl

macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/tcautm

clean:
	rm -f bin/tcautl
	rm -f bin/tcautm
	rm -f bin/tcaut.exe
	
purge:
	rm -f release/tcautl
	rm -f release/tcautm
	rm -f release/tcaut.exe
	rm -f tcaut-*.zip

out: purge
	cp bin/tcautl release/tcautl
	cp bin/tcautm release/tcautm
	cp bin/tcaut.exe release/tcaut.exe
	cp .ignore release/.ignore
	zip -9 -T -x "*.DS_Store*" -r tcaut-x-$(VERSION)-$(MOMENT).zip release/ 

out-linux: clean purge
	GOOS=linux GOARCH=amd64 go build -o bin/tcautl
	cp bin/tcautl release/tcautl
	cp .ignore release/.ignore
	zip -9 -T -x "*.DS_Store*" "*.exe" "*rgm*" "*tcautm*" -r tcaut-linux-$(VERSION)-$(MOMENT).zip release/ 







