BINARIES = tunme tunme-cat tunme-relay

.PHONY: all
all: ${BINARIES}

.PHONY: ${BINARIES}
${BINARIES}:
	CGO_ENABLED=0 go build ./cmd/$@

.PHONY: clean
clean:
	${RM} ${BINARIES}