package types

type StaticTemplates struct {
	CtlFile        string
	ReconcilerFile string
	TypesFile      string
	UtilsFile      string
}

type TransitionType string

const (
	TransitionTypeNormal TransitionType = "normal"
	TransitionTypeError  TransitionType = "error"
	TransitionTypeHappy  TransitionType = "happy"
)
