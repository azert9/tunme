package link

type ModuleLibrary interface {
	FindModule(name string) (Module, bool)
}

type BasicModuleLib struct {
	Modules map[string]Module
}

func (l *BasicModuleLib) Register(name string, module Module) {

	if l.Modules == nil {
		l.Modules = make(map[string]Module)
	}

	l.Modules[name] = module
}

func (l *BasicModuleLib) FindModule(name string) (Module, bool) {

	if l.Modules == nil {
		return nil, false
	}

	module, found := l.Modules[name]

	return module, found
}
