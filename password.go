package survey

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
)

/*
Password is like a normal Input but the text shows up as *'s and there is no default. Response
type is a string.

	password := ""
	prompt := &survey.Password{ Message: "Please type your password" }
	survey.AskOne(prompt, &password)
*/
type Password struct {
	Renderer
	Message string
	Help    string
}

type PasswordTemplateData struct {
	Password
	ShowHelp bool
	Icons    *IconSet
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var PasswordQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ .Icons.Help }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ .Icons.Question }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ .Icons.HelpInput }} for help]{{color "reset"}} {{end}}`

func (p *Password) Prompt(config *PromptConfig) (line interface{}, err error) {
	// render the question template
	out, err := core.RunTemplate(
		PasswordQuestionTemplate,
		PasswordTemplateData{Password: *p, Icons: &config.IconSet},
	)
	fmt.Fprint(terminal.NewAnsiStdout(p.Stdio().Out), out)
	if err != nil {
		return "", err
	}

	rr := p.NewRuneReader()
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	// no help msg?  Just return any response
	if p.Help == "" {
		line, err := rr.ReadLine('*')
		return string(line), err
	}

	cursor := p.NewCursor()

	// process answers looking for help prompt answer
	for {
		line, err := rr.ReadLine('*')
		if err != nil {
			return string(line), err
		}

		if string(line) == string(config.IconSet.HelpInput) {
			// terminal will echo the \n so we need to jump back up one row
			cursor.PreviousLine(1)

			err = p.Render(
				PasswordQuestionTemplate,
				PasswordTemplateData{Password: *p, ShowHelp: true, Icons: &config.IconSet},
			)
			if err != nil {
				return "", err
			}
			continue
		}
		return string(line), err
	}
}

// Cleanup hides the string with a fixed number of characters.
func (prompt *Password) Cleanup(val interface{}, config *PromptConfig) error {
	return nil
}
