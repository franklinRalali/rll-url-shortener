package rands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name       string
		args       args
		wantStrLen int
		wantErr    bool
	}{
		{
			name: "success generating random string",
			args: args{
				n: 8,
			},
			wantStrLen: 8,
			wantErr:    false,
		},
		{
			name: "error generating random string for n < 0",
			args: args{
				n: -1,
			},
			wantStrLen: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		randStr, err := RandString(tt.args.n)
		if tt.wantErr {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}

		assert.Equal(t, tt.wantStrLen, len(randStr))
	}
}
