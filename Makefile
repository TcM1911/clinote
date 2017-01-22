VERSION_FILE=cmd/root.go
VERSION=$(shell grep "const version string" $(VERSION_FILE) | cut -d '"' -f 2)

RELEASE_FILES=CHANGELOG.md LICENSE README.md TODO.md

build:
	go build -v

build_dev:
	go build -v -tags=dev

clean:
	rm -f clinote
	rm -f *.tar.gz *.tar.gz.sha256sum
	rm -rf target

build_386:
	mkdir -p target/clinote-$(VERSION)-i386
	GOOS=linux GOARCH=386 go build -v -a -o target/clinote-$(VERSION)-i386/clinote
	cp $(RELEASE_FILES) target/clinote-$(VERSION)-i386/

build_amd64:
	mkdir -p target/clinote-$(VERSION)-amd64
	GOOS=linux GOARCH=amd64 go build -v -a -o target/clinote-$(VERSION)-amd64/clinote
	cp $(RELEASE_FILES) target/clinote-$(VERSION)-amd64/

release_386:
	tar cfvz clinote-$(VERSION)-i386.tar.gz -C target clinote-$(VERSION)-i386
	sha256sum clinote-$(VERSION)-i386.tar.gz > clinote-$(VERSION)-i386.tar.gz.sha256sum

release_amd64:
	tar cfvz clinote-$(VERSION)-amd64.tar.gz -C target clinote-$(VERSION)-amd64
	sha256sum clinote-$(VERSION)-amd64.tar.gz > clinote-$(VERSION)-amd64.tar.gz.sha256sum

build_all: build_386 build_amd64

release_all: release_386 release_amd64

prep_release:
	sed -i 's/[0-9]\+.[0-9]\+.[0-9]\+-SNAPSHOT/$(VERSION)/' $(VERSION_FILE)
	git add $(VERSION_FILE) && git commit -m "Set release version"
	govendor list && govendor init && govendor add +external
	git add vendor/ && git commit -m "Vendor deps"

next_dev_cycle:
	sed -i 's/[0-9]\+.[0-9]\+.[0-9]\+/$(NEXT_VERSION)-SNAPSHOT/' $(VERSION_FILE)
	git add $(VERSION_FILE) && git commit -m "Set next dev cycle version"
	rm -rf vendor/ && git add vendor/ && git commit -m "Remove vendors"


