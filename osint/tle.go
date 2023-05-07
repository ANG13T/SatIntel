package osint

type TLE struct {
	CommonName string
	SatelliteCatalogNumber int
	ElsetClassificiation string
	InternationalDesignator string
	ElementSetEpoch float64
	FirstDerivativeMeanMotion float64
	SecondDerivativeMeanMotion string
	BDragTerm string
	ElementSetType int
	ElementNumber int
	ChecksumOne int
	OrbitInclination float64
	RightAscension float64
	Eccentrcity float64
	Perigee float64
	MeanAnamoly float64
	MeanMotion float64
	RevolutionNumber int
	ChecksumTwo int
}