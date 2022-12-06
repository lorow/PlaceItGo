package utils

import (
	"testing"
)

func TestConvertResolutionFormString(t *testing.T) {
	tests := []struct {
		name           string
		resolution     string
		delimiter      string
		expectedWidth  int
		expectedHeight int
		expectedError  bool
	}{
		{
			name:           "Valid resolution",
			resolution:     "1024x768",
			delimiter:      "x",
			expectedWidth:  1024,
			expectedHeight: 768,
			expectedError:  false,
		},
		{
			name:           "Invalid resolution format",
			resolution:     "1024/768",
			delimiter:      "x",
			expectedWidth:  -1,
			expectedHeight: -1,
			expectedError:  true,
		},
		{
			name:           "Invalid resolution value",
			resolution:     "1024xabc",
			delimiter:      "x",
			expectedWidth:  -1,
			expectedHeight: -1,
			expectedError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			width, height, err := ConvertResolutionFormString(test.resolution, test.delimiter)

			if width != test.expectedWidth {
				t.Errorf("unexpected width value. expected: %d, got: %d", test.expectedWidth, width)
			}

			if height != test.expectedHeight {
				t.Errorf("unexpected height value. expected: %d, got: %d", test.expectedHeight, height)
			}

			if err == nil && test.expectedError {
				t.Errorf("expected: %v", test.expectedError)
			}
		})
	}
}
