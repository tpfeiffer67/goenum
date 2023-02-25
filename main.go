package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexkappa/mustache"
)

const (
	goenumDataFolder    = "goenum" + string(os.PathSeparator)
	goenumFileExtension = ".goenum"
	templateFileName    = "goenum.template"
)

type Enum struct {
	Name    string
	Value   string
	IsFirst bool
}

func main() {
	packageName, err := getPackageName()
	check(err)

	template, err := readFileIntoString(goenumDataFolder + templateFileName)
	check(err)

	goenumFiles, err := filepath.Glob(goenumDataFolder + "*" + goenumFileExtension)
	check(err)

	for _, file := range goenumFiles {
		fileName := filepath.Base(file)
		typeName := buildTypeNameFromGoenumFileName(fileName)
		outputFileName := buildEnumGoFileNameFromTypeName(typeName)
		processEnumDescriptionFile(file, packageName, typeName, template, outputFileName)
	}
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func getPackageName() (string, error) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		return "", fmt.Errorf("Error list files *.go: %w", err)
	}

	for _, file := range files {
		s := readPackageNameFromGoFile(file)
		if s != "" {
			return s, nil
		}
	}

	return "", errors.New("Package name not found")
}

func readPackageNameFromGoFile(fileName string) (packageName string) {
	file, err := os.Open(fileName)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "package") {
			fields := strings.Fields(s)
			if len(fields) > 1 {
				return fields[1]
			}
		}
	}
	return ""
}

func buildTypeNameFromGoenumFileName(fileName string) string {
	// example : Alignment.goenum -> Alignment
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func buildEnumGoFileNameFromTypeName(typeName string) string {
	// example : Alignment -> enumalignment.go
	return "enum" + strings.ToLower(typeName) + ".go"
}

func processEnumDescriptionFile(inputFileName string, packageName string, typeName string, template string, outputFileName string) {
	enums, err := readEnumListFromFile(inputFileName)
	if err != nil {
		log.Println(err)
		return
	}

	var m map[string]interface{}
	m = make(map[string]interface{})
	m["Package"] = packageName
	m["EnumType"] = typeName
	m["EnumValues"] = enums
	m["EnumCount"] = len(enums)
	m["EnumLastValue"] = enums[len(enums)-1].Name

	output, err := processMustache(template, m)
	if err != nil {
		log.Println(err)
		return
	}
	err = writeStringToFile(outputFileName, output)
	if err != nil {
		log.Println(err)
	}
}

func readEnumListFromFile(fileName string) (enums []Enum, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Error read enums: %w", err)
	}
	defer file.Close()

	enums = make([]Enum, 0)

	scanner := bufio.NewScanner(file)
	var first bool = true
	var enum Enum
	var f1, f2 string
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		f1 = strings.TrimSpace(fields[0])
		if len(fields) > 1 {
			f2 = strings.TrimSpace(fields[1])
		} else {
			f2 = ""
		}
		enum.Name = f1
		enum.Value = f2
		enum.IsFirst = first
		if first {
			first = false
		}
		enums = append(enums, enum)
	}

	if err := scanner.Err(); err != nil {
		if err != nil {
			return nil, fmt.Errorf("Error read enums: %w", err)
		}
	}

	return enums, err
}

// https://mustache.github.io/mustache.5.html
func processMustache(template string, m map[string]interface{}) (result string, err error) {
	moustache := mustache.New() // moustache = mustache in french :-)

	err = moustache.ParseString(template)
	if err != nil {
		return "", fmt.Errorf("Error Mustache parse: %w", err)
	}

	result, err = moustache.RenderString(m)
	if err != nil {
		return "", fmt.Errorf("Error Mustache render: %w", err)
	}

	return result, err
}
