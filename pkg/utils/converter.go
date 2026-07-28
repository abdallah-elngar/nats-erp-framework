package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// ToString يحول إلى نص
func ToString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case time.Time:
		return val.Format(time.RFC3339)
	default:
		data, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		return string(data)
	}
}

// ToInt يحول إلى عدد صحيح
func ToInt(v interface{}) (int, error) {
	if v == nil {
		return 0, nil
	}

	switch val := v.(type) {
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case uint:
		return int(val), nil
	case uint8:
		return int(val), nil
	case uint16:
		return int(val), nil
	case uint32:
		return int(val), nil
	case uint64:
		return int(val), nil
	case float32:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

// ToFloat يحول إلى عدد عشري
func ToFloat(v interface{}) (float64, error) {
	if v == nil {
		return 0, nil
	}

	switch val := v.(type) {
	case float32:
		return float64(val), nil
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float", v)
	}
}

// ToBool يحول إلى منطقي
func ToBool(v interface{}) (bool, error) {
	if v == nil {
		return false, nil
	}

	switch val := v.(type) {
	case bool:
		return val, nil
	case int:
		return val != 0, nil
	case int64:
		return val != 0, nil
	case uint64:
		return val != 0, nil
	case string:
		return strconv.ParseBool(val)
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

// ToTime يحول إلى وقت
func ToTime(v interface{}) (time.Time, error) {
	if v == nil {
		return time.Time{}, nil
	}

	switch val := v.(type) {
	case time.Time:
		return val, nil
	case string:
		return time.Parse(time.RFC3339, val)
	default:
		return time.Time{}, fmt.Errorf("cannot convert %T to time", v)
	}
}

// ToJSON يحول إلى JSON
func ToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// FromJSON يحول من JSON
func FromJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
