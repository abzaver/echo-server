package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// Tests for the second example using the Proxy pattern (Car and CarProxy)

func TestCarProxy_DriverOldEnough(t *testing.T) {
	driver := &Driver{Age: 18} // Mock object: Driver with age 18
	carProxy := NewCarProxy(driver)

	// Capture stdout to check the output
	result := captureOutput(func() {
		carProxy.Drive()
	})

	expected := "Car is being driven\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCarProxy_DriverTooYoung(t *testing.T) {
	driver := &Driver{Age: 12} // Mock object: Driver with age 12
	carProxy := NewCarProxy(driver)

	// Capture stdout to check the output
	result := captureOutput(func() {
		carProxy.Drive()
	})

	expected := "Driver too young!\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCarProxy_DriverExactly16(t *testing.T) {
	driver := &Driver{Age: 16} // Mock object: Driver with age 16
	carProxy := NewCarProxy(driver)

	// Capture stdout to check the output
	result := captureOutput(func() {
		carProxy.Drive()
	})

	expected := "Car is being driven\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCarProxy_DriverBorderlineTooYoung(t *testing.T) {
	driver := &Driver{Age: 15} // Mock object: Driver with age 15 (just under the limit)
	carProxy := NewCarProxy(driver)

	// Capture stdout to check the output
	result := captureOutput(func() {
		carProxy.Drive()
	})

	expected := "Driver too young!\n"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// Helper function to capture stdout output
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
	}()
	os.Stdout = w

	f()
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
