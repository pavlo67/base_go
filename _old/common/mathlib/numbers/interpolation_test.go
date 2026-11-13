package numbers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterpolateByTable(t *testing.T) {
	table := [][2]float64{{1, 10}, {-1, 2}, {5, 20}}

	x, err := InterpolateByTable(-1, table)
	require.NoError(t, err)
	require.Equal(t, 2., x)

	x, err = InterpolateByTable(-2, table)
	require.Error(t, err)
	require.Equal(t, 0., x)

	x, err = InterpolateByTable(0, table)
	require.NoError(t, err)
	require.Equal(t, 6., x)

	x, err = InterpolateByTable(2, table)
	require.NoError(t, err)
	require.Equal(t, 12.5, x)

	x, err = InterpolateByTable(10, table)
	require.Error(t, err)
	require.Equal(t, 0., x)

}

func TestInterpolateByTwoPoints(t *testing.T) {
	tp1 := [2][2]float64{{1, 10}, {-1, 2}}

	x, err := InterpolateByTwoPoints(-1, tp1)
	require.NoError(t, err)
	require.Equal(t, 2., x)

	x, err = InterpolateByTwoPoints(1, tp1)
	require.NoError(t, err)
	require.Equal(t, 10., x)

	x, err = InterpolateByTwoPoints(0, tp1)
	require.NoError(t, err)
	require.Equal(t, 6., x)

	x, err = InterpolateByTwoPoints(2, tp1)
	require.Error(t, err)

	tp2 := [2][2]float64{{-1, 2}, {5, 20}}

	x, err = InterpolateByTwoPoints(-1, tp2)
	require.NoError(t, err)
	require.Equal(t, 2., x)

	x, err = InterpolateByTwoPoints(5, tp2)
	require.NoError(t, err)
	require.Equal(t, 20., x)

	x, err = InterpolateByTwoPoints(0, tp2)
	require.NoError(t, err)
	require.Equal(t, 5., x)

	x, err = InterpolateByTwoPoints(-2, tp2)
	require.Error(t, err)

	x, err = InterpolateByTwoPoints(20, tp2)
	require.Error(t, err)

	tp3 := [2][2]float64{{5, 20}, {5, 21}}

	x, err = InterpolateByTwoPoints(-10, tp3)
	require.Error(t, err)

	x, err = InterpolateByTwoPoints(-5, tp3)
	require.Error(t, err)

	x, err = InterpolateByTwoPoints(6, tp3)
	require.Error(t, err)

}
