package filter

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type ConditionContract interface {
	Name() string
	Check(data *ffprobe.ProbeData) bool
	GetPossibleConditions() []string
	New() ConditionContract
	SetCondition(condition string)
	GetEntry() fyne.Widget
}

type Filter struct {
	condition ConditionString
	value     string
}

type ConditionString string

var (
	equals          ConditionString = "equals"
	contains        ConditionString = "contains"
	notEquals       ConditionString = "not equals"
	greaterThan     ConditionString = "greater than"
	lessThan        ConditionString = "less than"
	greaterOrEquals ConditionString = "greater or equals"
	lessOrEquals    ConditionString = "less or equals"
)

func FromString(condition string) ConditionString {
	switch condition {
	case "equals":
		return equals
	case "contains":
		return contains
	case "not equals":
		return notEquals
	case "greater than":
		return greaterThan
	case "less than":
		return lessThan
	case "greater or equals":
		return greaterOrEquals
	case "less or equals":
		return lessOrEquals
	default:
		return ""
	}
}

func (f *Filter) checkString(value string) bool {
	switch f.condition {
	case equals:
		return value == f.value
	case contains:
		return strings.Contains(value, f.value)
	case notEquals:
		return value != f.value
	default:
		return false
	}
}

func (f *Filter) valueAsInt() (int, error) {
	return strconv.Atoi(f.value)
}

func (f *Filter) checkInt(value int) bool {
	valueInt, err := f.valueAsInt()

	if err != nil {
		return false
	}

	switch f.condition {
	case equals:
		return value == valueInt
	case greaterThan:
		return value > valueInt
	case lessThan:
		return value < valueInt
	case greaterOrEquals:
		return value >= valueInt
	case lessOrEquals:
		return value <= valueInt
	case notEquals:
		return value != valueInt
	default:
		return false
	}
}

func (f *Filter) checkStringToInt(value string) bool {
	valueInt, err := strconv.Atoi(value)

	if err != nil {
		return false
	}

	return f.checkInt(valueInt)
}
