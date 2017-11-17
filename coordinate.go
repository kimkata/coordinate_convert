package main

// conversion between WGS-84, GCJ-02, BD-09-LL, BD-09-MC
// conversion structure
//
// BD09LL --- GCJ-02 ----------------- WGS84
//   			|						  |
//		(BaiduMercartorConvert)	 (WebMercartorConvert)
//   			|						  |
//			BD09MC					WebMercator
//
//
// referered documents:
// http://blog.csdn.net/yu412346928/article/details/24419845
// http://ju.outofmemory.cn/entry/309208
// https://github.com/everalan/convertMC2LL
// http://blog.csdn.net/sinat_34719507/article/details/60904361
import (
	"errors"
	"math"
)

var InvalidGPS = errors.New("invalid gps!")

var x_PI = 3.14159265358979324 * 3000.0 / 180.0
var PI = 3.1415926535897932384626
var a = 6378245.0
var ee = 0.00669342162296594323
var earchHalfCir = 20037508.34

func inChina(lng, lat float64) bool {
	// 纬度3.86~53.55,经度73.66~135.05
	return lng > 73.66 && lng < 135.05 && lat > 3.86 && lat < 53.55
}

func IsInValidGps(lng, lat float64) bool {
	return !(lng >= -180 && lng <= 180 && lat >= -90 && lat <= 90)
}

func transLat(lng, lat float64) float64 {
	ret := -100.0 + 2.0*lng + 3.0*lat + 0.2*lat*lat + 0.1*lng*lat + 0.2*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*PI) + 20.0*math.Sin(2.0*lng*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*PI) + 40.0*math.Sin(lat/3.0*PI)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*PI) + 320*math.Sin(lat*PI/30.0)) * 2.0 / 3.0
	return ret
}

func transLng(lng, lat float64) float64 {
	ret := 300.0 + lng + 2.0*lat + 0.1*lng*lng + 0.1*lng*lat + 0.1*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*PI) + 20.0*math.Sin(2.0*lng*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lng*PI) + 40.0*math.Sin(lng/3.0*PI)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lng/12.0*PI) + 300.0*math.Sin(lng/30.0*PI)) * 2.0 / 3.0
	return ret
}

func GCJ02ToWGS84(lng, lat float64) (float64, float64, error) {
	if IsInValidGps(lng, lat) {
		return 0, 0, InvalidGPS
	}

	if !inChina(lng, lat) {
		return lng, lat, nil
	}

	var dlat = transLat(lng-105.0, lat-35.0)
	var dlng = transLng(lng-105.0, lat-35.0)
	var radlat = lat / 180.0 * PI
	var magic = math.Sin(radlat)
	magic = 1 - ee*magic*magic
	var Sqrtmagic = math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * Sqrtmagic) * PI)
	dlng = (dlng * 180.0) / (a / Sqrtmagic * math.Cos(radlat) * PI)
	var mglat = lat + dlat
	var mglng = lng + dlng
	return lng*2 - mglng, lat*2 - mglat, nil
}

func WGS84ToGCJ02(lng, lat float64) (float64, float64, error) {
	if IsInValidGps(lng, lat) {
		return 0, 0, InvalidGPS
	}

	if !inChina(lng, lat) {
		return lng, lat, nil
	}

	var dlat = transLat(lng-105.0, lat-35.0)
	var dlng = transLng(lng-105.0, lat-35.0)
	var radlat = lat / 180.0 * PI
	var magic = math.Sin(radlat)
	magic = 1 - ee*magic*magic
	var sqrtmagic = math.Sqrt(magic)
	dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * sqrtmagic) * PI)
	dlng = (dlng * 180.0) / (a / sqrtmagic * math.Cos(radlat) * PI)
	return lng + dlng, lat + dlat, nil
}

func GCJ02ToBD09(lng, lat float64) (float64, float64, error) {
	if IsInValidGps(lng, lat) {
		return 0, 0, InvalidGPS
	}
	var z = math.Sqrt(lng*lng+lat*lat) + 0.00002*math.Sin(lat*x_PI)
	var theta = math.Atan2(lat, lng) + 0.000003*math.Cos(lng*x_PI)
	return z*math.Cos(theta) + 0.0065, z*math.Sin(theta) + 0.006, nil
}

