package tunme

import (
	"fmt"
	"github.com/azert9/tunme/pkg/link"
	"github.com/azert9/tunme/pkg/link/builtin"
	"strings"
)

func OpenTunnel(args string) (link.Tunnel, error) {

	argSegments := strings.Split(args, ",")

	if len(argSegments) == 0 {
		return nil, fmt.Errorf("empty argument string")
	}
	moduleName := argSegments[0]
	argSegments = argSegments[1:]

	module, found := builtin.Modules.FindModule(moduleName)
	if !found {
		return nil, fmt.Errorf("module not found: %s", moduleName)
	}

	var positionalArgs []string
	var namedArgs []link.NamedArg

	for _, arg := range argSegments {

		sep := strings.IndexRune(arg, '=')

		if sep < 0 {
			positionalArgs = append(positionalArgs, arg)
		} else {
			namedArgs = append(namedArgs, link.NamedArg{arg[:sep], arg[sep+1:]})
		}
	}

	return module.Instantiate(positionalArgs, namedArgs)
}
