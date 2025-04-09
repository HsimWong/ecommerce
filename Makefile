all:
	go build -o binaries/server cmd/server/main.go 

clean:
	rm -rf binaries 