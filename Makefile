.PHONY: build install uninstall clean

build:
	@./build/build-release.sh

install: build
	@echo "Installing gct to /usr/local/bin..."
	@sudo install -m 755 ./build/compiled/gct /usr/local/bin/gct
	@echo "gct installed successfully!"

uninstall:
	@echo "Uninstalling gct from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/gct
	@echo "gct uninstalled successfully!"

clean:
	@echo "Cleaning up..."
	@rm -rf ./build/compiled