package converter_test

import (
	"fmt"
	"testing"

	"github.com/jademcosta/melanite/converter"
	"github.com/stretchr/testify/assert"
)

func TestValidationWorksForImageEncodingsThatAreSupported(t *testing.T) {

	var imageEncodingValidationTests = []struct {
		format         string
		expectedResult bool
	}{
		{"png", true},
		{"jpg", true},
		{"", false},
		{"jpeg", false},
		{"gif", false},
		{"bmp", false},
		{"svg", false},
	}

	for _, sample := range imageEncodingValidationTests {
		assert.Equal(t,
			sample.expectedResult, converter.IsValidImageEncoding(sample.format),
			fmt.Sprintf("should be %v", sample.expectedResult))
	}

}
