package jsonfilter

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"strconv"
	"strings"

	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	"github.com/ohler55/ojg/jp"
)

type Condition string

const (
	Equals          Condition = "equals"
	Contains        Condition = "contains"
	NotEquals       Condition = "not equals"
	GreaterThan     Condition = "greater than"
	LessThan        Condition = "less than"
	GreaterOrEquals Condition = "greater or equals"
	LessOrEquals    Condition = "less or equals"
)

type FilterType string

const (
	String FilterType = "string"
	Int    FilterType = "int"
	Bool   FilterType = "bool"
)

type Filters struct {
	Filters []filter `json:"filters"`
}

type Filter interface {
	GetType() FilterType
	GetStringCondition() []string
	HasDefaultValues() bool
	GetDefaultValues() []string
	Check(data *helper.FileMetadata, condition Condition, value string) bool
}

type filter struct {
	Name          string      `json:"name"`
	Conditions    []Condition `json:"conditions"`
	JsonPath      string      `json:"jsonPath"`
	Type          FilterType  `json:"type"`
	DefaultValues []string    `json:"values,omitempty" default:"[]"`
}

type Parser interface {
	Parse() (*Filters, error)
}

type parser struct {
	folder embed.FS
}

func NewParser(folder embed.FS) Parser {
	return &parser{folder: folder}
}

func (p parser) Parse() (*Filters, error) {
	dir, err := fs.ReadDir(p.folder, "filters")

	if err != nil {
		return nil, err
	}

	// open jsonpath file and parse json
	var filters []filter = make([]filter, len(dir))

	for i, file := range dir {
		if file.IsDir() {
			continue
		}
		// parse the file
		f, err := p.parseFile("filters/" + file.Name())
		if err != nil {
			return nil, err
		}

		filters[i] = f
	}

	return &Filters{
		Filters: filters,
	}, nil

}

func (p parser) parseFile(path string) (filter, error) {
	file, err := p.folder.ReadFile(path)
	if err != nil {
		return filter{}, err
	}

	// unmarshal json
	var f filter
	err = json.Unmarshal(file, &f)
	if err != nil {
		return filter{}, err
	}

	return f, nil
}

func (f filter) HasDefaultValues() bool {
	return len(f.DefaultValues) > 0
}

func (f filter) GetDefaultValues() []string {
	return f.DefaultValues
}

func (f filter) GetType() FilterType {
	return f.Type
}

func parseInt(value string) (int, error) {
	return strconv.Atoi(value)
}

func parseBool(value string) (bool, error) {
	return strconv.ParseBool(value)
}

// example of jsonPath : $.Format.BitRate
// example of jsonPath : $.Streams[0].Tags.Language
// where $ is the root of the json
// and . is the separator
func (f filter) Check(data *helper.FileMetadata, condition Condition, value string) bool {

	// Convertit les données en JSON

	// Compile le chemin JSONPath
	parse, err := jp.ParseString(f.JsonPath)

	if err != nil {
		log.Default().Println(err)
		return false
	}

	results := parse.Get(data)

	if len(results) == 0 {
		return false
	}

	// Convertit le résultat en chaîne de caractères pour comparaison

	for _, result := range results {

		switch condition {
		case Equals:
			switch f.Type {
			case Int:
				intValue, err := parseInt(value)
				if err != nil {
					log.Default().Println(err)
					return false
				}

				if result == intValue {
					return true
				}
			case Bool:
				boolValue, err := parseBool(value)
				if err != nil {
					log.Default().Println(err)
					return false
				}

				if result == boolValue {
					return true
				}
			}

			if fmt.Sprintf("%v", result) == value {
				return true
			}
		case Contains:
			switch f.Type {
			case String:
				return strings.Contains(fmt.Sprintf("%v", result), value)
			}
			return false
		case NotEquals:
			switch f.Type {
			case Int:
				intValue, err := parseInt(value)
				if err != nil {
					log.Default().Println(err)
					return false
				}

				if result != intValue {
					return true
				}
			case Bool:
				boolValue, err := parseBool(value)
				if err != nil {
					log.Default().Println(err)
					return false
				}

				if result != boolValue {
					return true
				}
			}

			if fmt.Sprintf("%v", result) != value {
				return true
			}
		case GreaterThan:
			if f.Type == Int {
				// parse value to int
				intValue, err := strconv.Atoi(value)
				if err != nil {
					log.Default().Println(err)
					return false
				}
				if result.(int) > intValue {
					return true
				}
			}
			return false
		case LessThan:
			if f.Type == Int {
				// parse value to int
				intValue, err := strconv.Atoi(value)
				if err != nil {
					log.Default().Println(err)
					return false
				}
				if result.(int) < intValue {
					return true
				}
			}
			return false
		case GreaterOrEquals:
			if f.Type == Int {
				// parse value to int
				intValue, err := strconv.Atoi(value)
				if err != nil {
					log.Default().Println(err)
					return false
				}
				if result.(int) >= intValue {
					return true
				}
			}
			return false
		case LessOrEquals:
			if f.Type == Int {
				// parse value to int
				intValue, err := strconv.Atoi(value)
				if err != nil {
					log.Default().Println(err)
					return false
				}
				if result.(int) <= intValue {
					return true
				}
			}
			return false

		}
	}

	return false

}

func (f filter) GetStringCondition() []string {
	var conditions []string
	for _, c := range f.Conditions {
		conditions = append(conditions, string(c))
	}
	return conditions
}
