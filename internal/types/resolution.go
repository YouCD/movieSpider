package types

type Resolution int

const (
	ResolutionOther Resolution = iota
	Resolution480P
	Resolution720P
	Resolution1080P
	Resolution2160P
)

func (r Resolution) Res() string {
	//nolint:exhaustive
	switch r {
	case Resolution480P:
		return "480"
	case Resolution720P:
		return "720"
	case Resolution1080P:
		return "1080"
	case Resolution2160P:
		return "2160"
	default:
		return ""
	}
}

func (r Resolution) ResolutionStr2Int(res string) Resolution {
	switch res {
	case "480":
		return Resolution480P
	case "720":
		return Resolution720P
	case "1080":
		return Resolution1080P
	case "2160":
		return Resolution2160P
	default:
		return Resolution1080P
	}
}
