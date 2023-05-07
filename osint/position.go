package osint

type SatelliteInfo struct {
    Satname           string `json:"satname"`
    Satid             int    `json:"satid"`
    Transactionscount int    `json:"transactionscount"`
}

type Position struct {
    Satlatitude  float64 `json:"satlatitude"`
    Satlongitude float64 `json:"satlongitude"`
    Sataltitude  float64 `json:"sataltitude"`
    Azimuth      float64 `json:"azimuth"`
    Elevation    float64 `json:"elevation"`
    Ra           float64 `json:"ra"`
    Dec          float64 `json:"dec"`
    Timestamp    int64   `json:"timestamp"`
}

type Response struct {
    SatelliteInfo      SatelliteInfo        `json:"info"`
    Positions []Position `json:"positions"`
}