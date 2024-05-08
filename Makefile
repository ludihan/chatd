.PHONY = all
all:
	go build client/cli/filesystem/filedaemon.go
	go build server/send.go

clean:
	rm filedaemon
	rm send
