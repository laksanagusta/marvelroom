package utils

import (
	"testing"
)

func TestToPtr(t *testing.T) {
	// Test int
	num := 42
	numPtr := ToPtr(num)
	if *numPtr != 42 {
		t.Errorf("Expected 42, got %d", *numPtr)
	}

	// Test string
	str := "hello"
	strPtr := ToPtr(str)
	if *strPtr != "hello" {
		t.Errorf("Expected 'hello', got %s", *strPtr)
	}

	// Test bool
	b := true
	bPtr := ToPtr(b)
	if *bPtr != true {
		t.Errorf("Expected true, got %v", *bPtr)
	}

	// Test float64
	f := 3.14
	fPtr := ToPtr(f)
	if *fPtr != 3.14 {
		t.Errorf("Expected 3.14, got %f", *fPtr)
	}
}

func TestFromPtr(t *testing.T) {
	// Test with valid pointer
	num := 42
	numPtr := ToPtr(num)
	result := FromPtr(numPtr)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with nil pointer
	var nilPtr *int
	result = FromPtr(nilPtr)
	if result != 0 {
		t.Errorf("Expected 0 (zero value), got %d", result)
	}

	// Test string with nil pointer
	var nilStrPtr *string
	strResult := FromPtr(nilStrPtr)
	if strResult != "" {
		t.Errorf("Expected empty string, got %s", strResult)
	}
}

func TestSafePtr(t *testing.T) {
	// Test with valid pointer
	num := 42
	numPtr := ToPtr(num)
	result := SafePtr(numPtr, 99)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with nil pointer
	var nilPtr *int
	result = SafePtr(nilPtr, 99)
	if result != 99 {
		t.Errorf("Expected 99 (default value), got %d", result)
	}

	// Test string with nil pointer
	var nilStrPtr *string
	strResult := SafePtr(nilStrPtr, "default")
	if strResult != "default" {
		t.Errorf("Expected 'default', got %s", strResult)
	}
}

func TestIsNilPtr(t *testing.T) {
	// Test with non-nil pointer
	num := 42
	numPtr := ToPtr(num)
	if IsNilPtr(numPtr) {
		t.Error("Expected false, got true")
	}

	// Test with nil pointer
	var nilPtr *int
	if !IsNilPtr(nilPtr) {
		t.Error("Expected true, got false")
	}
}

func TestNilPtr(t *testing.T) {
	nilPtr := NilPtr[int]()
	if nilPtr != nil {
		t.Error("Expected nil, got non-nil pointer")
	}
}
