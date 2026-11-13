package numbers

import (
	"math"
)

func Divide(total int, minPart, maxPart, minLastRatio float64) (num, part, partLast int) {
	minNum, maxNum := float64(total)/maxPart, float64(total)/minPart
	minNumInt, maxNumInt := int(math.Floor(minNum)), int(math.Floor(maxNum))
	if minNumInt > maxNumInt {
		// TODO??? signal the anomaly
		minNumInt, maxNumInt = maxNumInt, minNumInt
	}

	//log.Print(minNum, maxNum, minNumInt, maxNumInt)

	minRest, maxRest := total, 0
	var minRestNum, maxRestNum int
	for k := minNumInt; k <= maxNumInt; k++ {
		rest := total % k
		if rest == 0 {
			part := total / k
			return k, part, part
		}
		if rest < minRest {
			minRest = rest
			minRestNum = k
		}
		if rest > maxRest {
			maxRest = rest
			maxRestNum = k
		}
	}

	part = total / maxRestNum

	// log.Print(maxRestNum, maxRest, float64(maxRest)/float64(part))
	if lastRatio := float64(maxRest) / float64(part); lastRatio >= minLastRatio {
		return maxRestNum + 1, part, maxRest
	}

	part = total / minRestNum
	return minRestNum, part, part + minRest
}
