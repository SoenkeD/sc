# Cli Features


## gen
Executing `sc gen --name myctl` will generate 
all missing or regenerated code for the controller
based on the file at `CTL_ROOT/CTL_NAME/CTL_NAME.plantuml`.
It will always regenerate the `zz_gen_sm.go` file, 
which contains the paths and states of the state chart. 
It will not overwrite existing customizable files. 
When writing something it will use templates over the default generation.
So when generating an action it will use the template if available,
but it would not overwrite an existing action with that template.  

### (TODO) --clear
tbd.

### (TODO) --force
tbd.

## path
Executing `sc path --name myctl --route route.json` will generate
a visualization of the path taken through the state machine. 
It create / overwrites a file named `myctl.route.plantuml` 
in your current directory. 

![Example output](imgs/feature_log_example.png)


## import
Executing `sc import` will download the templates,
which are configured in `sc.yaml`,
to the target directories.
Afterwards the templates,
if also set as templates in `sc.yaml`,
are available to the `sc gen` command.
Currently only Github repositories are supported. 

### Add an import
1. In `sc.yaml` add the following with your own configuration.
```yaml
imports:
- repoOwner: "SoenkeD"
  repoName: "sc-go-templates"
  repoPath: "sc/templates/"
  localPath: "sc/templates/"
  # token: "abcdefg1245"
```

2. Add the template to `sc.yaml`
```yaml
templates:
  - dir: "sc/templates"
```
## init
Executing `sc init` initializes a project
based on a Github repository.
Find a sc template repository and
follow the instructions there. 

Golang example:
```bash
sc init --setup https://github.com/SoenkeD/sc-go-templates/main/sc/setup \
	--name myctl \
	--root $PWD/demo  \
	--module demo
```
`--name` is the name of the first controller to create \
`--root` is the desired root of the project (the directory should not exist) \
`--module` is the name of the desired Golang module e.g. `github.com/SoenkeD/sc`
`--container` can be used to enable podman (defaults to docker) 

## export
tbd.