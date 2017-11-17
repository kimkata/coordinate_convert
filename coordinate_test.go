package main

import (
	"fmt"
	"math"
	"testing"
)

const (
	PRECISION = 1e-4
)

func floatEqual(f1, f2 float64) bool {
	if math.Abs(f1-f2) <= PRECISION {
		return true
	}
	return false
}

// got test data from http://www.gpsspg.com/maps.htm
// note: the test may fail, because of precison reason and no endurable data source

func TestConvert_WGS84ToGCJ02(t *testing.T) {
	slng, slat := 116.4472222200, 39.9308333300
	expectLng, expectLat := 116.4534127584, 39.9321990290
	tlng, tlat, err := WGS84ToGCJ02(slng, slat)
	fmt.Printf("WGS84ToGCJ02 %f, %f expect %f, %f\n", slng, slat, expectLng, expectLat)
	fmt.Printf("WGS84ToGCJ02 %f, %f got    %f, %f, %v\n", slng, slat, tlng, tlat, err)
	if !(floatEqual(tlng, expectLng) && floatEqual(tlat, expectLat)) {
		t.Fatalf("WGS84ToGCJ02 convert failed!")
	}
}

func TestConvert_GCJ02ToWGS84(t *testing.T) {
	slng, slat := 116.4534127584, 39.9321990290
	expectLng, expectLat := 116.4472222200, 39.9308333300
	tlng, tlat, err := GCJ02ToWGS84(slng, slat)
	fmt.Printf("GCJ02ToWGS84 %f, %f expect %f, %f\n", slng, slat, expectLng, expectLat)
	fmt.Printf("GCJ02ToWGS84 %f, %f got    %f, %f, %v\n", slng, slat, tlng, tlat, err)
	if !(floatEqual(tlng, expectLng) && floatEqual(tlat, expectLat)) {
		t.Fatalf("GCJ02ToWGS84 convert failed!")
	}
}

func TestConvert_GCJ02ToBD09(t *testing.T) {
	slng, slat := 116.4534127584, 39.9321990290
	expectLng, expectLat := 116.4600064714, 39.9378464687
	tlng, tlat, err := GCJ02ToBD09(slng, slat)
	fmt.Printf("GCJ02ToBD09 %f, %f expect %f, %f\n", slng, slat, expectLng, expectLat)
	fmt.Printf("GCJ02ToBD09  %f, %f got    %f, %f, %v\n", slng, slat, tlng, tlat, err)
	if !(floatEqual(tlng, expectLng) && floatEqual(tlat, expectLat)) {
		t.Fatalf("GCJ02ToBD09 convert failed!")
	}
}

func TestConvert_BD09ToGCJ02(t *testing.T) {
	slng, slat := 116.4600064714, 39.9378464687
	expectLng, expectLat := 116.4534127584, 39.9321990290
	tlng, tlat, err := BD09ToGCJ02(slng, slat)
	fmt.Printf("BD09ToGCJ02 %f, %f expect %f, %f\n", slng, slat, expectLng, expectLat)
	fmt.Printf("BD09ToGCJ02  %f, %f got    %f, %f, %v\n", slng, slat, tlng, tlat, err)
	if !(floatEqual(tlng, expectLng) && floatEqual(tlat, expectLat)) {
		t.Fatalf("BD09ToGCJ02 convert failed!")
	}
}

func TestConvert_WGS84ToBD09(t *testing.T) {
	slng, slat := 116.4472222200, 39.9308333300
	expectLng, expectLat := 116.46000757747, 39.937860674084
	tlng, tlat, err := WGS84ToBD09(slng, slat)
	fmt.Printf("WGS84ToBD09 %f, %f expect %f, %f\n", slng, slat, expectLng, expectLat)
	fmt.Printf("WGS84ToBD09 %f, %f got    %f, %f, %v\n", slng, slat, tlng, tlat, err)
	if !(floatEqual(tlng, expectLng) && floatEqual(tlat, expectLat)) {
		t.Fatalf("WGS84ToBD09 convert failed!")
	}
}

