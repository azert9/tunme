package builtin

import (
	"github.com/azert9/tunme/pkg/modules"
	"github.com/azert9/tunme/pkg/modules/builtin/tcp"
	"github.com/azert9/tunme/pkg/modules/builtin/udp"
)

func _makeBuiltinModuleLib() modules.ModuleLibrary {

	var lib modules.BasicModuleLib

	lib.Register("tcp-client", tcp.ClientModule{})
	lib.Register("tcp-server", tcp.ServerModule{})

	lib.Register("udp-client", udp.ClientModule{})
	lib.Register("udp-server", udp.ServerModule{})

	return &lib
}

var Modules = _makeBuiltinModuleLib()
