package parseuml

import (
	"fmt"
	"strings"
)

func GenerateFromUml(umlStr string) (uml ParseUmlStage1, err error) {

	uml.Init()

	lines := strings.Split(umlStr, "\n")
	for lineIdx, line := range lines {
		area, parseErr := ParseLine(line, &uml)
		if parseErr != nil {
			err = fmtCompileErr(area, parseErr.Error(), line, lineIdx)
			return
		}
	}

	return
}

func fmtCompileErr(area, msg, line string, lineIdx int) error {
	return fmt.Errorf("fail(%s): compile line %d reason: %s line: %s", area, lineIdx+1, msg, line)
}
