package models

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}


func (e Error) ErrorCode() int {
	return e.Code
}