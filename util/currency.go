package util

var (
	USD = "USD"
	EUR = "EUR"
	INR = "INR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	return currency == USD || currency == EUR || currency == INR || currency == CAD
}
