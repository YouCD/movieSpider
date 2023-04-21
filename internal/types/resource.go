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
		return VideoTypeMovie
	case ResourceTV:
		return VideoTypeMovie
	default:
		return ""
	}

}
