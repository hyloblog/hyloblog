package trafficsize

import "fmt"

type Size uint64

const (
	One      = 1
	Thousand = 1e3
	Million  = 1e6
)

func (size Size) Abbrev(decimals int) string {
	switch {
	case size >= Million:
		return fmt.Sprintf("%.*fM", decimals, float64(size)/Million)
	case size >= Thousand:
		return fmt.Sprintf("%.*fk", decimals, float64(size)/Thousand)
	default:
		return fmt.Sprintf("%.*f", decimals, float64(size))
	}
}

func (size Size) String() string {
	return size.Abbrev(2)
}