func TestConvert_GCJ02ToBDMC(t *testing.T) {
	slng, slat := 108.95344, 34.265657
	expectMcX, expectMcY := 12128773.43, 4040249.00
	tlng, tlat, err := GCJ02ToBDMC(slng, slat)
	fmt.Printf("GCJ02ToBDMC %f, %f expect %f, %f\n", slng, slat, expectMcX, expectMcY)
	fmt.Printf("GCJ02ToBDMC %f, %f got    %f, %f, %v\n", slng, slat, tlng, tlat, err)
	if !(floatEqual(tlng, expectMcX) && floatEqual(tlat, expectMcY)) {
		t.Fatalf("GCJ02ToBDMC convert failed!")
	}
}

func TestConvert_BDMCToGCJ02(t *testing.T) {
	//北京钟楼
	// 12128773.43, 4040249.00 -> wgs84: 108.95344,34.265657
	mercartorX, mercartorY := 12128773.43, 4040249.00
	expectLng, expectLat := 108.95344, 34.265657
	tlng, tlat, err := BDMCToGCJ02(mercartorX, mercartorY)
	fmt.Printf("BDMCToGCJ02 %f, %f expect %f, %f\n", mercartorX, mercartorY, expectLng, expectLat)
	fmt.Printf("BDMCToGCJ02 %f, %f got    %f, %f, %v\n", mercartorX, mercartorY, tlng, tlat, err)
	if !(floatEqual(tlng, expectLng) && floatEqual(tlat, expectLat)) {
		t.Fatalf("BDMCToGCJ02 convert failed!")
	}
}

func TestConvert_BDMCToBD09(t *testing.T) {
	mercartorX, mercartorY := 12964409.440011, 4829229.7034175
	expectLng, expectLat := 116.460006, 39.937846
	tlng, tlat, err := BDMCToBD09(mercartorX, mercartorY)
	fmt.Printf("BDMCToBD09 %f, %f expect %f, %f\n", mercartorX, mercartorY, expectLng, expectLat)
	fmt.Printf("BDMCToBD09 %f, %f got    %f, %f, %v\n", mercartorX, mercartorY, tlng, tlat, err)
	if !(floatEqual(tlng, expectLng) && floatEqual(tlat, expectLat)) {
		t.Fatalf("BDMCToBD09 convert failed!")
	}
}

func TestConvert_BDMCToWGS84(t *testing.T) {
	//mercartorX, mercartorY := 12964409.440011, 4829229.7034175
	mercartorX, mercartorY := 12964409.500548, 4829229.8598471
	expectLng, expectLat := 116.4472222200, 39.9308333300
	tlng, tlat, err := BDMCToWGS84(mercartorX, mercartorY)
	fmt.Printf("BDMCToWGS84 %f, %f expect %f, %f\n", mercartorX, mercartorY, expectLng, expectLat)
	fmt.Printf("BDMCToWGS84 %f, %f got    %f, %f, %v\n", mercartorX, mercartorY, tlng, tlat, err)
	if !(floatEqual(tlng, expectLng) && floatEqual(tlat, expectLat)) {
		t.Fatalf("BDMCToWGS84 convert failed!")
	}
}

func TestConvert_WGS84ToBDMC(t *testing.T) {
	slng, slat := 116.4472222200, 39.9308333300
	expectMcX, expectMcY := 12964410.378148, 4829229.9704967
	tlng, tlat, err := WGS84ToBDMC(slng, slat)
	fmt.Printf("WGS84ToBDMC %f, %f expect %f, %f\n", slng, slat, expectMcX, expectMcY)
	fmt.Printf("WGS84ToBDMC %f, %f got    %f, %f, %v\n", slng, slat, tlng, tlat, err)
	if !(floatEqual(tlng, expectMcX) && floatEqual(tlat, expectMcY)) {
		t.Fatalf("WGS84ToBDMC convert failed!")
	}
}
