# SC Project Configuration
The `sc.yaml` file is used to configure a sc project. 
Most often a code language specific configuration template
should be used. 

## Configuration Documentation
```yaml
# directory where the different controllers are located
ctlDir: "src/controller"

# if true capitalizes the first character of file names
enableFileCapitalization: false

# if true marks auto generated files with the prefix `zz_gen_`
enableGeneratedFilePrefix: true

# list of exports configurations
# to export code files to templates when executing sc export
exports:
# file relative to controller root e.g. actions/MyAction.go
# you can also do actions/* to include all files which are not excluded
- src: "my-ctl/actions/*"
  to: "sc/templates/actions"
  # list of files to exclude from export
  # e.g MyExcludedAction.go
  excluded:
  - "MyExcludedAction.go"

# if true always regenerate actions and guards setup file
forceUnitSetupRegeneration: false

# the separator symbol for import paths e.g. / in example.com/example/example-sc
importPathSeparator: "/"

# list of template import configurations
imports:
# name of the Github repository owner
- repoOwner: "SoenkeD"
# name of the Github repository
  repoName: "sc-go-templates"
# directory inside source repository where the templates are located
  repoPath: "sc/templates/"
# local directory to save the templates to
  localPath: "sc/templates/"
# Github access token when using a private repository
  # token: "abcdefg12345"

# the application languages file suffix without the dot
language: "go"

# import path prefix - may be empty 
module: "example.com/example/example-sc"

# list of export configurations 
# root of the repository may be entered with a flag
# e.g. sc --root PWD
repoRoot: "/home/example/code/my-sc"

# list of templates package configurations 
# to import when executing sc import

templates:
#  relative path of the template directory to the root dir. 
# the template directory contains e.g. actions, guards, base, ...
  - dir: "sc/templates"
# the name of the controllers which should use this
# template package.
# controllers not in this list will ignore this package,
# if the list is empty.
    exclusive: 
    - "my-ctl"
```