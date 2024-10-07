package args

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

const (
	defaultConfigParam = "defaultConfigParam"
)

func TestNewClientFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		args           []string
		expectedResult clientArgs
		expectedError  bool
		errorStr       string
	}{
		{
			name:           "empty falgs",
			args:           []string{""},
			expectedResult: clientArgs{},
			expectedError:  true,
			errorStr:       fmt.Errorf("%q is required", configPathParam).Error(),
		},
		{
			name:           "empty config flag",
			args:           []string{fmt.Sprintf("-%s", configPathParam), ""},
			expectedResult: clientArgs{},
			expectedError:  true,
			errorStr:       fmt.Errorf("%q is required", configPathParam).Error(),
		},
		{
			name:           "config flag name typo",
			args:           []string{fmt.Sprintf("-%serr", configPathParam), ""},
			expectedResult: clientArgs{},
			expectedError:  true,
			errorStr:       fmt.Errorf("%q is required", configPathParam).Error(),
		},
		{
			name: "valid flags",
			args: []string{
				fmt.Sprintf("-%s", configPathParam), defaultConfigParam,
			},
			expectedResult: clientArgs{configPath: defaultConfigParam},
			expectedError:  false,
			errorStr:       "",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			/* ARRANGE */

			// Update os.Args to simulate different command line arguments
			os.Args = append([]string{"test_application"}, testCase.args...)

			// Reset the flag values for each test case
			flag.CommandLine = flag.NewFlagSet(testCase.name, flag.ContinueOnError)

			/* ACT */
			flags, err := NewClientFlags()

			/* ASSERT */
			// Assert expected error string
			if (err != nil) && (err.Error() != testCase.errorStr) {
				t.Fatalf("NewClientFlags() with args %v: expected error string [%s], got [%s]", testCase.args, testCase.errorStr, err.Error())
			}

			// Assert expected error
			if (err != nil) != testCase.expectedError {
				t.Fatalf("NewClientFlags() with args %v: expected error %v, got %v", testCase.args, testCase.expectedError, err != nil)
			}

			// Assert result
			if flags != testCase.expectedResult {
				t.Fatalf("NewClientFlags() with args %v: expected %v, got %v", testCase.args, testCase.expectedResult, flags)
			}
		})
	}
}
