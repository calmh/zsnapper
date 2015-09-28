package main

import (
	"reflect"
	"testing"
)

func TestExpandLines(t *testing.T) {
	in := []string{
		"a-$(echo t1; echo t2)-b",
		"a-$(echo t1; echo t2)-b-$(ls _testdata)-c",
	}

	expected := []string{
		"a-t1-b",
		"a-t2-b",
		"a-t1-b-bar-c",
		"a-t1-b-foo-c",
		"a-t2-b-bar-c",
		"a-t2-b-foo-c",
	}

	out := expandLines(in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Lines differ:\n%#v !=\n%#v", out, expected)
	}
}

func TestSelectLines(t *testing.T) {
	in := []string{
		"data/test",
		"data/test/child",
		"data/foo",
		"data/foo1",
		"data/foo2",
		"data/foo2/child",
		"data/test21/foo",
		"data/test21/bar",
		"data/test31/foo",
		"data/test31/bar",
	}

	pats := []string{
		"data/test",
		"data/foo*",
		"data/test2?/*",
	}

	expected := []string{
		"data/test",
		"data/foo",
		"data/foo1",
		"data/foo2",
		"data/test21/foo",
		"data/test21/bar",
	}

	out := selectLines(pats, in)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("Selections differ:\n%#v !=\n%#v", out, expected)
	}
}
