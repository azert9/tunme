package modules

import "errors"

var ErrTunnelClosed = errors.New("tunnel closed")
var ErrStreamRejected = errors.New("stream rejected")
