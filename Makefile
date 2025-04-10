all:
	go build -o binaries/server cmd/server/main.go 

prep:
	cp configs/config.yaml.example configs/config.yaml

clean:
	rm -rf binaries 