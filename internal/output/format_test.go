package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestModeFromFlags_JSON(t *testing.T) {
	m := ModeFromFlags(true, false)
	if m != ModeJSON {
		t.Errorf("got %v, want ModeJSON", m)
	}
}

func TestModeFromFlags_Plain(t *testing.T) {
	m := ModeFromFlags(false, true)
	if m != ModePlain {
		t.Errorf("got %v, want ModePlain", m)
	}
}

func TestModeFromFlags_Default(t *testing.T) {
	m := ModeFromFlags(false, false)
	if m != ModeTable {
		t.Errorf("got %v, want ModeTable", m)
	}
}

func TestModeFromFlags_JSONPrecedence(t *testing.T) {
	m := ModeFromFlags(true, true)
	if m != ModeJSON {
		t.Errorf("got %v, want ModeJSON (JSON takes precedence)", m)
	}
}

func TestMode_String(t *testing.T) {
	tests := []struct {
		mode Mode
		want string
	}{
		{ModeJSON, "json"},
		{ModePlain, "plain"},
		{ModeTable, "table"},
	}
	for _, tt := range tests {
		got := tt.mode.String()
		if got != tt.want {
			t.Errorf("Mode(%d).String() = %q, want %q", tt.mode, got, tt.want)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"key": "value"}
	if err := WriteJSON(&buf, data); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}

	var got map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("got %v", got)
	}
}

func TestWriteTSV(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"NAME", "VALUE"}
	rows := [][]string{{"a", "1"}, {"b", "2"}}

	if err := WriteTSV(&buf, headers, rows); err != nil {
		t.Fatalf("WriteTSV: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(lines), lines)
	}
	if !strings.Contains(lines[0], "NAME") {
		t.Errorf("header missing: %q", lines[0])
	}
	if !strings.Contains(lines[1], "a\t1") {
		t.Errorf("row 1: %q", lines[1])
	}
}

func TestWriteTSV_NoHeaders(t *testing.T) {
	var buf bytes.Buffer
	rows := [][]string{{"x", "y"}}

	if err := WriteTSV(&buf, nil, rows); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "x\ty") {
		t.Errorf("got %q", buf.String())
	}
}

func TestFormatter_Output_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, true, false, true)

	data := []map[string]string{{"id": "1"}}
	if err := f.Output(data, []string{"ID"}, [][]string{{"1"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"id"`) {
		t.Errorf("expected JSON output, got %q", buf.String())
	}
}

func TestFormatter_Output_Plain(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, false, true, true)

	if err := f.Output(nil, []string{"A", "B"}, [][]string{{"1", "2"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "A\tB") {
		t.Errorf("expected TSV header, got %q", buf.String())
	}
}

func TestFormatter_Output_Table(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, false, false, true) // noColor=true => SimpleTable

	if err := f.Output(nil, []string{"X"}, [][]string{{"val"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "X") || !strings.Contains(buf.String(), "val") {
		t.Errorf("expected table output, got %q", buf.String())
	}
}

func TestFormatter_OutputSingle_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, true, false, true)

	data := map[string]string{"k": "v"}
	if err := f.OutputSingle(data, [][2]string{{"Key", "v"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"k"`) {
		t.Errorf("expected JSON, got %q", buf.String())
	}
}

func TestFormatter_OutputSingle_Plain(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, false, true, true)

	if err := f.OutputSingle(nil, [][2]string{{"Name", "dev1"}, {"IP", "10.0.0.1"}}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Name\tIP") || !strings.Contains(out, "dev1\t10.0.0.1") {
		t.Errorf("expected TSV, got %q", out)
	}
}

func TestFormatter_OutputSingle_Table(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, false, false, true)

	if err := f.OutputSingle(nil, [][2]string{{"Key", "Value"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Key") {
		t.Errorf("expected KV output, got %q", buf.String())
	}
}

func TestWriteKV(t *testing.T) {
	var buf bytes.Buffer
	pairs := [][2]string{
		{"Name", "test"},
		{"IP", "10.0.0.1"},
	}
	if err := WriteKV(&buf, pairs, nil); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Name") || !strings.Contains(out, "10.0.0.1") {
		t.Errorf("got %q", out)
	}
}

func TestSimpleTable(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"A", "B"}
	rows := [][]string{{"1", "2"}, {"3", "4"}}

	if err := SimpleTable(&buf, headers, rows); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "A") || !strings.Contains(out, "3") {
		t.Errorf("got %q", out)
	}
}

func TestSimpleTable_NoHeaders(t *testing.T) {
	var buf bytes.Buffer
	if err := SimpleTable(&buf, nil, [][]string{{"x"}}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "x") {
		t.Errorf("got %q", buf.String())
	}
}

func TestRenderTable_NilColors(t *testing.T) {
	var buf bytes.Buffer
	// nil colors => falls back to SimpleTable
	if err := RenderTable(&buf, []string{"H"}, [][]string{{"v"}}, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "H") {
		t.Errorf("got %q", buf.String())
	}
}

func TestRenderTable_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderTable(&buf, nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for empty data, got %q", buf.String())
	}
}

func TestColors_NoColor(t *testing.T) {
	c := NewColors(true)
	if c.Enabled() {
		t.Error("expected colors disabled")
	}
	// Methods should return input unchanged.
	if c.Bold("x") != "x" {
		t.Error("Bold should passthrough")
	}
	if c.Success("x") != "x" {
		t.Error("Success should passthrough")
	}
	if c.Error("x") != "x" {
		t.Error("Error should passthrough")
	}
	if c.Warning("x") != "x" {
		t.Error("Warning should passthrough")
	}
	if c.Dim("x") != "x" {
		t.Error("Dim should passthrough")
	}
}

func TestIsColorEnabled_NoColor(t *testing.T) {
	if IsColorEnabled(true) {
		t.Error("expected false with noColor=true")
	}
}
