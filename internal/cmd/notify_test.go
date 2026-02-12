package cmd

import (
	"reflect"
	"testing"

	"github.com/dedene/lametric-cli/internal/api"
)

func TestBuildFrames_TextOnly(t *testing.T) {
	c := &NotifyCmd{Icon: ""}
	frames, err := c.buildFrames("Hello world")
	if err != nil {
		t.Fatalf("buildFrames: %v", err)
	}
	if len(frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(frames))
	}
	if frames[0].Text != "Hello world" {
		t.Errorf("text = %q, want %q", frames[0].Text, "Hello world")
	}
	if frames[0].Icon != "" {
		t.Errorf("icon = %q, want empty", frames[0].Icon)
	}
}

func TestBuildFrames_WithIcon(t *testing.T) {
	c := &NotifyCmd{Icon: "rocket"}
	frames, err := c.buildFrames("Test")
	if err != nil {
		t.Fatalf("buildFrames: %v", err)
	}
	if len(frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(frames))
	}
	// "rocket" alias resolves to animated rocket icon
	if frames[0].Icon != "a26304" {
		t.Errorf("icon = %q, want %q", frames[0].Icon, "a26304")
	}
}

func TestBuildFrames_Empty(t *testing.T) {
	c := &NotifyCmd{}
	_, err := c.buildFrames("")
	if err == nil {
		t.Fatal("expected error for empty frames")
	}
}

func TestBuildFrames_GoalFrame(t *testing.T) {
	c := &NotifyCmd{Goal: "50/100", Icon: "star"}
	frames, err := c.buildFrames("")
	if err != nil {
		t.Fatalf("buildFrames: %v", err)
	}
	if len(frames) != 1 {
		t.Fatalf("expected 1 frame (goal only), got %d", len(frames))
	}
	if frames[0].GoalData == nil {
		t.Fatal("expected GoalData")
	}
	if frames[0].GoalData.Current != 50 || frames[0].GoalData.End != 100 {
		t.Errorf("GoalData = %+v", frames[0].GoalData)
	}
}

func TestBuildFrames_ChartFrame(t *testing.T) {
	c := &NotifyCmd{Chart: "1,2,3"}
	frames, err := c.buildFrames("")
	if err != nil {
		t.Fatalf("buildFrames: %v", err)
	}
	if len(frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(frames))
	}
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(frames[0].ChartData, want) {
		t.Errorf("ChartData = %v, want %v", frames[0].ChartData, want)
	}
}

func TestParseGoal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *api.GoalData
		wantErr bool
	}{
		{"valid", "50/100", &api.GoalData{Start: 0, Current: 50, End: 100}, false},
		{"zero", "0/10", &api.GoalData{Start: 0, Current: 0, End: 10}, false},
		{"spaces", " 25 / 75 ", &api.GoalData{Start: 0, Current: 25, End: 75}, false},
		{"no slash", "50", nil, true},
		{"bad current", "abc/100", nil, true},
		{"bad max", "50/xyz", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGoal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseChart(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{"simple", "1,2,3,4,5", []int{1, 2, 3, 4, 5}, false},
		{"spaces", " 10 , 20 , 30 ", []int{10, 20, 30}, false},
		{"single", "42", []int{42}, false},
		{"invalid", "1,abc,3", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseChart(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
