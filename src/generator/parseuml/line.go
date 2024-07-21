package parseuml

import (
	"fmt"
	"strings"
)

const UML_PREFIX = "@startuml"

func ParseLine(line string, uml *ParseUmlStage1) (string, error) {
	// trim spaces
	line = strings.TrimSpace(line)

	// filter empty line
	if len(line) == 0 {
		return "", nil
	}

	// strip comment suffix
	lineParts := strings.SplitN(line, "'", 2)
	if len(lineParts) == 2 {
		line = strings.TrimSpace(lineParts[0])

		if len(line) == 0 {
			return "", nil
		}
	}

	// parse headline
	if strings.HasPrefix(line, UML_PREFIX) {
		words := strings.SplitN(line, " ", 2)
		if len(words) != 2 {
			return "head", fmt.Errorf("missing head")
		}
		uml.Name = words[1]
		return "", nil
	}

	// parse state group closing
	if strings.HasPrefix(line, "}") {
		uml.AddStateGroupClosing()
		return "", nil
	}

	// detect state group opening
	if StateRegex.MatchString(line) {
		prefix := "state "
		postfix := " {" // TODO: allow missing space idea trim
		stateGroup := line[len(prefix) : len(line)-len(postfix)]
		uml.AddStateGroup(stateGroup)
		return "", nil
	}

	// try parsing content line
	contentParts := strings.Split(line, ":")
	if len(contentParts) > 2 {
		return "content_line", fmt.Errorf("invalid number of content line parts: found %d expected %d", len(contentParts), 2)
	}

	// parse first part
	tokens := strings.Split(strings.TrimSpace(contentParts[0]), " ")
	for tokenIdx, token := range tokens {
		tokens[tokenIdx] = strings.TrimSpace(token)
	}

	// assume is an action and parse 2nd part
	if len(tokens) == 1 {
		if len(contentParts) != 2 {
			return "state_action", fmt.Errorf("expected ':' %d time got %d", 1, len(contentParts)-1)
		}

		action, actionParams, parseErr := ParseStateActionFromCode(contentParts[1])
		if parseErr != nil {
			return "state_action", fmt.Errorf(parseErr.Error())
		}
		uml.AddStateAction(contentParts[0], action, actionParams)
		return "", nil
	}

	// detect transition
	var contentPart2 string
	// if content part exists use it
	// else use empty string
	if len(contentParts) == 2 {
		contentPart2 = contentParts[1]
	}
	var parseErr error
	parsedTa, parseErr := ParseTransition(tokens, contentPart2)
	if parseErr != nil {
		return "transition", fmt.Errorf("%s", parseErr.Error())
	}
	uml.AddTransition(parsedTa)

	return "", nil
}
