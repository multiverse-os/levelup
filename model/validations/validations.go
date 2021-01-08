package validations

type Validation func(model interface{}) (bool, error)

////////////////////////////////////////////////////////////////////////////////

func IsEmpty(value interface{}) bool {
	switch value.(type) {
	case string:
		return IsZero(len(value.(string)))
	case int:
		return IsZero(value.(int))
	case []rune:
		return IsZero(len(value.([]rune)))
	case []string:
		return IsZero(len(value.([]string)))
	case []byte:
		return IsZero(len(value.([]byte)))
	case []int:
		return IsZero(len(value.([]int)))
	default:
		return false
	}
}

func NotEmpty(value interface{}) bool { return !IsEmpty(value) }

func IsNil(value interface{}) bool {
	switch value.(type) {
	case int:
		return IsZero(value.(int))
	case string:
		return IsBlank(value.(string))
	case error:
		return value.(error) == nil
	default:
		return value == nil
	}
}
func NotNil(value interface{}) bool { return !IsNil(value) }

////////////////////////////////////////////////////////////////////////////////
// VALIDATIONS::String
func IsBlank(str string) bool  { return IsZero(len(str)) }
func NotBlank(str string) bool { return !IsZero(len(str)) }

////////////////////////////////////////////////////////////////////////////////
// VALIDATION::Numeric
func IsZero(value int) bool  { return value == 0 }
func NotZero(value int) bool { return value != 0 }

func IsGreaterThan(gt, value int) bool         { return (gt > value) }
func IsGreaterOrEqualThan(gte, value int) bool { return (gte >= value) }
func IsLessThan(lt, value int) bool            { return (lt < value) }
func IsLessOrEqualThan(lte, value int) bool    { return (lte <= value) }

func IsBetween(start, end, value int) bool  { return start < value && value < end }
func NotBetween(start, end, value int) bool { return start > value && value > end }
