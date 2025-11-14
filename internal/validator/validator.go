package validator

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"
)

var (
	once     sync.Once
	validate *validator.Validate
	trans    ut.Translator
)

func InitValidatorTrans(locale string) {
	once.Do(func() { validatorTrans(locale) })
}

func validatorTrans(locale string) {
	var ok bool
	if validate, ok = binding.Validator.Engine().(*validator.Validate); !ok {
		panic("Failed to initialize the validator")
	}
	// 注册自定义验证消息
	registerValidation()
	// 注册一个获取json tag的自定义方法
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			label = field.Tag.Get("json")
			if label == "" {
				label = field.Tag.Get("form")
			}
		}

		if label == "-" {
			return ""
		}
		if label == "" {
			return field.Name
		}
		return label
	})

	zhT := zh.New() // 中文翻译器
	enT := en.New() // 英文翻译器
	uni := ut.New(enT, zhT, enT)

	trans, ok = uni.GetTranslator(locale)

	if !ok {
		panic("Initialize a language not supported by the validator")
	}
	var err error
	// 注册翻译器
	switch locale {
	case "en":
		err = enTranslations.RegisterDefaultTranslations(validate, trans)
	case "zh":
		err = zhTranslations.RegisterDefaultTranslations(validate, trans)
	default:
		err = enTranslations.RegisterDefaultTranslations(validate, trans)
	}
	if err != nil {
		panic("Failed to register translator when initializing validator")
	}
	// 注册自定义语言翻译
	err = customRegisTranslation()
	if err != nil {
		panic("Failed to register translator when initializing validator")
	}
}

func registerValidation() {
	reg := `^1([38][0-9]|4[579]|5[0-3,5-9]|6[6]|7[0135678]|9[89])\d{8}$`
	// 注册手机号验证规则
	err := validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(reg).MatchString(fl.Field().String())
	})
	if err != nil {
		panic("Failed to register the mobile phone number verification rule")
	}
}

type translation struct {
	tag             string
	translation     string
	override        bool
	customRegisFunc validator.RegisterTranslationsFunc
	customTransFunc validator.TranslationFunc
}

func customRegisTranslation() error {
	translations := []translation{
		{
			tag:         "mobile",
			translation: "{0}格式不正确",
			override:    false,
		},
	}

	return registerTranslation(translations)
}

func registerTranslation(translations []translation) (err error) {
	for _, t := range translations {
		if t.customTransFunc != nil && t.customRegisFunc != nil {
			err = validate.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)
		} else if t.customTransFunc != nil && t.customRegisFunc == nil {
			err = validate.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)
		} else if t.customTransFunc == nil && t.customRegisFunc != nil {
			err = validate.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFunc)
		} else {
			err = validate.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		}

		if err != nil {
			return
		}
	}
	return
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}
		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		zap.S().Warn("警告: 翻译字段错误: %#v", zap.Any("Error reason:", fe))
		return fe.(error).Error()
	}

	return t
}

func TranslateError(err error) string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}
	var msgs []string
	for _, e := range errs {
		msgs = append(msgs, e.Translate(trans))
	}
	return strings.Join(msgs, ",")
}
