.PHONY: all
all: clean concat run

.PHONY: clean
clean:
	rm -f embedded/embedded{,.zip} russian-doll

.PHONY: concat
concat: embedded/embedded.zip russian-doll
	cat embedded/embedded.zip >> russian-doll
	zip -A russian-doll

.PHONY: run
run: russian-doll
	./russian-doll

russian-doll:
	go build

embedded/embedded:
	cd embedded && go build

embedded/embedded.zip: embedded/embedded
	cd embedded && zip embedded embedded
