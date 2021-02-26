VERSION := $(if $(RELEASE_VERSION),$(RELEASE_VERSION),"master")

all: pre_clean ui darwin linux windows post_clean

pre_clean:
	rm -rf dist

ui:
	yarn prod
	mkdir dist
	cp .env.example dist/.env
	cp -r public dist

darwin:
	GOOS=darwin GOARCH=amd64 go build -o dist/satellite
	cd dist && zip -r satellite_$(VERSION)_darwin_amd64.zip public .env satellite
	rm -f dist/satellite

linux:
	GOOS=linux GOARCH=amd64 go build -o dist/satellite
	cd dist && zip -r satellite_$(VERSION)_linux_amd64.zip public .env satellite
	rm -f dist/satellite

windows:
	GOOS=windows GOARCH=amd64 go build -o dist/satellite.exe
	cd dist && zip -r satellite_$(VERSION)_windows_amd64.zip public .env satellite.exe
	rm -f dist/satellite.exe

post_clean:
	rm -rf dist/public dist/.env
