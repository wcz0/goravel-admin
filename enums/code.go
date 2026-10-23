package enums

type CodeEnum int

const (
	Success          CodeEnum = 200
	NotFound         CodeEnum = 404
	Unauthorized     CodeEnum = 401
	PermissionDenied CodeEnum = 403
	ServerError      CodeEnum = 500
	Failed           CodeEnum = 400
)

func (c CodeEnum) Message() string {
	switch c {
	case Success:
		return "Success"
	case NotFound:
		return "NotFound"
	case ServerError:
		return "ServerError"
	case Failed:
		return "Failed"
	default:
		return "Unknown error"
	}
}

type StatusEnum int

const (
	StatusSuccess StatusEnum = 0
	StatusFailed  StatusEnum = 1
)
