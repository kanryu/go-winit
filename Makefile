RUST_TARGET ?= x86_64-pc-windows-gnu
LIBWINIT_MANIFEST := third_party/libwinit/Cargo.toml
CARGO_TARGET_DIR ?= third_party/libwinit/target
BUILD_ROOT ?= build/libwinit/$(RUST_TARGET)

.PHONY: libwinit-debug libwinit-release libwinit-static

libwinit-debug:
	cargo build --manifest-path $(LIBWINIT_MANIFEST) --package libwinit --target $(RUST_TARGET) --target-dir $(CARGO_TARGET_DIR)
	mkdir -p $(BUILD_ROOT)/debug
	cp $(CARGO_TARGET_DIR)/$(RUST_TARGET)/debug/libWinit.a $(BUILD_ROOT)/debug/libwinit.a

libwinit-release:
	cargo build --manifest-path $(LIBWINIT_MANIFEST) --package libwinit --release --target $(RUST_TARGET) --target-dir $(CARGO_TARGET_DIR)
	mkdir -p $(BUILD_ROOT)/release
	cp $(CARGO_TARGET_DIR)/$(RUST_TARGET)/release/libWinit.a $(BUILD_ROOT)/release/libwinit.a

libwinit-static: libwinit-debug libwinit-release
