package gps

import "testing"

func TestWgs84ToBaidu09(t *testing.T) {
	// longitude, latitude
	//gnss := []float64{12122.841797, 3052.433838}
	gnss := []float64{12124.730469, 3104.423584}
	lon_src := gnss[0]
	lat_src := gnss[1]

	// convert gnss to wgs84
	lon_dd := lon_src / 100.0
	lat_dd := lat_src / 100.0

	lon_dd_i := float64(int(lon_dd))
	lat_dd_i := float64(int(lat_dd))

	lon_mm := (lon_dd - lon_dd_i) * 100.0
	lat_mm := (lat_dd - lat_dd_i) * 100.0

	lon_dd = lon_dd_i + lon_mm/60.0
	lat_dd = lat_dd_i + lat_mm/60.0

	// convert wgs84 to baidu09
	bd := Wgs84toBaidu09(lon_dd, lat_dd)
	t.Logf("lon/lat= %f,%f", bd[0], bd[1])
	t.Logf("lat/lon= %f,%f", bd[1], bd[0])
}
