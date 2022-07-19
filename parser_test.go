package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLoadFromString(t *testing.T) {
	str := fmt.Sprint("[owner]\nname = John Doe\n[database]\nserver = 192.0.2.62\nonly key = \n")
	var pr Parser
	pr.LoadFromString(str)
	got := pr.mp
	var want = make(map[string]Section)
	want["owner"] = Section{"name": "John Doe"}
	want["database"] = Section{"server": "192.0.2.62", "only key": ""}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}

}

func TestLoadFromFile(t *testing.T) {
	t.Run("File doesn't exist", func(t *testing.T) {
		filePath := "somefile.ini"
		var pr Parser
		got := pr.LoadFromFile(filePath).Error()
		want := "The file \"" + filePath + "\" is not found!"
		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}

	})
	t.Run("File exists", func(t *testing.T) {
		filePath := "testLoad.ini"
		var pr Parser
		got := pr.LoadFromFile(filePath)
		if got != nil {
			t.Errorf("expected %v but got %v", nil, got)
		}

	})

}
func TestGetSection(t *testing.T) {
	var pr Parser
	pr.LoadFromFile("testLoad.ini")
	got := pr.GetSections()
	var want = make(map[string]Section)
	want["owner"] = Section{"name": "John Doe"}
	want["database"] = Section{"server": "192.0.2.62", "only key": ""}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestGetSectionsNames(t *testing.T) {
	var pr Parser
	pr.LoadFromFile("testLoad.ini")
	got := pr.GetSectionNames()
	want := []string{"owner", "database"}
	//helper that compare two slices ignoring the order of elements.
	less := func(a, b string) bool { return a < b }
	if cmp.Diff(want, got, cmpopts.SortSlices(less)) != "" {
		t.Errorf("expected %q but got %q", want, got)
	}
}

func TestGet(t *testing.T) {
	t.Run("Section doesn't exist", func(t *testing.T) {
		var pr Parser
		sectionName := "not a section"
		pr.LoadFromFile("testLoad.ini")
		_, err := pr.Get(sectionName, "server")
		got := err.Error()
		want := "Section " + sectionName + " not found"
		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}

	})
	t.Run("Key doesn't exist", func(t *testing.T) {
		var pr Parser
		key := "not a key"
		pr.LoadFromFile("testLoad.ini")
		_, err := pr.Get("database", key)
		got := err.Error()
		want := "No value found for the key " + key
		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}

	})
	t.Run("valid section and key", func(t *testing.T) {
		var pr Parser
		pr.LoadFromFile("testLoad.ini")
		got, _ := pr.Get("database", "server")
		want := "192.0.2.62"
		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}
	})
}

func TestSet(t *testing.T) {
	t.Run("New section", func(t *testing.T) {
		var pr Parser
		pr.LoadFromFile("testLoad.ini")
		pr.Set("newSection", "newKey", "newVal")
		got, _ := pr.Get("newSection", "newKey")
		want := "newVal"
		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}
	})

	t.Run("New key", func(t *testing.T) {
		var pr Parser
		pr.LoadFromFile("testLoad.ini")
		pr.Set("database", "newKey", "newVal")
		got, _ := pr.Get("database", "newKey")
		want := "newVal"
		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}
	})
	t.Run("Defined section and key", func(t *testing.T) {
		var pr Parser
		pr.LoadFromFile("testLoad.ini")
		pr.Set("database", "server", "newVal")
		got, _ := pr.Get("database", "server")
		want := "newVal"
		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}
	})
}
func TestString(t *testing.T) {
	str := fmt.Sprint("[owner]\nname = John Doe\n[database]\nserver = 192.0.2.62\nonly key = \n")
	var pr Parser
	pr.LoadFromString(str)
	want := pr.mp
	gotStr := pr.String()
	var pr2 Parser
	pr2.LoadFromString(gotStr)
	got := pr2.GetSections()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected\n%v but got\n%v", want, got)
	}

}
func TestSaveToFile(t *testing.T) {
	var pr Parser
	pr.LoadFromFile("testLoad.ini")
	pr.SaveToFile("testSave.ini")
}
