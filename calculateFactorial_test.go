package main

import (
	"errors"
	"testing"
)

func Test_ClaculateFactorial(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want int
		err  error
	}{
		{
			name: "factorial of -5",
			arg:  -5,
			want: 0,
			err:  nil,
		},
		{
			name: "factorial of 0",
			arg:  0,
			want: 0,
			err:  nil,
		},
		{
			name: "factorial of 3",
			arg:  3,
			want: 6,
			err:  ErrIncorrectMessage,
		},
		{
			name: "factorial of 20",
			arg:  20,
			want: 2432902008176640000,
			err:  ErrIncorrectMessage,
		},
		{
			name: "factorial of 21",
			arg:  21,
			want: 0,
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := calculateFactorial(tt.arg); got != tt.want || errors.Is(err, tt.err) {
				t.Errorf("calculateFactorial() = %v with error %v, want %v with error %v", got, err, tt.want, tt.err)
			}
		})
	}
}
