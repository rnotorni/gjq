package gjq

import (
	"testing"
)

// TestNewGJQ .
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
