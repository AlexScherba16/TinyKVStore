package helpers_test

import (
	"TinyKVStore/internal/helpers"
	"testing"
)

func TestNewClientFlags(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult int
	}{
		{
			name:           "empty buffer size",
			input:          "",
			expectedResult: -1,
		},
		{
			name:           "invalid input, invalid data pattern",
			input:          "invalid input d@ta",
			expectedResult: -1,
		},
		{
			name:           "valid size data with lower multiplier postfix",
			input:          "1kb",
			expectedResult: 1 * 1024,
		},
		{
			name:           "valid size data with upper multiplier postfix",
			input:          "123KB",
			expectedResult: 123 * 1024,
		},
		{
			name:           "valid size data with combined multiplier postfix",
			input:          "123mB",
			expectedResult: 123 * 1024 * 1024,
		},
		{
			name:           "invalid multiplier postfix",
			input:          "123mBb",
			expectedResult: -1,
		},
		{
			name:           "invalid data, several multipliers",
			input:          "Kb 123 Mb",
			expectedResult: -1,
		},
		{
			name:           "invalid data, several sizes",
			input:          "456Kb 123 Mb",
			expectedResult: -1,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */

			/* ACT */
			result := helpers.ParseBufferSize(testCase.input)

			/* ASSERT */
			// Assert result
			if result != testCase.expectedResult {
				t.Fatalf("ParseBufferSize() with input %v: expected %v, got %v", testCase.input, testCase.expectedResult, result)
			}
		})
	}
}
