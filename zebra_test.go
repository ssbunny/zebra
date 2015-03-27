package zebra

import (
	"testing"
)

func TestServeHTTP(t *testing.T) {
	tester := NewTestHelper(t)

	z := New()
	m1 := &demoMiddleware{}
	m2 := &demoMiddleware2{}
	m3 := &demoMiddleware3{}
	z.Use(m1)
	z.Use(m2)
	z.Use(m3)

	z.ServeHTTP(nil, tester.NewRequest("GET", "/foo"))

	tester.AssertTrue(m1.excuted)
	tester.AssertTrue(m2.excuted)
	tester.AssertFalse(m3.excuted)
	tester.AssertTrue(m1.callbacked)
}

func TestRun(t *testing.T) {
	go New().Run()
	go New().RunOnAddr(":8888")
}

///-----------------------------------------------

type demoMiddleware struct {
	excuted    bool
	callbacked bool
}

func (m *demoMiddleware) Excute(c *Context) bool {
	m.excuted = true
	return true
}

func (m *demoMiddleware) Callback(c *Context) {
	m.callbacked = true
}

type demoMiddleware2 struct {
	excuted bool
}

func (m *demoMiddleware2) Excute(c *Context) bool {
	m.excuted = true
	return false
}

type demoMiddleware3 struct {
	excuted bool
}

func (m *demoMiddleware3) Excute(c *Context) bool {
	m.excuted = true
	return true
}
