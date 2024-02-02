package formatter

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

func removeTrailingNewLine(input string) string {
	return strings.TrimSuffix(input, "\n")
}

// FormatInput takes a string and returns a string with the input string formatted as a JSON string to be used as input to a model.
func FormatInput(input string) (string, error) {
	buf := new(bytes.Buffer)
	e := json.NewEncoder(buf)
	err := e.Encode(input)
	if err != nil {
		return "", err
	}
	// We do not want to include the trailing newline in the output
	result := removeTrailingNewLine(buf.String())
	return result, nil
}

type DistillFormat struct {
	InputText  string
	OutputText string
}

// FormatJsonLine takes two strings and returns a string with the input and output strings formatted as a JSON object.
// Input object should already be formatted with the remove trailiing new lines formatter
func FormatJsonLine(input string, output string) (string, error) {
	buf := new(bytes.Buffer)
	df := DistillFormat{InputText: input, OutputText: output}
	tmpl, err := template.New("jsonl").Parse(`{"input_text": "{{.InputText}}", "output_text": "{{.OutputText}}"}`)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(buf, df)
	if err != nil {
		return "", err
	}
	result := removeTrailingNewLine(buf.String())
	return result, nil
}
