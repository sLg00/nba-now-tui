package types

type DateProvider interface {
	GetCurrentDate() (string, error)
	GetCurrentSeason() string
}
