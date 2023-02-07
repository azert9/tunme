package tunme

import (
	"fmt"
	"github.com/azert9/tunme/pkg/modules"
	"github.com/azert9/tunme/pkg/modules/builtin"
	"strings"
)

func OpenTunnel(args string) (Tunnel, error) {

	argSegments := strings.Split(args, ",")

	if len(argSegments) == 0 {
		return nil, fmt.Errorf("empty argument string")
	}
	moduleName := argSegments[0]
	argSegments = argSegments[1:]

	module, found := builtin.Modules.FindModule(moduleName)
	if !found {
		return nil, fmt.Errorf("modules not found: %s", moduleName)
	}

	var positionalArgs []string
	var namedArgs []modules.NamedArg

	for _, arg := range argSegments {

		sep := strings.IndexRune(arg, '=')

		if sep < 0 {
			positionalArgs = append(positionalArgs, arg)
		} else {
			namedArgs = append(namedArgs, modules.NamedArg{arg[:sep], arg[sep+1:]})
		}
	}

	return module.Instantiate(positionalArgs, namedArgs)
}
