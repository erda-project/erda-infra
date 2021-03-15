package errorresp

import (
	"fmt"
	"net/http"

	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints/i18n"
)

// Error Codes i18n templates
var (
	templateMissingParameter      = i18n.NewTemplate("MissingParameter", "缺少参数 %s")
	templateInvalidParameter      = i18n.NewTemplate("InvalidParameter", "参数错误 %s")
	templateInvalidState          = i18n.NewTemplate("InvalidState", "状态异常 %s")
	templateNotLogin              = i18n.NewTemplate("NotLogin", "未登录")
	templateAccessDenied          = i18n.NewTemplate("AccessDenied", "无权限")
	templateNotFound              = i18n.NewTemplate("NotFound", "资源不存在")
	templateAlreadyExists         = i18n.NewTemplate("AlreadyExists", "资源已存在")
	templateInternalError         = i18n.NewTemplate("InternalError", "异常 %s")
	templateErrorVerificationCode = i18n.NewTemplate("ErrorVerificationCode", "验证码错误 %s")
)

// MissingParameter .
func (e *APIError) MissingParameter(err string) *APIError {
	return e.dup().appendCode(http.StatusBadRequest, "MissingParameter").
		appendLocaleTemplate(templateMissingParameter, err)
}

// InvalidParameter .
func (e *APIError) InvalidParameter(err interface{}) *APIError {
	return e.dup().appendCode(http.StatusBadRequest, "InvalidParameter").
		appendLocaleTemplate(templateInvalidParameter, toString(err))
}

// InvalidState .
func (e *APIError) InvalidState(err string) *APIError {
	return e.dup().appendCode(http.StatusBadRequest, "InvalidState").
		appendLocaleTemplate(templateInvalidState, err)
}

// NotLogin .
func (e *APIError) NotLogin() *APIError {
	return e.dup().appendCode(http.StatusUnauthorized, "NotLogin").
		appendLocaleTemplate(templateNotLogin)
}

// AccessDenied .
func (e *APIError) AccessDenied() *APIError {
	return e.dup().appendCode(http.StatusForbidden, "AccessDenied").
		appendLocaleTemplate(templateAccessDenied)
}

// NotFound .
func (e *APIError) NotFound() *APIError {
	return e.dup().appendCode(http.StatusNotFound, "NotFound").
		appendLocaleTemplate(templateNotFound)
}

// IsNotFound .
func IsNotFound(e error) bool {
	return getCode(e) == "NotFound"
}

// AlreadyExists .
func (e *APIError) AlreadyExists() *APIError {
	return e.dup().appendCode(http.StatusConflict, "AlreadyExists").
		appendLocaleTemplate(templateAlreadyExists)
}

// InternalError .
func (e *APIError) InternalError(err error) *APIError {
	return e.dup().appendCode(http.StatusInternalServerError, "InternalError").
		appendLocaleTemplate(templateInternalError, err.Error())
}

// ErrorVerificationCode .
func (e *APIError) ErrorVerificationCode(err error) *APIError {
	return e.dup().appendCode(http.StatusInternalServerError, "ErrorVerificationCode").
		appendLocaleTemplate(templateErrorVerificationCode, err.Error())
}

func toString(err interface{}) string {
	switch t := err.(type) {
	case string:
		return err.(string)
	case error:
		return err.(error).Error()
	default:
		_ = t
		return fmt.Sprintf("%v", err)
	}
}

func getCode(e error) string {
	switch t := e.(type) {
	case *APIError:
		return t.code
	default:
		return ""
	}
}
