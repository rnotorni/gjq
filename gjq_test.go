package gjq

import (
	"testing"
)

func TestNewGJQ(t *testing.T) {
	tests := []struct {
		name      string
		script    string
		isSuccess bool
	}{
		{
			name:      "success",
			script:    `.hoge|.fuga`,
			isSuccess: true,
		},
		{
			name:      "success",
			script:    `hoge|`,
			isSuccess: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gjq, err := NewGJQ(tt.script)
			if gjq != nil {
				defer gjq.Close()
			}

			if tt.isSuccess != (err == nil) {
				t.Errorf("failed")
			}
		})
	}
}


func Test_GJQRun(t *testing.T) {
	tests := []struct {
		name string
		script    string
		input string
		expected string
	}{
		{
			name:      "success",
			script:    `.hoge`,
			input: `{"hoge":{"fuga":"piyo"}}`,
			expected: `{"fuga":"piyo"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gjq,_ := NewGJQ(tt.script)
			defer gjq.Close()
			if actual, err := gjq.Run(tt.input); tt.expected != actual {
				t.Errorf("%v", err)
				t.Errorf("%v != %v", tt.expected, actual)
			}
		})
	}
}
