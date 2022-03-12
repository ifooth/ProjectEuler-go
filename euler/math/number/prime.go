package number

import (
	"context"
	"fmt"
	"math"

	"github.com/ifooth/projecteuler-go/euler/math/itertools"
)

const (
	FirstPrime  = int64(2)
	SecondPrime = int64(3)
)

// IsPrime : 素数检测法
// Simple methods https://en.wikipedia.org/wiki/Primality_test
func IsPrime(num int64) bool {
	if num < 2 {
		return false
	}
	if num%2 == 0 && num != 2 {
		return false
	}
	if num%3 == 0 && num != 3 {
		return false
	}

	// +2序列必须加1, int向上浮动1
	for i := int64(5); i <= SqrtInt(num); i += 6 {
		if num%i == 0 || num%(i+2) == 0 {
			return false
		}
	}

	return true
}

type PrimeGeneratorOption struct {
	Start   int64
	End     int64
	Reverse bool
}

// PrimeGenerator : 素数迭代器
func PrimeGenerator(ctx context.Context, opts ...PrimeGeneratorOption) <-chan int64 {
	result := make(chan int64)
	go func() {
		defer close(result)
		result <- FirstPrime
		result <- SecondPrime

		num := SecondPrime
		for {
			num += 2
			if IsPrime(num) {
				select {
				case <-ctx.Done():
					return
				case result <- num:
				}
			}

		}
	}()
	return result
}

// 因子生成器
// 12 = 1 * 2 * 2 * 3
func FactorsGenerator(num int64) <-chan int64 {
	result := make(chan int64)

	go func() {
		defer close(result)

		result <- 1
		factor, limit := FirstPrime, math.Sqrt(float64(num))
		for float64(factor) <= limit {
			if num%factor == 0 {
				result <- factor

				num /= factor
				limit = math.Sqrt(float64(num))
			} else {
				factor += 1
			}
		}
		if num > 1 {
			result <- num
		}
	}()
	return result
}

// Factors 因数分解 12 = 1^1 * 2^2 * 3^1
// return {1: 1, 2: 2, : 3: 1}
func Factors(num int64) map[int64]int64 {
	factorMap := map[int64]int64{}
	for factor := range FactorsGenerator(num) {
		factorMap[factor] += 1
	}
	return factorMap
}

// ProperDivisors 真因子 小于n且整除n的正整数, 不包含自己
// 12 = [1, 2, 3, 4, 6]
func ProperDivisors(num int64) []int64 {
	factorMap := Factors(num)
	divisorSet := map[int64]struct{}{1: {}}

	fmt.Println(factorMap)

	limit := num / 2

	for prime, power := range factorMap {
		fmt.Println("prime, power", prime, power)
		for p := range itertools.IterInt(power + 1) {
			temp := map[int64]struct{}{}
			for d := range divisorSet {
				divisor := d * PowInt(prime, p)
				if divisor > limit {
					continue
				}
				temp[divisor] = struct{}{}
				fmt.Println("divisorSet", p, d, divisor, temp)
			}
			divisorSet = temp
		}
	}

	divisors := make([]int64, len(divisorSet))
	for d := range divisorSet {
		divisors = append(divisors, d)
	}
	return divisors
}
