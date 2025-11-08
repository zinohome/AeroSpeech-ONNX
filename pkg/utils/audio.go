package utils

import (
	"encoding/binary"
	"math"
)

// SamplesInt16ToFloat 将int16样本转换为float32样本
func SamplesInt16ToFloat(data []byte) []float32 {
	if len(data)%2 != 0 {
		return nil
	}

	samples := make([]float32, len(data)/2)
	for i := 0; i < len(samples); i++ {
		// 读取little-endian int16
		sample := int16(binary.LittleEndian.Uint16(data[i*2 : i*2+2]))
		// 转换为float32 (-1.0 到 1.0)
		samples[i] = float32(sample) / 32768.0
	}

	return samples
}

// SamplesFloatToInt16 将float32样本转换为int16样本
func SamplesFloatToInt16(samples []float32) []byte {
	data := make([]byte, len(samples)*2)
	for i, sample := range samples {
		// 限制范围到[-1.0, 1.0]
		sample = float32(math.Max(-1.0, math.Min(1.0, float64(sample))))
		// 转换为int16
		sampleInt16 := int16(sample * 32767.0)
		// 写入little-endian
		binary.LittleEndian.PutUint16(data[i*2:i*2+2], uint16(sampleInt16))
	}

	return data
}

// ConvertPCM16ToFloat32 转换PCM16音频数据到float32样本
func ConvertPCM16ToFloat32(data []byte) []float32 {
	return SamplesInt16ToFloat(data)
}

// ConvertFloat32ToPCM16 转换float32样本到PCM16音频数据
func ConvertFloat32ToPCM16(samples []float32) []byte {
	return SamplesFloatToInt16(samples)
}

// ResampleAudio 重采样音频（简化实现，实际应使用专业库）
func ResampleAudio(samples []float32, fromRate, toRate int) []float32 {
	if fromRate == toRate {
		return samples
	}

	ratio := float64(toRate) / float64(fromRate)
	newLength := int(float64(len(samples)) * ratio)
	resampled := make([]float32, newLength)

	for i := 0; i < newLength; i++ {
		srcIndex := float64(i) / ratio
		index := int(srcIndex)
		frac := srcIndex - float64(index)

		if index+1 < len(samples) {
			// 线性插值
			resampled[i] = float32(float64(samples[index])*(1-frac) + float64(samples[index+1])*frac)
		} else if index < len(samples) {
			resampled[i] = samples[index]
		}
	}

	return resampled
}

