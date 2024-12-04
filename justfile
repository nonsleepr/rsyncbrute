rsyncbrute:
        CGO_ENABLED=0 GOOS=linux go build -o rsyncbrute main.go

test: rsyncbrute
        ./rsyncbrute --host 127.0.0.1 --port 1873 --share myshare --usernames users.txt --passwords passwords.txt

build: rsyncbrute

clean:
        rm rsyncbrute
