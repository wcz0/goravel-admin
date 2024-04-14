package enums

type CodeEnum int

const (
	Success     CodeEnum = 200
	NotFound    CodeEnum = 404
	ServerError CodeEnum = 500
	Fail        CodeEnum = 400
)

func (c CodeEnum) Message() string {
	switch c {
	case Success:
		return "Success"
	case NotFound:
		return "NotFound"
	case ServerError:
		return "ServerError"
	case Fail:
		return "Fail"
	default:
		return "Unknown error"
	}
}
