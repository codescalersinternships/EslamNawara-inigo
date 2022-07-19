# INI Parser
Simple parser that provides variety of functionalities for reading from and writing to ini files.
 
## Functionalities:
 - `LoadFromString` Initializes the parser with the content of the string argument.
 - `LoadFromFile` It is similar to the `LoadFromString` function but gets its content from a file which its bath is given as argument. 
 - `GetSectionNames` Returns all the sections in the parser.
 - `GetSections` Returns the parser's content converted into a map.
 - `Get` Returns the value with key 'key' and section 'sectionName'.
 - `Set` Sets the value with key 'key' and section 'sectionName' to value 'val'.
 - `ToString` Returns a string copy contains the data stored in the parser.
 - `SaveToFile` Overwrites the file with path 'filePath' with the value of the parser or create a file with similar name if not exist and Write the parser value to it.
**<ins>Note that </ins>** it is reqired to initialize the parser by using one of the functions(LoadFromString, LoadFromFile) before using any of the other functions to make sure that the parser is not empty.

## What are the code modules?
 - `parser.go` Contains the actual implimentation of the parser.
 - `test_parser.go` Contains testers for the parser to make sure every thing is going well.
 - `testLoad.ini` Helper file used to test the parser by reading from it. 
 - `testSave.ini` Helper file used to test the parser by writing to it.

## To use the parser follow the commands:
 ```
 git clone https://github.com/codescalersinternships/EslamNawara-inigo/tree/main/parser
 cd parser
 go test -v
 ```
  After testing is done successfully you can go ahead and use the library by importing it to your go file
