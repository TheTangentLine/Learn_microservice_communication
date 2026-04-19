MDBOOK := $(HOME)/.cargo/bin/mdbook
BOOK_DIR := mdbook

.PHONY: serve build clean

serve:
	$(MDBOOK) serve $(BOOK_DIR) --open

build:
	$(MDBOOK) build $(BOOK_DIR)

clean:
	rm -rf $(BOOK_DIR)/book
