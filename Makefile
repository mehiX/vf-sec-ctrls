.PHONY: default
default: binary

dist:
	mkdir -p dist

.PHONY: clean
clean:
	rm -rf dist
	rm -f cover.out

.PHONY: binary
binary:
	./script/make.sh binary

.PHONY: crossbinary-default
crossbinary-default:
	./script/make.sh crossbinary-default

.PHONY: test-unit
test-unit:
	./script/make.sh test-unit