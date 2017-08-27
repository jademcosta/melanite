package imageaction_test

import (
	"testing"

	"github.com/jademcosta/melanite/imageaction"
	"github.com/stretchr/testify/assert"
)

func TestItReportsTheCorrectMessage(t *testing.T) {
	err := imageaction.Error{Message: "Testing!"}

	assert.Equal(t, "Testing!", err.Error(), "Should be equal")
}
