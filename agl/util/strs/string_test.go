package strs_test

import (
	"testing"

	. "github.com/leclecr04/go-tool/agl/util/strs"
	"github.com/stretchr/testify/assert"
)

func TestDepunct(t *testing.T) {
	assert.Equal(t, "AaaBbb", Depunct("aaa.bbb", true))
	assert.Panics(t, func() { Depunct("aaa.bbb.", true) })

	assert.Equal(t, "ReviewerIDs", Depunct("reviewer_ids", true))
	assert.Equal(t, "reviewer_ids", CamelToSnake("ReviewerIDs"))

	assert.Equal(t, "C2CProfile", Depunct("c2c_profile", true))
	assert.Equal(t, "c2c_profile", CamelToSnake("C2CProfile"))

	assert.Equal(t, "BoardJSON", Depunct("board_json", true))
	assert.Equal(t, "board_json", CamelToSnake("BoardJSON"))
}

func TestStringSetToSlice(t *testing.T) {
	dit := map[string]bool{
		"bca": true,
		"bac": true,
		"abc": true,
		"cba": true,
		"acb": true,
		"cab": true,
	}
	assert.Equal(t, []string{"abc", "acb", "bac", "bca", "cab", "cba"}, SetToSlice(dit))
}
