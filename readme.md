# SC - State Chart
State Chart is a tool that generates application code based on state charts
to help solving complex tasks. 
Using state charts in programming has many advantages. 
Some of the most important are:
- Easier to understand & serves as a communication basis
- Decouple overall behavior from components
- Scales good with increasing complexity 
- Supports documenting without additional effort

To get more information about state charts (in programming) 
or read about [the concept of SC](docs/concept.md),
consider visiting [statechart.dev](https://statecharts.dev),
which contains a nice collection of information and explanations 
regarding this topic. 

## Install sc
The build-in version check and update capabilities require a build flag to be set.
To do this there are currently two options:

1. Automatic
- Install from Github `go install github.com/SoenkeD/sc@main`
- Update it to set the build flag `sc version --update`

2. Manual
- Find [the most recent (long) commit hash of the main branch](https://github.com/SoenkeD/sc-go-templates/commits/main/) 
(e.g. 78bd6bc89736f4ac35ca7aaf1fe0f9cba4e31159) 
- Replace `COMMIT_HASH` with the recent has and execute `GOFLAGS=-ldflags=-X=main.commitHash=COMMIT_HASH go install github.com/SoenkeD/sc@main`

## Get Started
1. Install the `sc` tool from GitHub: \
`go install github.com/SoenkeD/sc@main`

2. (optional) [Read the documentation](docs/readme.md)

3. Decide on a language on following the the available guides and examples
- [Golang](https://github.com/SoenkeD/sc-go-templates)
- [Java](https://github.com/SoenkeD/sc-java-templates)