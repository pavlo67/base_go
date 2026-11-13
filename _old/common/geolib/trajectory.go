package geolib

import (
	"fmt"
	"math"
)

type Trajectory []Point

type TrajectoryPoint struct {
	Point
	Bearing
}

func GetTrajectory(from Point, directions []Direction) Trajectory {
	trajectory := make(Trajectory, len(directions)+1)
	trajectory[0] = from

	for i, direction := range directions {
		to := from.PointAtDirection(direction)
		trajectory[i+1], from = to, to
	}

	return trajectory
}

const onPoints = "on geolib.Trajectory.Points()"

func (t Trajectory) Points(velocity, fps float64) ([]TrajectoryPoint, error) {
	if !(velocity > 0) || math.IsInf(velocity, 1) {
		return nil, fmt.Errorf("wrong velocity: %f / "+onPoints, velocity)
	} else if !(fps > 0) || math.IsInf(fps, 1) {
		return nil, fmt.Errorf("wrong fps: %f / "+onPoints, velocity)
	} else if len(t) < 1 {
		return nil, nil
	} else if len(t) == 1 {
		return []TrajectoryPoint{{t[0], 0}}, nil
	}

	currentDirection := 0
	currentDistance := t[0].DistanceTo(t[1])
	currentBearing := t[0].BearingTo(t[1])

	tp := TrajectoryPoint{t[0], currentBearing}
	var passedDistance float64

	trajectoryPoints := []TrajectoryPoint{tp}

	piece := velocity / fps

	for {
		if passedDistance+piece >= currentDistance {
			if currentDirection >= len(t)-2 {
				return append(
					trajectoryPoints,
					TrajectoryPoint{tp.PointAtDirection(Direction{currentBearing, piece}), currentBearing},
				), nil
			}

			passedDistance = passedDistance + piece - currentDistance

			currentDirection++
			currentDistance = t[currentDirection].DistanceTo(t[currentDirection+1])
			currentBearing = t[currentDirection].BearingTo(t[currentDirection+1])

			if passedDistance < currentDistance {
				tp = TrajectoryPoint{
					t[currentDirection].PointAtDirection(Direction{currentBearing, passedDistance}),
					currentBearing,
				}
			} else {
				// setting "temporary point"
				passedDistance -= piece

				tp = TrajectoryPoint{
					t[currentDirection+1].PointAtDirection(Direction{currentBearing + 180, passedDistance}),
					currentBearing,
				}

				continue
			}

		} else {
			passedDistance += piece
			tp = TrajectoryPoint{
				tp.PointAtDirection(Direction{currentBearing, piece}),
				currentBearing,
			}
		}

		trajectoryPoints = append(trajectoryPoints, tp)
	}
}
