package builtin

import (
	"tunme/pkg/link"
	"tunme/pkg/link/builtin/tcp"
)

func _makeBuiltinModuleLib() link.ModuleLibrary {

	var lib link.BasicModuleLib

	lib.Register("tcp-client", tcp.ClientModule{})
	lib.Register("tcp-server", tcp.ServerModule{})

	return &lib
}

var Modules link.ModuleLibrary = _makeBuiltinModuleLib()
