## Design

* All goroutines must use the WaitGroup of the tunnel.
* All channel operations must be done through a select with the close channel.
