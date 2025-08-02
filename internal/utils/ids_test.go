package utils_test

import (
	"testing"

	"github.com/open-ug/conveyor/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomID(t *testing.T) {
	// Generate first ID
	id1, err1 := utils.GenerateRandomID()
	assert.NoError(t, err1)
	assert.Len(t, id1, 24, "ID should be 24 hex characters long")

	// Generate second ID
	id2, err2 := utils.GenerateRandomID()
	assert.NoError(t, err2)
	assert.Len(t, id2, 24, "ID should be 24 hex characters long")

	// Ensure they are not equal (highly unlikely due to randomness)
	assert.NotEqual(t, id1, id2, "Generated IDs should be unique")
}
