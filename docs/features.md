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

## export
tbd.