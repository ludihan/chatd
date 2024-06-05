.PHONY: all clean
all:
	go build -o start-server ./server && go build -o chat ./tui-client

clean:
	rm start-server
	rm chat
