.PHONY: all run clean
all:
	go build -o start-server ./server

run: all
	./start-server server-config.toml

clean:
	rm start-server
