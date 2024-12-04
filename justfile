rsyncbrute:
        CGO_ENABLED=0 GOOS=linux go build -o rsyncbrute main.go

build: rsyncbrute

clean:
        rm rsyncbrute
