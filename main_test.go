package main

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_filterTestCase(t *testing.T) {
	as := assert.New(t)

	input := `?   	github.com/chyroc/debug-test-usage/a	[no test files]
--- FAIL: Test_RandFail (0.00s)
    --- FAIL: Test_RandFail/child-test (0.00s)
        --- FAIL: Test_RandFail/child-test/child-test-child (0.00s)
FAIL
FAIL	github.com/chyroc/debug-test-usage/b	0.399s
FAIL`

	res, err := filterFailTestCase(strings.NewReader(input))
	as.Nil(err)
	as.Equal([]string{"Test_RandFail", "Test_RandFail/child-test", "Test_RandFail/child-test/child-test-child"}, res)
}

func Test_MayFail(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	// 0 1 2 3
	if rand.Int63()%4 != 0 {
		// 0.75
		t.Fail()
	}
}
