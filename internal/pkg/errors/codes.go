package errors

import "net/http"

// AppError 统一业务错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

// WithStatus 返回对应的 HTTP 状态码
func (e *AppError) WithStatus() int {
	switch {
	case e.Code >= 10000 && e.Code < 20000:
		return http.StatusBadRequest
	case e.Code >= 20000 && e.Code < 30000:
		return http.StatusUnauthorized
	case e.Code >= 30000 && e.Code < 40000:
		return http.StatusNotFound
	case e.Code >= 40000 && e.Code < 50000:
		return http.StatusNotFound
	case e.Code >= 70000 && e.Code < 80000:
		return http.StatusBadRequest
	case e.Code >= 50000:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}

// ── 认证模块 10000~19999 ──

var (
	ErrInvalidPhone    = &AppError{Code: 10001, Message: "手机号格式不正确"}
	ErrInvalidCode     = &AppError{Code: 10002, Message: "验证码格式不正确"}
	ErrCodeExpired     = &AppError{Code: 10003, Message: "验证码已过期或不存在"}
	ErrCodeSendFreq    = &AppError{Code: 10004, Message: "验证码发送过于频繁，请60秒后再试"}
	ErrCodeSendFailed  = &AppError{Code: 10005, Message: "验证码发送失败"}
	ErrInvalidRole     = &AppError{Code: 10006, Message: "角色参数无效"}
)

// ── 用户模块 20000~29999 ──

var (
	ErrUnauthorized    = &AppError{Code: 20001, Message: "请先登录"}
	ErrTokenExpired    = &AppError{Code: 20002, Message: "登录已过期，请重新登录"}
	ErrForbidden       = &AppError{Code: 20003, Message: "权限不足"}
	ErrUserNotFound    = &AppError{Code: 20004, Message: "用户不存在"}
	ErrNameRequired    = &AppError{Code: 20005, Message: "姓名不能为空"}
	ErrCertRequired    = &AppError{Code: 20006, Message: "公司名称和职位不能为空"}
)

// ── 简历模块 30000~39999 ──

var (
	ErrResumeNotFound   = &AppError{Code: 30001, Message: "简历不存在"}
	ErrResumeTitleEmpty = &AppError{Code: 30002, Message: "简历标题不能为空"}
	ErrInvalidJSON      = &AppError{Code: 30003, Message: "简历内容格式无效"}
)

// ── 职位模块 40000~49999 ──

var (
	ErrJobNotFound    = &AppError{Code: 40001, Message: "职位不存在"}
	ErrJobTitleEmpty  = &AppError{Code: 40002, Message: "职位标题不能为空"}
	ErrJobNotOwned    = &AppError{Code: 40003, Message: "职位不存在或不属于你"}
	ErrInvalidStatus  = &AppError{Code: 40004, Message: "无效的状态值"}
)

// ── 投递模块 50000~59999 ──

var (
	ErrAppNotFound     = &AppError{Code: 50001, Message: "投递记录不存在"}
	ErrAlreadyApplied  = &AppError{Code: 50002, Message: "你已投递过该职位"}
	ErrAppNotOwned     = &AppError{Code: 50003, Message: "投递记录不存在或不在你的权限范围内"}
)

// ── 面试模块 60000~69999 ──

var (
	ErrInvNotFound      = &AppError{Code: 60001, Message: "面试邀约不存在"}
	ErrInvNotAuthorized = &AppError{Code: 60002, Message: "无权操作该面试邀约"}
	ErrAddressInvalid   = &AppError{Code: 60003, Message: "地址格式无效"}
	ErrAppNotFoundForInv = &AppError{Code: 60004, Message: "投递记录不存在"}
)

// ── 文件模块 70000~79999 ──

var (
	ErrFileTooLarge    = &AppError{Code: 70001, Message: "文件大小不能超过20MB"}
	ErrInvalidFileType = &AppError{Code: 70002, Message: "不支持的文件格式，仅支持 PDF、DOC、DOCX、JPG、PNG"}
	ErrFileRequired    = &AppError{Code: 70003, Message: "请选择文件"}
	ErrUploadFailed    = &AppError{Code: 70004, Message: "文件上传失败"}
)

// ── 通用 90000~99999 ──

var (
	ErrInternal    = &AppError{Code: 90001, Message: "服务器内部错误"}
	ErrInvalidParam = &AppError{Code: 90002, Message: "请求参数无效"}
	ErrNotFound    = &AppError{Code: 90003, Message: "资源不存在"}
)
