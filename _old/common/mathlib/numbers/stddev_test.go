package numbers

import (
	"math"
	"testing"
)

func TestStdDev(t *testing.T) {
	type data []float64
	tests := []struct {
		name string
		args data
		want float64
	}{
		{name: "test empty list", args: data{}, want: 0},
		{name: "test 1-element list", args: data{1}, want: 0},
		{name: "test common list", args: data{1, 3, 5}, want: math.Sqrt(8. / 3.)},
		{name: "test sample data list", args: data{1.2219, 1.225, 1.2216, 1.2221, 1.2292}, want: 0.0028917814578560746},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StdDev(tt.args); got != tt.want {
				t.Errorf("HullStdDev() = %v, wantDistance %v", got, tt.want)
			}
		})
	}
}
