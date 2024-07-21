package types

type TemplatePackage struct {
	// Relative path of the template directory to the root dir
	Dir string `validate:"required"`

	// The name of the controllers which should use this
	// template package.
	// Controllers not in this list will ignore this package,
	// if the list is empty.
	Exclusive []string
}

type Export struct {
	// controller from where to export items from
	Controller string `validate:"required"`

	// items to export
	Items []ExportItem `validate:"required"`
}

type ExportItem struct {
	// file relative to controller root e.g. actions/MyAction.go
	// you can also do actions/* to include all files which are not excluded
	Src string `validate:"required"`

	// export dir root relative to execution path to export to e.g sc/templates/examples/actions
	To string `validate:"required"`

	// files to ignore when exporting
	// only usable when using * ending for the source
	Excluded []string
}

type Import struct {
	RepoOwner string `validate:"required"`
	RepoName  string `validate:"required"`
	RepoPath  string `validate:"required"`
	LocalPath string `validate:"required"`
	Token     string
}

type Config struct {
	CtlDir                     string `validate:"required"`
	EnableFileCapitalization   bool
	EnableGeneratedFilePrefix  bool
	Exports                    []Export
	ForceUnitSetupRegeneration bool
	ImportPathSeparator        string `validate:"required"`
	Imports                    []Import
	Language                   string `validate:"required"`
	Module                     string
	RepoRoot                   string `validate:"required"`
	Templates                  []TemplatePackage
}
