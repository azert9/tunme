BINARIES = tunme tunme-cat tunme-relay tunme-test tunme-tun

.PHONY: all
all: ${BINARIES}

.PHONY: ${BINARIES}
${BINARIES}:
	go build -ldflags '-extldflags "-fno-PIC -static"' -buildmode pie -tags 'osusergo netgo static_build' ./cmd/$@

.PHONY: clean
clean:
	${RM} ${BINARIES}
	${RM} coverage.*

.PHONY: test
test:
	go test -cover -coverprofile=coverage.cov -coverpkg=./... ./...

.PHONY: coverage
coverage: test
	go tool cover -o coverage.html -html ./coverage.cov