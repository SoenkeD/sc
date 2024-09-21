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
	RepoOwner string `yaml:"repoOwner" validate:"required"`
	RepoName  string `yaml:"repoName" validate:"required"`
	RepoPath  string `yaml:"repoPath" validate:"required"`
	LocalPath string `yaml:"localPath" validate:"required"`
	Token     string `yaml:"token"`
}

type Config struct {
	CtlDir                     string            `yaml:"ctlDir" validate:"required"`
	EnableFileCapitalization   bool              `yaml:"enableFileCapitalization"`
	EnableGeneratedFilePrefix  bool              `yaml:"EnableGeneratedFilePrefix"`
	Exports                    []Export          `yaml:"exports"`
	ForceUnitSetupRegeneration bool              `yaml:"forceUnitSetupRegeneration"`
	ImportPathSeparator        string            `yaml:"importPathSeparator" validate:"required"`
	Imports                    []Import          `yaml:"imports"`
	Language                   string            `yaml:"language" validate:"required"`
	Module                     string            `yaml:"module"`
	RepoRoot                   string            `yaml:"repoRoot"`
	Templates                  []TemplatePackage `yaml:"templates"`
}
