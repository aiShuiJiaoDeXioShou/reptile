package res

import "time"

type ResponseType struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
	Time    time.Time   `json:"time"`
}

func Ok(data interface{}) ResponseType {
	return ResponseType{
		Code:    200,
		Message: "success",
		Data:    data,
		Time:    time.Now(),
	}
}

// Fail 失败
func Fail(msg string) ResponseType {
	return ResponseType{
		Code:    -1,
		Message: msg,
		Data:    nil,
		Time:    time.Now(),
	}
}

// Error 错误
func Error(msg string) ResponseType {
	return ResponseType{
		Code:    500,
		Message: msg,
		Data:    nil,
		Time:    time.Now(),
	}
}

// 未登入
func Unlogin(msg string) ResponseType {
	return ResponseType{
		Code:    401,
		Message: msg,
		Data:    nil,
		Time:    time.Now(),
	}
}

// 未授权
func Unauthorized(msg string) ResponseType {
	return ResponseType{
		Code:    403,
		Message: msg,
		Data:    nil,
		Time:    time.Now(),
	}
}

// 格式错误
func FormatError(msg string) ResponseType {
	return ResponseType{
		Code:    400,
		Message: msg,
		Data:    nil,
		Time:    time.Now(),
	}
}
