package formatter_test

import (
	"testing"

	"github.com/andreyvit/diff"
	"github.com/nohe427/token-pipeline/formatter"
)

func TestFormatInput(t *testing.T) {

	var cases = []struct {
		input  string
		expect string
	}{
		{"hello world", "\"hello world\""},
		{`summarize the following github issue comments:
["Well @nohe427, I did think that was a good idea", "comment 2 [https://nohe.dev](https://nohe.dev)"]
`, "\"summarize the following github issue comments:\\n[\\\"Well @nohe427, I did think that was a good idea\\\", \\\"comment 2 [https://nohe.dev](https://nohe.dev)\\\"]\\n\""},
	}

	for _, tt := range cases {
		result, _ := formatter.FormatInput(tt.input)
		if result != tt.expect {
			t.Errorf("expected %s, got %s\ndiff %v", tt.expect, result, diff.CharacterDiff(tt.expect, result))
		}
	}
}

func TestFormatJsonLine(t *testing.T) {
	var cases = []struct {
		input  string
		output string
		expect string
	}{
		{"hello world", "goodbye world", "{\"input_text\": \"hello world\", \"output_text\": \"goodbye world\"}"},
	}
	for _, tt := range cases {
		result, _ := formatter.FormatJsonLine(tt.input, tt.output)
		if result != tt.expect {
			t.Errorf("expected %s, got %s\ndiff %v", tt.expect, result, diff.CharacterDiff(tt.expect, result))
		}
	}
}
