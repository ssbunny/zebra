package zebra

import (
	"fmt"
	"net/http"
	"testing"
)

type TestHelper interface {
	AssertTrue(bool)
	AssertFalse(bool)
	AssertEqual(string, string)

	NewRequest(string, string) *http.Request
}

type testHelper struct {
	t *testing.T
}

func NewTestHelper(t *testing.T) TestHelper {
	return &testHelper{t}
}

func (helper *testHelper) AssertTrue(val bool) {
	if !val {
		helper.err("// TODO Error Message!")
	}
}

func (helper *testHelper) AssertFalse(val bool) {
	helper.AssertTrue(!val)
}

func (helper *testHelper) AssertEqual(des, obj string) {
	if des != obj {
		helper.err(fmt.Sprintf(` want "%s", but "%s" `, des, obj))
	}
}

func (helper *testHelper) NewRequest(method, path string) *http.Request {
	req, _ := http.NewRequest(method, "http://localhost:3000"+path, nil)
	return req
}

func (helper *testHelper) err(msg string) {
	helper.t.Error(msg)
}
