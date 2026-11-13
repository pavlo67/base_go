package numbers

import "testing"

func TestDivide(t *testing.T) {
	tests := []struct {
		name         string
		total        int
		minPart      float64
		maxPart      float64
		minLastRatio float64
		wantNum      int
		wantPart     int
		wantPartLast int
	}{
		{
			name:         "exact",
			total:        25,
			minPart:      4,
			maxPart:      6,
			minLastRatio: 0.5,
			wantNum:      5,
			wantPart:     5,
			wantPartLast: 5,
		},
		{
			name:         "bigRest",
			total:        26,
			minPart:      4,
			maxPart:      6,
			minLastRatio: 0.6,
			wantNum:      5,
			wantPart:     5,
			wantPartLast: 6,
		},
		{
			name:         "smallRest",
			total:        27,
			minPart:      4,
			maxPart:      6,
			minLastRatio: 0.5,
			wantNum:      5,
			wantPart:     6,
			wantPartLast: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNum, gotPart, gotPartLast := Divide(tt.total, tt.minPart, tt.maxPart, tt.minLastRatio)
			if gotNum != tt.wantNum {
				t.Errorf("Divide() gotNum = %v, wantDistance %v", gotNum, tt.wantNum)
			}
			if gotPart != tt.wantPart {
				t.Errorf("Divide() gotPart = %v, wantDistance %v", gotPart, tt.wantPart)
			}
			if gotPartLast != tt.wantPartLast {
				t.Errorf("Divide() gotPartLast = %v, wantDistance %v", gotPartLast, tt.wantPartLast)
			}
		})
	}
}
