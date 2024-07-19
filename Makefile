# Define variables
BINARY_DIR := amis-ux
BINARY_BBS := $(BINARY_DIR)/amisbbs
BINARY_SETUP := $(BINARY_DIR)/setup
CONFIG_FILE := config.json
TAR_FILE := amis-ux.tar.gz
DEB_PACKAGE := amis-ux.deb
DEB_DIR := deb-package
DEB_DIR_FILES := $(DEB_DIR)/DEBIAN/control $(DEB_DIR)/usr/local/amis-ux

# Define Go build flags
GO_BUILD_FLAGS :=

# Targets
all: build copy-config

build:
	@mkdir -p $(BINARY_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BINARY_BBS) bbs/main.go
	go build $(GO_BUILD_FLAGS) -o $(BINARY_SETUP) setup/main.go

copy-config:
	cp $(CONFIG_FILE) $(BINARY_DIR)/

clean:
	rm -f $(BINARY_BBS) $(BINARY_SETUP)
	rm -f $(TAR_FILE)
	rm -rf $(DEB_DIR) $(DEB_PACKAGE)

package: all
	tar -czf $(TAR_FILE) $(BINARY_DIR)

deb: package
	@mkdir -p $(DEB_DIR)/DEBIAN
	@mkdir -p $(DEB_DIR)/usr/local/amis-ux
	cp -r $(BINARY_DIR)/* $(DEB_DIR)/usr/local/amis-ux/
	cp $(CONFIG_FILE) $(DEB_DIR)/usr/local/amis-ux/
	echo "Package: amis-ux" > $(DEB_DIR)/DEBIAN/control
	echo "Version: 1.0" >> $(DEB_DIR)/DEBIAN/control
	echo "Section: base" >> $(DEB_DIR)/DEBIAN/control
	echo "Priority: optional" >> $(DEB_DIR)/DEBIAN/control
	echo "Architecture: all" >> $(DEB_DIR)/DEBIAN/control
	echo "Maintainer: Rick Collette <megalith@root.sh>" >> $(DEB_DIR)/DEBIAN/control
	echo "Description: AMIS-UX Bulletin Board System" >> $(DEB_DIR)/DEBIAN/control
	dpkg-deb --build $(DEB_DIR) $(DEB_PACKAGE)

.PHONY: all build copy-config clean package deb
