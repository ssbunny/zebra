package zebra

import (
	"testing"
)

func TestGet(t *testing.T) {
	z := New()
	r := NewRouter()
	s := ""

	r.Get("/foo", func(c *Captain) {
		s += "foo1"
	})
	r.Get("/foo", func(c *Captain) {
		s += "foo2"
	})
	r.Get(
		"/bar",
		func(c *Captain) {
			s += "bar1"
		},
		func(c *Captain) {
			s += "bar2"
		},
	)
	z.Use(r)

	tester := NewTestHelper(t)
	z.ServeHTTP(nil, tester.NewRequest("GET", "/foo"))
	z.ServeHTTP(nil, tester.NewRequest("GET", "/bar"))

	tester.AssertEqual("foo1foo2bar1bar2", s)
}

func TestAny(t *testing.T) {
	z := New()
	r := NewRouter()
	s := ""
	r.Get("/foo", func(c *Captain) {
		s += "foo1"
	})
	r.Any("/foo", func(c *Captain) {
		s += "foo2"
	})
	r.Any("/bar", func(c *Captain) {
		s += "bar"
	})

	z.Use(r)

	tester := NewTestHelper(t)
	z.ServeHTTP(nil, tester.NewRequest("GET", "/foo"))
	z.ServeHTTP(nil, tester.NewRequest("GET", "/bar"))

	tester.AssertEqual("foo1foo2bar", s)
}

func TestOtherMethods(t *testing.T) {
	z := New()
	r := NewRouter()
	s := ""
	r.Post("/foo", func(c *Captain) {
		s += "post"
	})
	r.Put("/foo", func(c *Captain) {
		s += "put"
	})
	r.Delete("/foo", func(c *Captain) {
		s += "delete"
	})
	r.Options("/foo", func(c *Captain) {
		s += "options"
	})
	z.Use(r)

	tester := NewTestHelper(t)
	z.ServeHTTP(nil, tester.NewRequest("POST", "/foo"))
	z.ServeHTTP(nil, tester.NewRequest("PUT", "/foo"))
	z.ServeHTTP(nil, tester.NewRequest("DELETE", "/foo"))
	z.ServeHTTP(nil, tester.NewRequest("OPTIONS", "/foo"))

	tester.AssertEqual("postputdeleteoptions", s)
}

func TestPathRules(t *testing.T) {
	z := New()
	r := NewRouter()
	s := ""

	tester := NewTestHelper(t)

	r.Get("/foo/:id", func(c *Captain) {
		tester.AssertEqual(c.Path("id"), "111")
		s += "1"
	})
	r.Get("/foo/:id/bar/:name", func(c *Captain) {
		tester.AssertEqual(c.Path("id"), "222")
		tester.AssertEqual(c.Path("name"), "aaa")
		s += "2"
	})
	r.Get("/foo/:first/:second/:third/:forth", func(c *Captain) {
		tester.AssertEqual(c.Path("first"), "1")
		tester.AssertEqual(c.Path("second"), "2")
		tester.AssertEqual(c.Path("third"), "3")
		tester.AssertEqual(c.Path("forth"), "4")
		s += "3"
	})
	r.Get("/bar/:bar{name[\\d]+}", func(c *Captain) {
		tester.AssertEqual(c.Path("bar"), "name123")
		s += "4"
	})
	z.Use(r)

	z.ServeHTTP(nil, tester.NewRequest("GET", "/foo/111"))
	z.ServeHTTP(nil, tester.NewRequest("GET", "/foo/222/bar/aaa"))
	z.ServeHTTP(nil, tester.NewRequest("GET", "/foo/1/2/3/4"))
	z.ServeHTTP(nil, tester.NewRequest("GET", "/bar/name123"))

	tester.AssertEqual("1234", s)
}

func TestWriteJSON(t *testing.T) {
	// TODO 测试各种JSON转化是否正确
}
