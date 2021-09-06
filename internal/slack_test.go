package internal

import (
	"testing"
	"time"
)

func TestGetTimeFromSlackTimeStamp(t *testing.T) {
	type args struct {
		ts string
	}
	tests := []struct {
		name     string
		args     args
		expected time.Time
	}{
		{
			name: "normal",
			args: args{
				ts: "1630923736.000200",
			},
			expected: time.Unix(1630923736, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := GetTimeFromSlackTimeStamp(tt.args.ts)
			if err != nil {
				t.Errorf("failed to call GetTimeFromSlackTimeStamp()")
			}
			if tt.expected != actual {
				t.Errorf("expected = %v, actual = %v", tt.expected, actual)
			}
		})
	}
}
