package sets

func SignatureCyclicInt(values []int) []int {
	if len(values) < 1 {
		return nil
	}

	maxII := []int{0}
	for i := 1; i < len(values); i++ {
		if values[i] > values[maxII[0]] {
			maxII = []int{i}
		} else if values[i] == values[maxII[0]] {
			maxII = append(maxII, i)
		}
	}

	maxI := maxII[0]

M:
	for m := 1; m < len(maxII); m++ {
		mI := maxII[m]

		for i := 1; i < len(values); i++ {
			if values[(maxI+i)%len(values)] > values[(mI+i)%len(values)] {
				continue M
			} else if values[(maxI+i)%len(values)] < values[(mI+i)%len(values)] {
				maxI = mI
				break
			}
		}
	}

	signature := make([]int, len(values))

	for i := range values {
		signature[i] = values[(maxI+i)%len(values)]
	}

	return signature
}
