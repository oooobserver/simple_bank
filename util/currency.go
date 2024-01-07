package util

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case "EUR", "USD", "RMB":
		return true
	}
	return false
}