func BD09ToGCJ02(lng, lat float64) (float64, float64, error) {
	if IsInValidGps(lng, lat) {
		return 0, 0, InvalidGPS
	}

	var x = lng - 0.0065
	var y = lat - 0.006
	var z = math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*x_PI)
	var theta = math.Atan2(y, x) - 0.000003*math.Cos(x*x_PI)
	return z * math.Cos(theta), z * math.Sin(theta), nil
}

func BD09ToWGS84(lng, lat float64) (float64, float64, error) {
	glng, glat, err := BD09ToGCJ02(lng, lat)
	if err != nil {
		return 0, 0, InvalidGPS
	}
	return GCJ02ToWGS84(glng, glat)
}

func WGS84ToBD09(lng, lat float64) (float64, float64, error) {
	tlng, tlat, err := WGS84ToGCJ02(lng, lat)
	if err != nil {
		return 0, 0, InvalidGPS
	}
	return GCJ02ToBD09(tlng, tlat)
}

func WebMCToWGS84(mercartorX, mercartorY float64) (float64, float64, error) {
	if !(mercartorX >= -earchHalfCir && mercartorX <= earchHalfCir) {
		return 0, 0, InvalidGPS
	}

	lng := mercartorX / earchHalfCir * 180
	lat := mercartorY / earchHalfCir * 180
	lat = 180 / PI * (2*math.Atan(math.Exp(lat*PI/180)) - PI/2)
	return lng, lat, nil
}

func WGS84ToWebMC(lng, lat float64) (float64, float64, error) {
	if IsInValidGps(lng, lat) {
		return 0, 0, InvalidGPS
	}
	x := lng * earchHalfCir / 180
	y := math.Log(math.Tan((90+lat)*PI/360)) / (PI / 180)
	y = y * earchHalfCir / 180
	return x, y, nil
}

// constants for baidu mercartor conversion
var mcband = []float64{12890594.86, 8362377.87, 5591021, 3481989.83, 1678043.12, 0}
var mc2ll = [][]float64{
	[]float64{1.410526172116255e-8, 0.00000898305509648872, -1.9939833816331, 200.9824383106796, -187.2403703815547, 91.6087516669843, -23.38765649603339, 2.57121317296198, -0.03801003308653, 17337981.2},
	[]float64{-7.435856389565537e-9, 0.000008983055097726239, -0.78625201886289, 96.32687599759846, -1.85204757529826, -59.36935905485877, 47.40033549296737, -16.50741931063887, 2.28786674699375, 10260144.86},
	[]float64{-3.030883460898826e-8, 0.00000898305509983578, 0.30071316287616, 59.74293618442277, 7.357984074871, -25.38371002664745, 13.45380521110908, -3.29883767235584, 0.32710905363475, 6856817.37},
	[]float64{-1.981981304930552e-8, 0.000008983055099779535, 0.03278182852591, 40.31678527705744, 0.65659298677277, -4.44255534477492, 0.85341911805263, 0.12923347998204, -0.04625736007561, 4482777.06},
	[]float64{3.09191371068437e-9, 0.000008983055096812155, 0.00006995724062, 23.10934304144901, -0.00023663490511, -0.6321817810242, -0.00663494467273, 0.03430082397953, -0.00466043876332, 2555164.4},
	[]float64{2.890871144776878e-9, 0.000008983055095805407, -3.068298e-8, 7.47137025468032, -0.00000353937994, -0.02145144861037, -0.00001234426596, 0.00010322952773, -0.00000323890364, 826088.5},
}
var llband = []float64{75, 60, 45, 30, 15, 0}
var ll2mc = [][]float64{
	[]float64{-0.0015702102444, 111320.7020616939, 1704480524535203, -10338987376042340, 26112667856603880, -35149669176653700, 26595700718403920, -10725012454188240, 1800819912950474, 82.5},
	[]float64{0.0008277824516172526, 111320.7020463578, 647795574.6671607, -4082003173.641316, 10774905663.51142, -15171875531.51559, 12053065338.62167, -5124939663.577472, 913311935.9512032, 67.5},
	[]float64{0.00337398766765, 111320.7020202162, 4481351.045890365, -23393751.19931662, 79682215.47186455, -115964993.2797253, 97236711.15602145, -43661946.33752821, 8477230.501135234, 52.5},
	[]float64{0.00220636496208, 111320.7020209128, 51751.86112841131, 3796837.749470245, 992013.7397791013, -1221952.21711287, 1340652.697009075, -620943.6990984312, 144416.9293806241, 37.5},
	[]float64{-0.0003441963504368392, 111320.7020576856, 278.2353980772752, 2485758.690035394, 6070.750963243378, 54821.18345352118, 9540.606633304236, -2710.55326746645, 1405.483844121726, 22.5},
	[]float64{-0.0003218135878613132, 111320.7020701615, 0.00369383431289, 823725.6402795718, 0.46104986909093, 2351.343141331292, 1.58060784298199, 8.77738589078284, 0.37238884252424, 7.45},
}

