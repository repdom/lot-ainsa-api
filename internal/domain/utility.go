package domain

import "math"

func RoundToTwoDecimals(value float64) float64 {
	var round float64
	if ((value * 100) - math.Floor(value*100)) >= 0.5 {
		round = math.Ceil(value*100) / 100
	} else {
		round = math.Floor(value*100) / 100
	}
	//	log.Printf("roundToTwoDecimals(%f) \n", value)
	//	log.Println("round = ", round)
	return round
}
