package types

type Resource int

const (
	ResourceOther Resource = iota
	ResourceMovie
	ResourceTV
)

func (r Resource) Typ() string {
	switch r {
	case ResourceMovie:
		return "movie"
	case ResourceTV:
		return "tv"
	default:
		return ""
	}

}
