package parser

import (
	"errors"
	"os"
	"strconv"
	s "strings"
)

type Parser struct {
	mp map[string]Section
}

type Section map[string]string

const (
	commentChar     = ";"
	newSectionChar  = "["
	endSectionChar  = "]"
	keyValSeparator = "="
)

/*
LoadFromString function converts a string that follows the ini format and
convert it to a map of the section Name to a map of key-value pairs.
the function return an error message with any line that doesn't follow ini format.
*/
func (pr *Parser) LoadFromString(content string) error {
	return pr.parseData(content, "", false)
}

/*
LoadFromFile function converts an ini file into a map of the section Name-map of key-value pairs
the function returns error message if the file doesn't exist or have invalid format
*/
func (pr *Parser) LoadFromFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if check(err) {
		return errors.New("The file in the path\"" + filePath + "\" is not found!")
	}
	return pr.parseData(string(content), filePath, true)
}

/*
GetSectionNames returns a slice contains all sections in the parser map
*/
func (pr Parser) GetSectionNames() []string {
	length := len(pr.mp)
	sectionNames := make([]string, length)
	i := 0
	for key := range pr.mp {
		sectionNames[i] = key
		i++
	}
	return sectionNames
}

func (pr Parser) GetSections() map[string]Section {
	return pr.mp
}

func (pr Parser) Get(sectionName, key string) (string, error) {
	if pr.mp[sectionName] == nil {
		return "", errors.New("Section " + sectionName + " not found")
	}
	if pr.mp[sectionName][key] == "" {
		return "", errors.New("No value found for the key " + key)
	}
	return pr.mp[sectionName][key], nil
}

func (pr Parser) Set(sectionName, key, val string) {
	pr.mp[sectionName][key] = val
}

func (pr Parser) String() string {
	stringValue := ""
	for sec, secVal := range pr.mp {
		stringValue += "[" + sec + "]" + "\n"
		for key, val := range secVal {
			stringValue += key + " : " + val + "\n"
		}
	}
	return stringValue
}

func (pr Parser) SaveToFile(filePath string) {
	file, _ := os.Create(filePath)
	defer file.Close()
	file.Write([]byte(pr.String()))
}

func check(e error) bool {
	return e != nil
}

func (pr *Parser) clearParser() {
	for k := range pr.mp {
		delete(pr.mp, k)
	}
}

func (pr Parser) generateError(errorMsg, filePath, content string, isFile bool) error {
	pr.clearParser()
	var tp string
	if isFile {
		tp = filePath
	} else {
		tp = content
	}

	if isFile {
		return errors.New(errorMsg + " in the file " + tp)
	}
	return errors.New(errorMsg + " in the string\n" + tp)

}

func (pr *Parser) parseData(content, filePath string, isFile bool) error {
	pr.mp = make(map[string]Section) //initialize the parser.
	var newSection Section
	var sectionName string
	lines := s.Split(content, "\n")

	//Go through the text line by line.
	for index, line := range lines {
		s.Trim(line, " ")

		//skip comments and empty lines.
		if s.Index(line, commentChar) == 0 || line == "" {
			continue
		}

		//check if a line contains a comment in the middle of it line and skip it.
		actualLine := s.Split(line, commentChar)[0]

		//check if the line has no key.
		if s.Index(actualLine, keyValSeparator) == 0 {
			return pr.generateError("Key not found in line "+strconv.Itoa(index+1), filePath, content, isFile)
		} else if s.Contains(actualLine, newSectionChar) { //check if the line is a start of a new section.
			//check if there is a section that is previously processed and save it.
			if sectionName != "" {
				pr.mp[sectionName] = newSection
			}

			newSection = make(Section)
			sectionName = s.Split(actualLine, newSectionChar)[1]
			temp := s.Split(sectionName, endSectionChar)
			if !s.Contains(sectionName, endSectionChar) {
				return pr.generateError("Invalid section in line "+strconv.Itoa(index+1), filePath, content, isFile)
			} else if len(temp) > 1 {
				temp[1] = s.Trim(temp[1], " ")
				if temp[1] != "" {
					return pr.generateError("Too much data for the section name in line "+strconv.Itoa(index+1), filePath, content, isFile)
				}
			}
			sectionName = s.Split(sectionName, endSectionChar)[0]
			sectionName = s.Trim(sectionName, " ")

		} else {
			//check if the line is actually in key = value format
			temp := s.Split(actualLine, keyValSeparator)
			key := s.Trim(temp[0], " ")
			val := ""
			if len(temp) == 2 {
				val = s.Trim(temp[1], " ")
			} else if len(temp) > 2 {
				pr.clearParser()
				return pr.generateError("Too much values for one key in line "+strconv.Itoa(index+1), filePath, content, isFile)
			}
			newSection[key] = val
		}
	}
	if sectionName != "" {
		pr.mp[sectionName] = newSection
	}
	return nil

}
