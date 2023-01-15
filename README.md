# tunme

`tunme` is a simple program for creating tunnels. It is self-contained and statically linked.
By using a modular approach, `tunme` allows to encapsulate different applications into different communication
protocols.

## Examples

Create a TCP tunnel between two machines:
```sh
sudo tunme tun tcp-server,:80 --address=10.0.0.1/24
```
```sh
sudo tunme tun tcp-client,example.com:80 --address=10.0.0.2/24
```

Exchange data between two machines behind a NAT, using a third machine as a relay:
```sh
tunme relay tcp-server,:8080 tcp-server,:8081
```
```sh
tunme cat tcp-client,example.com:8080
```
```sh
tunme cat tcp-client,example.com:8081
```

## TODO

* Middlewares (for e.g. cryptography)
* SOCKS
* Pre-built binaries
* Find a way of allowing multiples clients on the same TCP port (would be useful for relays)