package validate

import (
	"net/http"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/labstack/echo/v4"
)

func NewReqValidator() *ReqValidator {
	validate := validator.New()
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")

	zhTranslations.RegisterDefaultTranslations(validate, trans)

	return &ReqValidator{
		validator:  validate,
		translator: trans,
	}
}

type ReqValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func (rv *ReqValidator) Validate(i interface{}) error {
	if err := rv.validator.Struct(i); err != nil {
		errs := err.(validator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errs.Translate(rv.translator))
	}

	return nil
}
