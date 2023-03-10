.PHONY: all
all:
	go build -tags app_cat,app_relay,app_tcp,app_tun .
