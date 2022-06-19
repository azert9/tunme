BINARIES = tunme tunme-cat tunme-relay tunme-tun

.PHONY: all
all: ${BINARIES}

.PHONY: ${BINARIES}
${BINARIES}:
	go build -ldflags '-extldflags "-fno-PIC -static"' -buildmode pie -tags 'osusergo netgo static_build' ./cmd/$@

.PHONY: clean
clean:
	${RM} ${BINARIES}