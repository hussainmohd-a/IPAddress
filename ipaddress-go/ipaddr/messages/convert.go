package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	path := "src/github.com/seancfoley/ipaddress/ipaddress-go/ipaddr/"
	mappings, err := readPropertiesFile(path + "IPAddressResources.properties")
	if err != nil {
		log.Fatal(err)
	}
	source := writeSourceFile(mappings)
	ioutil.WriteFile(path+"ipaddressresources.go", source, 0644)
}

func writeSourceFile(mappings map[string]string) []byte {
	indexMappings := make(map[string]int)
	indices := make([]int, len(mappings))
	valsArray := make([]string, len(mappings))
	i := 0
	valLen := 0
	// create the mappings from string to index into slice, from slice entry to string index
	for key, val := range mappings {
		indexMappings[key] = i
		indices[i] = valLen
		valsArray[i] = val
		valLen += len(val)
		i++
	}
	// now prepare the source code for each of the three elements, the map, the slice, and the string
	mappingsStr := "\n"
	for key, val := range indexMappings {
		mappingsStr += "`" + key + "`: " + strconv.Itoa(val) + ",\n"
	}
	indicesStr := ""
	for i, val := range indices {
		if i%10 == 0 {
			indicesStr += "\n"
		}
		indicesStr += strconv.Itoa(val) + ","
	}
	strStr := "\n"
	for i, val := range valsArray {
		if i > 0 {
			strStr += "+\n"
		}
		strStr += "`" + val + "`"
	}

	bytes :=
		`package ipaddr

var keyStrMap = map[string]int {` + mappingsStr + `
}

var strIndices = []int{` + indicesStr + `
}

var strVals =` + strStr + `

func lookupStr(key string) (result string, ok bool) {
	var index int
	if index, ok = keyStrMap[key]; ok {
		start, end := strIndices[index], strIndices[index+1]
		result = strVals[start:end]
	}
	return
}
`
	return []byte(bytes)
}

func readPropertiesFile(filename string) (map[string]string, error) {
	config := make(map[string]string)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			firstChar := line[0]
			if firstChar != '#' && firstChar != '=' {
				if divIndex := strings.Index(line, "="); divIndex > 0 && divIndex < len(line)-1 {
					key := line[:divIndex]
					value := line[divIndex+1:]
					config[key] = value
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return config, nil
}
