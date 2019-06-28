package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateSha512(t *testing.T) {
	var test = []struct {
		Text         string
		ExpectedHash string
	}{
		{"codersrank", "9c2d1e94334bb1efbf38fcc68a0cae807ec67c38aa82c6a3df61b84bb8febeefbb0555fd09f0933bed849f7ee2f447e6f8a1b95344a10ca55e9991b55d2b5397"},
		{"golang", "df84c5d44709cfeb8a22c8cf006ac926c92c6823d37e112f2c68a22890e61615f97ad1d4eb1d3e043442063886b4ce2f15eaa73ea8ff769808fc76d47f607ec5"},
		{"$-!@%-&", "a371f26f8a8f460ecf4b2b789aca1780a49b5aece9a29d2d0c71fbd8b375d8348fd7dbfb3f25676bede77b39fa45ba8e9ea24825313edcd3fcaa4d60f429596d"},
	}

	for _, tc := range test {
		hash := CalculateSha512(tc.Text)
		assert.Equal(t, tc.ExpectedHash, hash, "The hash's should be equals")
	}
}

func TestCalculateSha512NOK(t *testing.T) {
	var test = []struct {
		Text         string
		ExpectedHash string
	}{
		{"codersrank", "1"},
	}

	for _, tc := range test {
		hash := CalculateSha512(tc.Text)
		assert.NotEqual(t, tc.ExpectedHash, hash, "The hash's should NOT be equals")
	}
}
