package types

type TelemetryData struct {
	Distance float64 `json:"distance"`
	OBUID    int     `json:"obuID"`
	Unix     int64   `json:"unix"`
}

type OBUData struct {
	OBUID int     `json:"obuID"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

type Invoice struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	Fee           float64 `json:"fee"`
}
