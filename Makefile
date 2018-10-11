VERSION_FILE=cmd/root.go
VERSION=$(shell grep "const version string" $(VERSION_FILE) | cut -d '"' -f 2)

DIST_FOLDER=dist
RELEASE_FILES=CHANGELOG.md LICENSE README.md TODO.md

build:
	go build -v -o clinote ./cmd/...

clean:
	rm -f clinote
	rm -f *.tar.gz *.tar.gz.sha256sum
	rm -rf target
	rm -f *.cov
	rm -f coverage_*
	rm -rf $(DIST_FOLDER)

test_evernote:
	go test -v ./evernote/...

test:
	go test -v ./...

coverage_evernote:
	go test -coverprofile=coverage_evernote ./evernote

coverage: coverage_evernote
	for i in $(shell ls coverage_*); do \
	cat $$i | tail -n +2 >> profile.cov; \
	done

build_386:
	mkdir -p target/clinote-$(VERSION)-i386
	GOOS=linux GOARCH=386 go build -v -a -o target/clinote-$(VERSION)-i386/clinote ./cmd
	cp $(RELEASE_FILES) target/clinote-$(VERSION)-i386/

build_amd64:
	mkdir -p target/clinote-$(VERSION)-amd64
	GOOS=linux GOARCH=amd64 go build -v -a -o target/clinote-$(VERSION)-amd64/clinote ./cmd
	cp $(RELEASE_FILES) target/clinote-$(VERSION)-amd64/

release_386:
	mkdir -p $(DIST_FOLDER)
	tar cfvz $(DIST_FOLDER)/clinote-$(VERSION)-i386.tar.gz -C target clinote-$(VERSION)-i386
	cd dist && sha256sum clinote-$(VERSION)-i386.tar.gz > clinote-$(VERSION)-i386.tar.gz.sha256sum

release_amd64:
	mkdir -p $(DIST_FOLDER)
	tar cfvz $(DIST_FOLDER)/clinote-$(VERSION)-amd64.tar.gz -C target clinote-$(VERSION)-amd64
	cd dist && sha256sum clinote-$(VERSION)-amd64.tar.gz > clinote-$(VERSION)-amd64.tar.gz.sha256sum

release_macos:
	mkdir -p $(DIST_FOLDER)
	tar cfvz $(DIST_FOLDER)/clinote-$(VERSION)-macos.tar.gz -C target clinote-$(VERSION)-macos
	cd dist && shasum -a 256 clinote-$(VERSION)-macos.tar.gz > clinote-$(VERSION)-macos.tar.gz.sha256sum

build_macos:
	mkdir -p target/clinote-$(VERSION)-macos
	GOOS=darwin GOARCH=amd64 go build -v -a -o target/clinote-$(VERSION)-macos/clinote ./cmd
	cp $(RELEASE_FILES) target/clinote-$(VERSION)-macos/

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

.PHONY: no_targets__ list
no_targets__:
list:
	sh -c "$(MAKE) -p no_targets__ | awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ {split(\$$1,A,/ /);for(i in A)print A[i]}' | grep -v '__\$$' | sort"

