all: clean windows linux macos out

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/auditw.exe

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/auditl

macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/auditm

clean:
	rm -f bin/auditl
	rm -f bin/auditm
	rm -f bin/auditw.exe
	
purge:
	rm -f bin/auditl
	rm -f bin/auditm
	rm -f bin/auditw.exe
	rm -f release.zip

out:
	rm -f release.zip
	rm -f release/auditl
	rm -f release/auditm
	rm -f release/auditw.exe
	cp bin/auditl release/auditl
	cp bin/auditm release/auditm
	cp bin/auditw.exe release/auditw.exe
	cp .ignore release/.ignore
	rm -f release/.DS_Store
	zip -9 -T -r release.zip release/

