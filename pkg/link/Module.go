package link

type NamedArg struct {
	Name  string
	Value string
}

type Module interface {
	Instantiate(positionalArgs []string, namedArgs []NamedArg) (Tunnel, error)
}
