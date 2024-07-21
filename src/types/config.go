package types

type TemplatePackage struct {
	// Relative path of the template directory to the root dir
	Dir string `validate:"required"`

	// The name of the controllers which should use this
	// template package.
	// Controllers not in this list will ignore this package.
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
	RepoRoot                   string `validate:"required"`
	Language                   string `validate:"required"`
	ImportPathSeparator        string `validate:"required"`
	CtlDir                     string `validate:"required"`
	Module                     string
	EnableGeneratedFilePrefix  bool
	EnableFileCapitalization   bool
	ForceUnitSetupRegeneration bool
	Templates                  []TemplatePackage
	Exports                    []Export
	Imports                    []Import
}
