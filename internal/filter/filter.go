package filter

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type ConditionContract interface {
	Name() string
	GetPossibleConditions() []string
	New() ConditionContract
	CheckGlobal(data *ffprobe.ProbeData) bool
	CheckStream(data *ffprobe.Stream) bool
	SetCondition(condition string)
	GetEntry() fyne.Widget
}

type ConditionString string

var (
	Equals          ConditionString = "equals"
	Contains        ConditionString = "contains"
	NotEquals       ConditionString = "not equals"
	GreaterThan     ConditionString = "greater than"
	LessThan        ConditionString = "less than"
	GreaterOrEquals ConditionString = "greater or equals"
	LessOrEquals    ConditionString = "less or equals"
)

type Filter struct {
	Condition ConditionString
	Value     string
}

func FromString(condition string) ConditionString {
	switch condition {
	case "equals":
		return Equals
	case "contains":
		return Contains
	case "not equals":
		return NotEquals
	case "greater than":
		return GreaterThan
	case "less than":
		return LessThan
	case "greater or equals":
		return GreaterOrEquals
	case "less or equals":
		return LessOrEquals
	default:
		return ""
	}
}

func (f *Filter) CheckString(value string) bool {
	switch f.Condition {
	case Equals:
		return value == f.Value
	case Contains:
		return strings.Contains(value, f.Value)
	case NotEquals:
		return value != f.Value
	default:
		return false
	}
}

func (f *Filter) ValueAsInt() (int, error) {
	return strconv.Atoi(f.Value)
}

func (f *Filter) CheckInt(value int) bool {
	valueInt, err := f.ValueAsInt()

	if err != nil {
		return false
	}

	switch f.Condition {
	case Equals:
		return value == valueInt
	case GreaterThan:
		return value > valueInt
	case LessThan:
		return value < valueInt
	case GreaterOrEquals:
		return value >= valueInt
	case LessOrEquals:
		return value <= valueInt
	case NotEquals:
		return value != valueInt
	default:
		return false
	}
}

func (f *Filter) CheckStringToInt(value string) bool {
	valueInt, err := strconv.Atoi(value)

	if err != nil {
		return false
	}

	return f.CheckInt(valueInt)
}
