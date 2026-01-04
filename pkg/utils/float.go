// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package utils

func AverageFloat32(numbers []float32) float32 {
	if len(numbers) == 0 {
		return 0
	}

	sum := float32(0)
	for _, num := range numbers {
		sum += num
	}

	return sum / float32(len(numbers))
}
