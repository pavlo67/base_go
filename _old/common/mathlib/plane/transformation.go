package plane

import (
	"fmt"
	"math"
)

type Transformation struct {
	From Point2
	To   Point2
}

const onCalculateRotationAndScale = "on CalculateRotationAndScale()"

func CalculateRotationAndScale(pts []Transformation, rotDeviationMax, scaleDeviationMax float64) (XToYAngle, float64, error) {
	if len(pts) != 3 {
		return 0, 0, fmt.Errorf("wrong points (%d), must be 3 items exactly / "+onCalculateRotationAndScale, len(pts))
	}

	cOut, cIn := make(Contour, len(pts)), make(Contour, len(pts))

	for i, p := range pts {
		if i > 0 && pts[i-1].From.X == p.From.X && pts[i-1].From.Y == p.From.Y {
			return 0, 0, fmt.Errorf("pts[%d] == pts[%d] / "+onCalculateRotationAndScale, i-1, i)
		}
		cOut[i], cIn[i] = p.To, p.From
	}

	rotInn, distInn := cOut.Rotations(), cOut.Distances()
	rotOut, distOut := cIn.Rotations(), cIn.Distances()

	//fmt.Printf("rotInn: %v\n", rotInn)
	//fmt.Printf("rotOut: %v\n", rotOut)

	var rotDeltaSum XToYAngle
	var scalesSum float64

	rotDeltas := make([]XToYAngle, len(pts))
	scales := make([]float64, len(pts))

	for i := 0; i < len(cOut); i++ {
		rotDeltas[i] = (rotInn[i] - rotOut[i]).Canon()
		rotDeltaSum += rotDeltas[i]
		scales[i] = distInn[i] / distOut[i]
		scalesSum += scales[i]
	}

	rotDeltaAvg := (rotDeltaSum / XToYAngle(len(pts))).Canon()
	scaleAvg := scalesSum / float64(len(pts))

	if scaleAvg <= 0 {
		return 0, 0, fmt.Errorf("wrong scaleAvg (%f) for scales: %v / "+onCalculateRotationAndScale, scaleAvg, scales)
	}

	rotDeltaDev := math.Abs(float64((rotInn[0] - rotOut[0] - rotDeltaAvg).Canon()))
	scaleDev := math.Abs(1 - scales[0]/scaleAvg)
	for i := 1; i < len(pts); i++ {
		rotDeltaDevI := math.Abs(float64((rotInn[i] - rotOut[i] - rotDeltaAvg).Canon()))
		if rotDeltaDevI > rotDeltaDev {
			rotDeltaDev = rotDeltaDevI
		}
		scaleDevI := math.Abs(1 - scales[i]/scaleAvg)
		if scaleDevI > scaleDev {
			scaleDev = scaleDevI
		}
	}

	if rotDeltaDev > rotDeviationMax {
		return 0, 0, fmt.Errorf("rotDeltaDev (%f) > rotDeviationMax (%f) for rotation deltas: %v, rotDeltaAvg: %f / "+onCalculateRotationAndScale,
			rotDeltaDev, rotDeviationMax, rotDeltas, rotDeltaAvg)
	}

	if scaleDev > scaleDeviationMax {
		return 0, 0, fmt.Errorf("scaleDev (%f) > scaleDeviationMax (%f) for scales: %v / "+onCalculateRotationAndScale,
			scaleDev, scaleDeviationMax, scales)
	}

	return rotDeltaAvg, scaleAvg, nil
}
