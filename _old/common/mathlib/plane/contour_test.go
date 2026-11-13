package plane

//func TestContour_Shorten(t *testing.T) {
//	tests := []struct {
//		name        string
//		cntr        Contour
//		distanceMax float64
//		want        Contour
//	}{
//		{
//			name:        "",
//			cntr:        Contour{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
//			distanceMax: 10,
//			want:        nil,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.cntr.Shorten(tt.distanceMax); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Shorten() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestContour_ConvexHull(t *testing.T) {
//	tests := []struct {
//		name string
//		cntr Contour
//		want []image.Point
//	}{
//		{
//			name: "",
//			cntr: Contour{{0, 0}, {1, 0}, {1, 1}},
//			want: []image.Point{{0, 0}, {1, 0}, {1, 1}},
//		},
//		{
//			name: "",
//			cntr: Contour{{0, 0}, {1, 1}, {1, 0}},
//			want: []image.Point{{0, 0}, {1, 0}, {1, 1}},
//		},
//		{
//			name: "",
//			cntr: Contour{{0, 0}, {1, 1}, {0.75, 0.5}, {1, 0}},
//			want: []image.Point{{0, 0}, {1, 0}, {1, 1}},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := tt.cntr.ConvexHull()
//			sort.Slice(got, func(i, j int) bool { return got[i].X < got[j].X || (got[i].X == got[j].X && got[i].Y < got[j].Y) })
//
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ConvexHull() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
