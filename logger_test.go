package zebra

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	tester := NewTestHelper(t)

	tester.AssertTrue(logger != nil)
}

func TestParse(t *testing.T) {
	tester := NewTestHelper(t)
	req := tester.NewRequest("GET", "/foo")
	cxt := newContext(nil, req)

	logger := NewLogger()
	logger.SetFormat(":url")
	r := logger.parse(cxt)
	tester.AssertEqual(r, "/foo")
}
