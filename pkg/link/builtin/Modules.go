package builtin

import (
	"tunme/pkg/link"
	"tunme/pkg/link/builtin/tcp"
	"tunme/pkg/link/builtin/udp"
)

func _makeBuiltinModuleLib() link.ModuleLibrary {

	var lib link.BasicModuleLib

	lib.Register("tcp-client", tcp.ClientModule{})
	lib.Register("tcp-server", tcp.ServerModule{})

	lib.Register("udp-client", udp.ClientModule{})
	lib.Register("udp-server", udp.ServerModule{})

	return &lib
}

var Modules = _makeBuiltinModuleLib()
