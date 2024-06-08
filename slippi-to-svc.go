package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	slippi "github.com/pmcca/go-slippi"
)

/*	- Work through csv mapping
	- Read all slp from /slippi
	- Write all csv to /output
	- Concurrency
*/

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	filePath := "Day 3-Game_20210718T094500.slp"
	game, err := slippi.ParseGame(filePath)
	if err != nil {
		log.Fatal(err)
	}

	gameToJSON(game)
	gameToCSV(game)

}

// Write .slp to .json
func gameToJSON(g interface{}) {
	jd, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	file, err := os.Create("test.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jd)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}
}

// Write .slp to .csv -- WIP
func gameToCSV(g interface{}) error {
	file, err := os.Create("frames.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	jd, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}

	var d map[string]interface{}
	if err := json.Unmarshal(jd, &d); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	data := d["Data"].(map[string]interface{})
	frames := data["Frames"].(map[string]interface{})

	// Get All Frame Nums
	frameMap := []int{}
	for key := range frames {
		i, _ := strconv.Atoi(key)
		frameMap = append(frameMap, i)
	}
	sort.Ints(frameMap)

	frameNumber := "0"
	sample := frames[frameNumber].(map[string]interface{})
	flatMap := make(map[string]string)
	flattenMap(sample, "", flatMap)

	keys, err := sortedNumericKeys(flatMap)
	if err != nil {
		return err
	}

	k := append([]string{"Index"}, keys...)
	if err := writer.Write(k); err != nil {
		return err
	}

	// Everything to here is keep

	// This is iterating beyond the elements of the list (by one) ...
	for _, j := range frameMap {
		frameNumber = strconv.Itoa(j)

		frame := frames[frameNumber].(map[string]interface{})
		flatMap := make(map[string]string)
		flattenMap(frame, "", flatMap)

		row := []string{}
		row = append([]string{frameNumber}, row...)
		for _, j := range keys {
			row = append(row, flatMap[j])
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil

}

/*
flattenMap is a recursive function that flattens a nested map into a flat map with composite keys.
nestedMap: The
*/
func flattenMap(nestedMap map[string]interface{}, prefix string, flatMap map[string]string) {
	for key, value := range nestedMap {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			flattenMap(v, fullKey, flatMap)
		case []interface{}:
			for i, item := range v {
				itemKey := fmt.Sprintf("%s.%d", fullKey, i)
				if itemMap, ok := item.(map[string]interface{}); ok {
					flattenMap(itemMap, itemKey, flatMap)
				} else {
					flatMap[itemKey] = fmt.Sprintf("%v", item)
				}
			}
		default:
			flatMap[fullKey] = fmt.Sprintf("%v", v)
		}
	}
}

// TODO: This is not returning in the correct order: (-1, -10, -100, ... 0, 10, 100)
func sortedNumericKeys(m map[string]string) ([]string, error) {
	keys := make([]string, 0, len(m))
	positiveKeys := []string{}
	negativeKeys := []string{}

	for key := range m {
		if strings.HasPrefix(key, "-") {
			negativeKeys = append(negativeKeys, key)
		} else {
			positiveKeys = append(positiveKeys, key)
		}
	}

	sort.Strings(positiveKeys)
	sort.Strings(negativeKeys)

	// Concatenate negative keys followed by positive keys
	keys = append(keys, negativeKeys...)
	keys = append(keys, positiveKeys...)

	return keys, nil
}
