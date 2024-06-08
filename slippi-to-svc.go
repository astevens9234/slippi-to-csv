// Convert .slp files to .csv or .json files for downstream processing and Deep Learning.
//
// This package contains two callables;
//   - gameToJSON
//   - gameToCSV
//
// .slp files will be read from /slp and written to either /csv or /json depending on function call.
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

/*	TODO:
	- Concurrency
	- Add directory structure
	- ...
*/

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	bulkProcessing()

	// Create a single file -- reference
	// filePath := "slp/Day 3-Game_20210718T094500.slp"
	// game, err := slippi.ParseGame(filePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// gameToJSON(game, "json/Sample.json")
	// gameToCSV(game, "csv/Sample.csv")

}

// Write .slp to .json
func gameToJSON(g interface{}, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jd, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		panic(err)
	}
	_, err = file.Write(jd)
	if err != nil {
		panic(err)
	}
}

// Write .slp to .csv
func gameToCSV(g interface{}, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	jd, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		panic(err)
	}

	var d map[string]interface{}
	if err := json.Unmarshal(jd, &d); err != nil {
		panic(err)
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

	keys, err := sortedKeys(flatMap)
	if err != nil {
		return err
	}

	k := append([]string{"Index"}, keys...)
	if err := writer.Write(k); err != nil {
		return err
	}

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

// flattenMap is a recursive function that flattens a nested maps into a flat map with composite keys.
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

func sortedKeys(m map[string]string) ([]string, error) {
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

	keys = append(keys, negativeKeys...)
	keys = append(keys, positiveKeys...)

	return keys, nil
}

func bulkProcessing() {
	f, _ := os.ReadDir("slp")
	for _, j := range f {
		// This anonymous function allows the program to skip over corrupted files
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("panic occured: ", r)
				}
			}()

			filePath := "slp/" + j.Name()
			fileName := "csv/" + j.Name()[:len(j.Name())-4] + ".csv"

			game, err := slippi.ParseGame(filePath)
			if err != nil {
				panic(err)
			}

			err = gameToCSV(game, fileName)
			if err != nil {
				panic(err)
			}
		}()
	}
}
