.PHONY: all install clean

TARGETNAMEWIN				:= Atlantr-Extreme-Win
TARGETNAMELINUX				:= Altantr-Extreme-Unix

all: windows linux

windows:
	GOOS=windows GOARCH=amd64 go build -o $(TARGETNAMEWIN)

linux:
	GOOS=linux GOARCH=amd64 go build -o $(TARGETNAMELINUX)

clean:
	rm $(TARGETNAMELINUX) $(TARGETNAMEWIN)

