package ierror

import (
	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints/i18n"
)

// IAPIError .
type IAPIError interface {
	Render(locale i18n.LocaleResource) string
	Code() string
	HttpCode() int
}
