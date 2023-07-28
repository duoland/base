package gps

import "testing"

func TestWgs84ToBaidu09(t *testing.T) {
	// longitude, latitude
	gnss := []float64{12122.841797, 3052.433838}
	bd := Wgs84toBaidu09(gnss[0]/100.0, gnss[1]/100.0)
	t.Logf("lon/lat= %f,%f", bd[0], bd[1])
}
