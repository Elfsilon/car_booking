package appraiser

type RentAppraiser interface {
	Appraise(days int) float64
}
