package utils

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestSamplesInt16ToFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []float32
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: []float32{},
		},
		{
			name:     "zero sample",
			input:    []byte{0x00, 0x00},
			expected: []float32{0.0},
		},
		{
			name:     "max positive",
			input:    []byte{0xFF, 0x7F},
			expected: []float32{32767.0 / 32768.0},
		},
		{
			name:     "max negative",
			input:    []byte{0x00, 0x80},
			expected: []float32{-32768.0 / 32768.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SamplesInt16ToFloat(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(got))
				return
			}

			for i := range got {
				if math.Abs(float64(got[i]-tt.expected[i])) > 0.0001 {
					t.Errorf("Index %d: expected %f, got %f", i, tt.expected[i], got[i])
				}
			}
		})
	}
}

func TestSamplesFloatToInt16(t *testing.T) {
	tests := []struct {
		name     string
		input    []float32
		check    func([]byte) bool
	}{
		{
			name:  "empty input",
			input: []float32{},
			check: func(b []byte) bool {
				return len(b) == 0
			},
		},
		{
			name:  "zero sample",
			input: []float32{0.0},
			check: func(b []byte) bool {
				if len(b) != 2 {
					return false
				}
				return binary.LittleEndian.Uint16(b) == 0
			},
		},
		{
			name:  "max positive",
			input: []float32{1.0},
			check: func(b []byte) bool {
				if len(b) != 2 {
					return false
				}
				val := int16(binary.LittleEndian.Uint16(b))
				return val > 0 && val <= 32767
			},
		},
		{
			name:  "max negative",
			input: []float32{-1.0},
			check: func(b []byte) bool {
				if len(b) != 2 {
					return false
				}
				val := int16(binary.LittleEndian.Uint16(b))
				return val < 0 && val >= -32768
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SamplesFloatToInt16(tt.input)
			if !tt.check(got) {
				t.Errorf("Check failed for input %v", tt.input)
			}
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	original := []float32{-1.0, -0.5, 0.0, 0.5, 1.0}

	// float32 -> int16 -> float32
	int16Data := SamplesFloatToInt16(original)
	converted := SamplesInt16ToFloat(int16Data)

	if len(converted) != len(original) {
		t.Fatalf("Length mismatch: expected %d, got %d", len(original), len(converted))
	}

	// 允许一定的精度损失
	for i := range original {
		diff := math.Abs(float64(converted[i] - original[i]))
		if diff > 0.01 {
			t.Errorf("Index %d: expected %f, got %f (diff: %f)", i, original[i], converted[i], diff)
		}
	}
}

func TestSamplesInt16ToFloatOddLength(t *testing.T) {
	// 测试奇数长度的输入
	input := []byte{0x00, 0x00, 0xFF}
	result := SamplesInt16ToFloat(input)
	if result != nil {
		t.Error("Expected nil for odd-length input")
	}
}

func TestConvertPCM16ToFloat32(t *testing.T) {
	input := []byte{0x00, 0x00, 0xFF, 0x7F}
	result := ConvertPCM16ToFloat32(input)
	if len(result) != 2 {
		t.Errorf("Expected 2 samples, got %d", len(result))
	}
}

func TestConvertFloat32ToPCM16(t *testing.T) {
	input := []float32{0.0, 1.0, -1.0}
	result := ConvertFloat32ToPCM16(input)
	if len(result) != 6 { // 3 samples * 2 bytes
		t.Errorf("Expected 6 bytes, got %d", len(result))
	}
}

func TestResampleAudio(t *testing.T) {
	tests := []struct {
		name     string
		samples  []float32
		fromRate int
		toRate   int
		check    func([]float32) bool
	}{
		{
			name:     "same rate",
			samples:  []float32{0.0, 0.5, 1.0},
			fromRate: 16000,
			toRate:   16000,
			check: func(result []float32) bool {
				return len(result) == 3
			},
		},
		{
			name:     "upsample",
			samples:  []float32{0.0, 1.0},
			fromRate: 16000,
			toRate:   32000,
			check: func(result []float32) bool {
				return len(result) >= 2 && len(result) <= 4
			},
		},
		{
			name:     "downsample",
			samples:  []float32{0.0, 0.5, 1.0, 0.5, 0.0},
			fromRate: 32000,
			toRate:   16000,
			check: func(result []float32) bool {
				return len(result) >= 2 && len(result) <= 5
			},
		},
		{
			name:     "empty input",
			samples:  []float32{},
			fromRate: 16000,
			toRate:   32000,
			check: func(result []float32) bool {
				return len(result) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResampleAudio(tt.samples, tt.fromRate, tt.toRate)
			if !tt.check(result) {
				t.Errorf("ResampleAudio() check failed for %s", tt.name)
			}
		})
	}
}

func TestSamplesFloatToInt16Clamping(t *testing.T) {
	// 测试超出范围的值会被限制
	input := []float32{-2.0, -1.0, 0.0, 1.0, 2.0}
	result := SamplesFloatToInt16(input)
	
	if len(result) != 10 { // 5 samples * 2 bytes
		t.Errorf("Expected 10 bytes, got %d", len(result))
	}
	
	// 验证值被限制在[-1.0, 1.0]范围内
	converted := SamplesInt16ToFloat(result)
	for i, val := range converted {
		if val < -1.0 || val > 1.0 {
			t.Errorf("Sample %d: value %f is out of range [-1.0, 1.0]", i, val)
		}
	}
}

