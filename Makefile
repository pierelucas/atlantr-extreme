.PHONY: all install clean

TARGETNAMEWIN				:= Atlantr-Extreme-Win.exe
TARGETNAMELINUX				:= Altantr-Extreme-Unix
TARGETNAMEARM				:= Atlantr-Extreme-Arm64
TARGETNAMEMAC				:= Atlantr-Exreme-Mac

all: windows linux arm64 darwin

windows:
	GOOS=windows GOARCH=amd64 go build -o $(TARGETNAMEWIN)

linux:
	GOOS=linux GOARCH=amd64 go build -o $(TARGETNAMELINUX)

arm64:
	GOOS=linux GOARCH=arm64 go build -o $(TARGETNAMEARM)

darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(TARGETNAMEMAC)

clean:
	rm $(TARGETNAMELINUX) $(TARGETNAMEWIN)

