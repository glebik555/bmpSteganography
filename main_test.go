package main

import (
	"github.com/fatih/color"
	"reflect"
	"testing"
)

func TestSendEnc(t *testing.T) {
	testTable := []struct {
		message  []uint
		expected []uint
	}{
		{
			message:  []uint{1, 0, 1, 1, 1, 1, 1, 1, 1},
			expected: []uint{1, 0, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			message:  []uint{0, 0, 0, 0, 1},
			expected: []uint{0, 0, 0, 0, 1},
		},
		{
			message:  []uint{1, 1},
			expected: []uint{1, 1},
		},
	}

	for _, testCase := range testTable {
		img, pixels := selectFile("pic\\normal.bmp")
		startTranmission(testCase.message, img, pixels, "pic\\outimage.bmp")
		result := startDecode("pic\\outimage.bmp")
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Incorrect result. Result: %q \n Expected: %q", result, testCase.expected)
		} else {
			color.Green("Ok")
		}
	}
}
