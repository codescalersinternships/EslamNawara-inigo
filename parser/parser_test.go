package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetSection(t *testing.T) {
	var pr Parser
	pr.LoadFromFile("testLoad.ini")

	got := pr.GetSections()

	var want = make(map[string]Section)
	want["owner"] = Section{"name": "John Doe"}
	want["database"] = Section{"server": "192.0.2.62"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestGetSectionsNames(t *testing.T) {
	var pr Parser
	err := pr.LoadFromFile("testLoad.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	got := pr.GetSectionNames()

	want := []string{"owner", "database"}

	//helper that compare two slices ignoring the order of elements.
	less := func(a, b string) bool { return a < b }
	if cmp.Diff(want, got, cmpopts.SortSlices(less)) != "" {
		t.Errorf("expected %q but got %q", want, got)
	}
}

func TestGet(t *testing.T) {
	var pr Parser
	pr.LoadFromFile("testLoad.ini")

	got, _ := pr.Get("database", "server")

	want := "192.0.2.62"
	if got != want {
		t.Errorf("expected %q but got %q", want, got)
	}
}

func TestSet(t *testing.T) {
	var pr Parser
	pr.LoadFromFile("testLoad.ini")

	pr.Set("database", "Server", "changed value")
	got, _ := pr.Get("database", "Server")

	want := "changed value"
	if got != want {
		t.Errorf("expected %q but got %q", want, got)
	}
}

func TestSaveToFile(t *testing.T) {
	var pr Parser
	pr.LoadFromFile("testLoad.ini")
	pr.SaveToFile("testSave.ini")
}