func BDMCToGCJ02(mercartorX, mercartorY float64) (float64, float64, error) {

	mercartorX, mercartorY = math.Abs(mercartorX), math.Abs(mercartorY)
	var f []float64
	for i := 0; i < len(mcband); i++ {
		if mercartorY >= mcband[i] {
			f = mc2ll[i]
			break
		}
	}
	if len(f) == 0 {
		for i := 0; i < len(mcband); i++ {
			if -mercartorY <= -mcband[i] {
				f = mc2ll[i]
				break
			}
		}
	}
	return convert(mercartorX, mercartorY, f)
}

func convert(lng, lat float64, f []float64) (float64, float64, error) {
	if len(f) == 0 {
		return 0, 0, InvalidGPS
	}
	tlng := f[0] + f[1]*math.Abs(lng)
	cc := math.Abs(lat) / f[9]

	var tlat float64
	for i := 0; i <= 6; i++ {
		tlat += (f[i+2] * math.Pow(cc, float64(i)))
	}

	if lng < 0 {
		tlng *= -1
	}
	if lat < 0 {
		tlat *= -1
	}
	return tlng, tlat, nil
}

func GCJ02ToBDMC(lng, lat float64) (float64, float64, error) {
	lng = getLoop(lng, -180, 180)
	lat = getRange(lat, -74, 74)
	var f []float64
	for i := 0; i < len(llband); i++ {
		if lat >= llband[i] {
			f = ll2mc[i]
			break
		}
	}
	if len(f) > 0 {
		for i := len(llband) - 1; i >= 0; i-- {
			if lat <= -llband[i] {
				f = ll2mc[i]
				break
			}
		}
	}
	return convert(lng, lat, f)
}

func getLoop(lng, min, max float64) float64 {
	for lng > max {
		lng -= (max - min)
	}
	for lng < min {
		lng += (max - min)
	}
	return lng
}

func getRange(lat, min, max float64) float64 {
	if min != 0 {
		lat = math.Max(lat, min)
	}
	if max != 0 {
		lat = math.Min(lat, max)
	}
	return lat
}

func BDMCToWGS84(mercartorX, mercartorY float64) (float64, float64, error) {
	gcjLng, gcjLat, _ := BDMCToGCJ02(mercartorX, mercartorY)
	return GCJ02ToWGS84(gcjLng, gcjLat)
}

func BDMCToBD09(mercartorX, mercartorY float64) (float64, float64, error) {
	gcjLng, gcjLat, _ := BDMCToGCJ02(mercartorX, mercartorY)
	return GCJ02ToBD09(gcjLng, gcjLat)
}

func WGS84ToBDMC(lng, lat float64) (float64, float64, error) {
	gcjLng, gcjLat, err := WGS84ToGCJ02(lng, lat)
	if err != nil {
		return 0, 0, err
	}
	return GCJ02ToBDMC(gcjLng, gcjLat)
}

func main() {
	// you should delete main method and copy this file to your project
}
