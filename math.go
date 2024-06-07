package d

import "math"

type Math struct{}

// 按照保留参数中指定的小数位数返回浮点数。
func (m Math) RoundFloats(decimalPlaces int, value float64) float64 {
	return math.Round((value*math.Pow(10, float64(decimalPlaces+1)))/10) / math.Pow(10, float64(decimalPlaces))
}

// addAndRoundFloats 接收多个浮点数和保留小数位数的参数，将这些浮点数相加，
// 然后按照保留参数中指定的小数位数返回浮点数。
func (m Math) AddFloats(decimalPlaces int, nums ...float64) float64 {
	var sum float64
	for _, num := range nums {
		sum += num
	}
	return m.RoundFloats(decimalPlaces, sum)
}

// subtractFloats 接收多个浮点数，将第一个浮点数减去后续的浮点数，返回结果。
func (m Math) SubtractFloats(decimalPlaces int, nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	result := nums[0]
	for i := 1; i < len(nums); i++ {
		result -= nums[i]
	}
	return m.RoundFloats(decimalPlaces, result)
}

// multiplyFloats 接收多个浮点数，将这些浮点数相乘，返回结果。
func (m Math) MultiplyFloats(decimalPlaces int, nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	result := 1.0
	for _, num := range nums {
		result *= num
	}
	return m.RoundFloats(decimalPlaces, result)
}

// divideFloats 接收两个浮点数，将第一个浮点数除以第二个浮点数，返回结果。
func (m Math) DivideFloats(decimalPlaces int, a, b float64) float64 {
	if b == 0 {
		return math.Inf(1) // 正无穷大
	}
	return m.RoundFloats(decimalPlaces, a/b)
}
