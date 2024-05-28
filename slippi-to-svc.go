package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	slippi "github.com/pmcca/go-slippi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	filePath := "Day 3-Game_20210718T094500.slp"
	game, err := slippi.ParseGame(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(game.Data.GameStart.GameMode)

	// Write the data to a CSV file
	if err := writeCsv("test.csv", game); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("CSV file written successfully!")
	}
}

func gameToSlice(v interface{}) []string {
	var result []string
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		switch field.Kind() {
		case reflect.String:
			result = append(result, field.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result = append(result, strconv.FormatInt(field.Int(), 10))
		case reflect.Float32, reflect.Float64:
			result = append(result, strconv.FormatFloat(field.Float(), 'f', -1, 64))
		case reflect.Bool:
			result = append(result, strconv.FormatBool(field.Bool()))
		case reflect.Struct:
			result = append(result, gameToSlice(field.Interface())...)
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				result = append(result, fmt.Sprint(field.Index(j)))
			}
		default:
			result = append(result, fmt.Sprint(field))
		}
	}
	return result
}

func writeCsv(filename string, data interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		log.Fatal(val.Kind())
	}

	// Assuming all structs have the same fields and structure, use the first struct for the header
	if val.Len() > 0 {
		firstElement := val.Index(0)
		header := getHeader(firstElement.Interface())
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	// Write the data
	for i := 0; i < val.Len(); i++ {
		row := gameToSlice(val.Index(i).Interface())
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// Helper function to get the header from a struct
func getHeader(v interface{}) []string {
	var header []string
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if field.Type.Kind() == reflect.Struct {
			nestedHeader := getHeader(val.Field(i).Interface())
			header = append(header, nestedHeader...)
		} else {
			header = append(header, field.Name)
		}
	}
	return header
}
