# Code Structure
In your state chart folder (location of the `.plantuml` file)
the `sc` tool will generate some directories & files. 
Some of the files have the prefix `zz_gen_`.
This prefix marks the files as being fully generated and 
should not be edited at all. 
These files can safely be removed as they will be generated if missing. 

The structure of your state chart folder will look the following:

- `actions/` \
This directory contains the code & tests for your actions. 
You should write your action code in the files without the generation prefix.

- `guards/` \
This directory contains the code & tests for your guards. action
You should write your guard code in the files without the generation prefix.

- `state/` \
This directory contains your controllers context and state.
Both files must exist.
You can add files in here as you wish.
It must not include any imports of the actions & guards module 
in order to prevent import cycles. 

- e.g. `demo.plantuml` \
This file contains the state chart as code.
It is the main artifact of each controller.

- `zz_gen_ctl.go` \
This generated file contains the entry point to this controller.
You call its `InitCtl()` function to obtain an initialized controller.
This file is generated and should not be edited. 

- `zz_gen_sm.go` \
This file contains the generated state machine,
which reflects your state chart from e.g `demo.plantuml`.
This file is generated every time you run the `sc gen` command.
Any changes will be removed then. Do NOT edit it. 