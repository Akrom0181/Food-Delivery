package etc

import (
	"math"
)

// Yer radiusi (kilometr)
const EarthRadius = 6371.0

// CalculateDistance - Ikkita koordinata orasidagi masofani hisoblash (km)
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Radyanga o'tkazamiz
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Delta (farq) qiymatlarni hisoblash
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Masofani hisoblash
	distance := EarthRadius * c

	return distance
}
