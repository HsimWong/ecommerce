all:
	go build -o binaries/server cmd/server/main.go 

prep:
	cp configs/config.yaml.example configs/config.yaml

test:
	go build -o binaries/test cmd/test/test.go

clean:
	rm -rf binaries 