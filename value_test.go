package redisson

import (
	"math"
	"testing"
)

const strValue = "test string"
const intValue = 1
const floatValue = 3.1415
const boolValue = true

func TestValue(t *testing.T) {
	value := NewValue("")
	if !value.IsEmpty() {
		t.Errorf("expected empty value")
	}
	if value.V() == nil {
		t.Errorf("expected non-null value, received '%v'", value.V())
	}
}
func TestStringValue(t *testing.T) {
	value := NewValue(strValue)
	if value.IsEmpty() {
		t.Errorf("expected non-empty value")
	}
	if value.AsString() != strValue {
		t.Errorf("expected '%s', received '%s'", strValue, value)
	}
	if value.AsInt() != math.MinInt {
		t.Errorf("expected '%d', received '%d'", math.MinInt, value.AsInt())
	}
	if value.AsFloat() != math.MinInt {
		t.Errorf("expected '%d', received '%f'", math.MinInt, value.AsFloat())
	}
	if value.AsBool() != false {
		t.Errorf("expected '%v', received '%v'", false, value.AsBool())
	}
}
func TestIntValue(t *testing.T) {
	value := NewValue(intValue)
	if value.IsEmpty() {
		t.Errorf("expected non-empty value")
	}
	if value.AsInt() != intValue {
		t.Errorf("expected '%d', received '%d'", intValue, value.AsInt())
	}
}
func TestFloatValue(t *testing.T) {
	value := NewValue(floatValue)
	if value.IsEmpty() {
		t.Errorf("expected non-empty value")
	}
	if value.AsFloat() != floatValue {
		t.Errorf("expected '%f', received '%f'", floatValue, value.AsFloat())
	}
}
func TestBoolValue(t *testing.T) {
	value := NewValue(boolValue)
	if value.IsEmpty() {
		t.Errorf("expected non-empty value")
	}
	if value.AsBool() != boolValue {
		t.Errorf("expected '%v', received '%v'", boolValue, value.AsBool())
	}

	value = NewValue("")
	if value.AsBool() != false {
		t.Errorf("expected '%v', received '%v'", false, value.AsBool())
	}
}
