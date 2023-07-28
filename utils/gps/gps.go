package gps

import "math"

const (
	X_PI = 3.14159265358979324 * 3000.0 / 180.0
	PI   = 3.1415926535897932384626
	A    = 6378245.0
	EE   = 0.00669342162296594323
)

func outOfChina(lon, lat float64) bool {
	if (lon < 72.004 || lon > 137.8347) ||
		(lat < 0.8293 || lat > 55.8271) {
		return true
	}
	return false
}

func transformLatitude(lon, lat float64) (ret float64) {
	ret = -100.0 + 2.0*lon + 3.0*lat + 0.2*lat*lat + 0.1*lon*lat + 0.2*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*PI) + 20.0*math.Sin(2.0*lon*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*PI) + 40.0*math.Sin(lat/3.0*PI)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*PI) + 320*math.Sin(lat*PI/30.0)) * 2.0 / 3.0
	return ret
}

func transformLongitude(lon, lat float64) (ret float64) {
	ret = 300.0 + lon + 2.0*lat + 0.1*lon*lon + 0.1*lon*lat + 0.1*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*PI) + 20.0*math.Sin(2.0*lon*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lon*PI) + 40.0*math.Sin(lon/3.0*PI)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lon/12.0*PI) + 300.0*math.Sin(lon/30.0*PI)) * 2.0 / 3.0
	return ret
}

func Wgs84toGcj02(wgsLon, wgsLat float64) (ret []float64) {
	if outOfChina(wgsLon, wgsLat) {
		return []float64{wgsLon, wgsLat}
	}
	dlat := transformLatitude(wgsLon-105.0, wgsLat-35.0)
	dlng := transformLongitude(wgsLon-105.0, wgsLat-35.0)
	radlat := wgsLat / 180.0 * PI
	magic := math.Sin(radlat)
	magic = 1 - EE*magic*magic
	sqrtmagic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((A * (1 - EE)) / (magic * sqrtmagic) * PI)
	dlng = (dlng * 180.0) / (A / sqrtmagic * math.Cos(radlat) * PI)
	mglat := wgsLat + dlat
	mglng := wgsLon + dlng
	return []float64{mglng, mglat}
}

func Gcj02toWgs84(gcjLon, gcjLat float64) (ret []float64) {
	if outOfChina(gcjLon, gcjLat) {
		return []float64{gcjLon, gcjLat}
	}
	dlat := transformLatitude(gcjLon-105.0, gcjLat-35.0)
	dlng := transformLongitude(gcjLon-105.0, gcjLat-35.0)
	radlat := gcjLat / 180.0 * PI
	magic := math.Sin(radlat)
	magic = 1 - EE*magic*magic
	sqrtmagic := math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((A * (1 - EE)) / (magic * sqrtmagic) * PI)
	dlng = (dlng * 180.0) / (A / sqrtmagic * math.Cos(radlat) * PI)
	mglat := gcjLat + dlat
	mglng := gcjLon + dlng
	return []float64{gcjLon*2 - mglng, gcjLat*2 - mglat}
}

func Gcj02toBaidu09(gcjLon, gcjLat float64) []float64 {
	z := math.Sqrt(gcjLon*gcjLon+gcjLat*gcjLat) + 0.00002*math.Sin(gcjLat*X_PI)
	theta := math.Atan2(gcjLat, gcjLon) + 0.000003*math.Cos(gcjLon*X_PI)
	bdLon := z*math.Cos(theta) + 0.0065
	bdLat := z*math.Sin(theta) + 0.006
	return []float64{bdLon, bdLat}
}

func Baidu09toGcj02(bdLon, bdLat float64) (ret []float64) {
	x := bdLon - 0.0065
	y := bdLat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*X_PI)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*X_PI)
	gcjLon := z * math.Cos(theta)
	gcjLat := z * math.Sin(theta)
	return []float64{gcjLon, gcjLat}
}

func Wgs84toBaidu09(wgsLon, wgsLat float64) (ret []float64) {
	gcj := Wgs84toGcj02(wgsLon, wgsLat)
	return Gcj02toBaidu09(gcj[0], gcj[1])
}

func Baidu09toWgs84(bdLon, bdLat float64) (ret []float64) {
	gcj := Baidu09toGcj02(bdLon, bdLat)
	return Gcj02toWgs84(gcj[0], gcj[1])
}
