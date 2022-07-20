package parser

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

var (
	testLoad = "testingFiles/testLoad.ini"
	testSave = "testingFiles/testSave.ini"
)
var (
	validMap    = map[string]Section {
		"owner":    {"name": "John Doe"},
		"database": {"server": "192.0.2.62", "only key": ""},
	}
	valNoKey        = "[section]\n=value without key"
	badSection      = "[Bad section"
	dataWithSection = "[More] data with section"
	multiVal        = "[Multiple]\nvalues=for=a key"

	validString     = 
    `[owner]
    name = John Doe
    [database]
    server = 192.0.2.62
    only key = `
)

func TestLoadFromString(t *testing.T) {
	t.Run("Good string format", func(t *testing.T) {
		var pr Parser
		want := validMap
		pr.LoadFromString(validString)
		got := pr.mp
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
	t.Run("Empty string", func(t *testing.T) {
		var pr Parser
		pr.LoadFromString("\n\n")
		got := pr.mp
		var want = make(map[string]Section)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})

	t.Run("value without key", func(t *testing.T) {
		var pr Parser
		got := pr.LoadFromString(valNoKey).Error()
		want := "Key not found in line 2 in the string\n" + valNoKey
		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
	t.Run("Bad section", func(t *testing.T) {
		var pr Parser
		got := pr.LoadFromString(badSection).Error()
		want := "Invalid section in line 1 in the string\n" + badSection
		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
	t.Run("More data with section", func(t *testing.T) {
		var pr Parser
		got := pr.LoadFromString(dataWithSection).Error()
		want := "Too much data for the section name in line 1 in the string\n" + dataWithSection
		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
	t.Run("Multiple values for a key", func(t *testing.T) {
		var pr Parser
		got := pr.LoadFromString(multiVal).Error()
		want := "Too much values for one key in line 2 in the string\n" + multiVal
		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
	t.Run("Value without section in a bad file", func(t *testing.T) {
		filePath := "testingFiles/badFile.ini"
		var pr Parser
		got := pr.LoadFromFile(filePath).Error()
		want := "File contains values doesn't belong to any section in line 1 in the file " + filePath
		if got != want {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
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
		filePath := testLoad
		var pr Parser
		got := pr.LoadFromFile(filePath)
		if got != nil {
			t.Errorf("expected %v but got %v", nil, got)
		}
	})
}
func TestGetSection(t *testing.T) {
	var pr Parser
	pr.LoadFromFile(testLoad)
	got := pr.GetSections()
	var want = validMap
    if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestGetSectionsNames(t *testing.T) {
	var pr Parser
	pr.LoadFromFile(testLoad)
	got := pr.GetSectionNames()
	want := []string{"database", "owner"}
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v but got %v", want, got)
	}
}

func TestGet(t *testing.T) {
	t.Run("Section doesn't exist", func(t *testing.T) {
		var pr Parser
		sectionName := "not a section"
		pr.LoadFromFile(testLoad)
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
		pr.LoadFromFile(testLoad)
		_, err := pr.Get("database", key)
		got := err.Error()
		want := "No value found for the key " + key
		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}

	})
	t.Run("valid section and key", func(t *testing.T) {
		var pr Parser
		pr.LoadFromFile(testLoad)
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
		pr.LoadFromFile(testLoad)
		pr.Set("newSection", "newKey", "newVal")
		got, _ := pr.Get("newSection", "newKey")
		want := "newVal"
		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}
	})

	t.Run("New key", func(t *testing.T) {
		var pr Parser
		pr.LoadFromFile(testLoad)
		pr.Set("database", "newKey", "newVal")
		got, _ := pr.Get("database", "newKey")
		want := "newVal"
		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}
	})
	t.Run("Defined section and key", func(t *testing.T) {
		var pr Parser
		pr.LoadFromFile(testLoad)
		pr.Set("database", "server", "newVal")
		got, _ := pr.Get("database", "server")
		want := "newVal"
		if got != want {
			t.Errorf("expected %q but got %q", want, got)
		}
	})
}
func TestString(t *testing.T) {
	var pr Parser
	pr.LoadFromString(validString)
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
	pr.LoadFromFile(testLoad)
	pr.SaveToFile(testSave)
	var got, want Parser
	got.LoadFromFile(testLoad)
	want.LoadFromFile(testSave)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected\n%v but got\n%v", want, got)
	}
}

// Example functions

func ExampleParser_LoadFromFile() {
	// Create a parser and fill it with data parsed from a file
	var pr Parser
	pr.LoadFromFile(testLoad)

	// Create another parser and fill it with data parsed from the same file
	var pr2 Parser
	pr2.LoadFromFile(testLoad)

	fmt.Println(reflect.DeepEqual(pr, pr2))
	//output: true
}

func ExampleParser_LoadFromString() {
	// Create a parser and fill it with data parsed from testLoad file
	var pr Parser
	pr.LoadFromFile(testLoad)

	// Create generated parser from the string resulted from the same file
	var pr2 Parser
	pr2.LoadFromString(pr.String())

	fmt.Println(reflect.DeepEqual(pr, pr2))
	//output: true
}

func ExampleParser_GetSectionNames() {
	// Create a parser and fill it with data parsed from testLoad file
	var pr Parser
	pr.LoadFromFile(testLoad)

	sectionNames := pr.GetSectionNames()

	// Sort the resulting slice to always match the output example
	sort.Strings(sectionNames)

	fmt.Println(sectionNames)
	// output: [database owner]
}

func ExampleParser_GetSections() {
	// Create a parser and fill it with data parsed from testLoad
	var pr Parser
	pr.LoadFromFile(testLoad)

	// Create another parser and fill it with data parsed from the same file
	var pr2 Parser
	pr2.LoadFromFile(testLoad)

	fmt.Println(reflect.DeepEqual(pr, pr2))
	// output: true
}

func ExampleParser_Get() {
	// Create a parser and fill it with data parsed from testLoad
	var pr Parser
	pr.LoadFromFile(testLoad)

	valueField, _ := pr.Get("owner", "name")
	fmt.Println(valueField)
	// output: John Doe
}

func ExampleParser_Set() {
	// Create a parser and fill it with data parsed from testLoad
	var pr Parser
	pr.LoadFromFile(testLoad)

	// Sets the entity with key "name" in section "owner" to value "person"
	pr.Set("owner", "name", "person")
	valueField, _ := pr.Get("owner", "name")
	fmt.Println(valueField)
	// output: person
}

func ExampleParser_String() {
	// Create a parser and fill it with data parsed from testLoad
	var pr Parser
	pr.LoadFromFile(testLoad)

	// Parse implements String so it is printed with ini form
	fmt.Println(pr)
}

func ExampleParser_SaveToFile() {
	// Create new parser
	var pr Parser

	// load data from testLoad file to the parser
	pr.LoadFromFile(testLoad)

	// Add data from testLoad file to testSave file
	pr.SaveToFile(testSave)
}
