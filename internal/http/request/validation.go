package request

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	transalation_en "github.com/go-playground/validator/v10/translations/en"
	transalation_id "github.com/go-playground/validator/v10/translations/id"
	"github.com/labstack/gommon/log"
	"github.com/yaza-putu/golang-starter-mongo-api/internal/database"
	"github.com/yaza-putu/golang-starter-mongo-api/internal/http/response"
	"github.com/yaza-putu/golang-starter-mongo-api/internal/pkg/i18n"
	"github.com/yaza-putu/golang-starter-mongo-api/internal/pkg/logger"
	filePkg "github.com/yaza-putu/golang-starter-mongo-api/pkg/file"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	optFunc func(*attr)

	attr struct {
		M map[string]string
	}
)

func defaultParam() attr {
	return attr{
		M: map[string]string{},
	}
}

func CustomMessage(m map[string]string) optFunc {
	return func(p *attr) {
		p.M = m
	}
}

// Validation form
func Validation(s any, opts ...optFunc) (response.DataApi, error) {
	o := defaultParam()

	for _, fn := range opts {
		fn(&o)
	}

	uni := ut.New(id.New(), en.New(), id.New())
	trans, found := uni.GetTranslator(i18n.Locale())

	if !found {
		logger.New(errors.New("translator not found"), logger.SetType(logger.FATAL))
	}

	v := validator.New()
	switch i18n.Locale() {
	case "id":
		if err := transalation_id.RegisterDefaultTranslations(v, trans); err != nil {
			log.Fatal(err)
			return response.Api(response.SetMessage(err)), err
		}

		break
	default:
		if err := transalation_en.RegisterDefaultTranslations(v, trans); err != nil {
			logger.New(err, logger.SetType(logger.FATAL))
			return response.Api(response.SetMessage("bad request"), response.SetError(err)), err
		}

		break
	}

	_ = v.RegisterTranslation("unique", trans, func(ut ut.Translator) error {
		return ut.Add("unique", i18n.T(i18n.Localize{Key: "validations.unique"}), true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("filetype", trans, func(ut ut.Translator) error {
		return ut.Add("filetype", i18n.T(i18n.Localize{Key: "validations.filetype"}), true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("filetype", fe.Field(), fe.Param())
		return t
	})

	_ = v.RegisterTranslation("when", trans, func(ut ut.Translator) error {
		return ut.Add("when", i18n.T(i18n.Localize{Key: "validations.when"}), true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("when", fe.Field())
		return t
	})

	// register validation
	err := v.RegisterValidation("unique", unique)

	if err != nil {
		return response.Api(response.SetMessage(err)), err
	}

	err = v.RegisterValidation("filetype", filetype)

	if err != nil {
		return response.Api(response.SetMessage(err)), err
	}

	err = v.RegisterValidation("when", when)

	if err != nil {
		return response.Api(response.SetMessage(err)), err
	}

	for k, ms := range o.M {
		rError := v.RegisterTranslation(k, trans, func(ut ut.Translator) error {
			return ut.Add(k, fmt.Sprintf("{0} %s", ms), true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(k, fe.Field())
			return t
		})

		if rError != nil {
			return response.Api(response.SetError(rError), response.SetMessage(i18n.T(i18n.Localize{Key: "validations.badrequest"}))), rError
		}

	}

	r := v.Struct(s)
	m := map[string]interface{}{}

	if r != nil {
		for _, r := range r.(validator.ValidationErrors) {
			m[strings.ToLower(r.Field())] = []string{r.Translate(trans)}
		}

		return response.Api(response.SetCode(422), response.SetMessage(i18n.T(i18n.Localize{Key: "validations.badrequest"})), response.SetError(m)), r
	}

	return response.Api(response.SetCode(200)), nil
}

func filetype(fl validator.FieldLevel) bool {
	param := fl.Param()
	params := strings.Split(param, " ")

	file := fl.Field().Interface().(multipart.File)
	if reflect.TypeOf(file).Kind() == reflect.Ptr && reflect.TypeOf(file).Elem().Name() == "File" {
		return false
	}

	if len(params) < 1 {
		return false
	}

	return filePkg.DetectContentType(file, params)
}

func unique(fl validator.FieldLevel) bool {
	param := fl.Param()
	params := strings.Split(param, ":")
	switch len(params) {
	case 1:
		return true
	case 2:
		dField := fl.Field().String()
		col := strings.ToLower(params[1])
		count, err := database.Mongo.Collection(params[0]).CountDocuments(context.Background(), bson.M{col: dField})
		if err != nil {
			logger.New(err)
		}
		if count > 0 {
			return false
		} else {
			return true
		}
	case 3:
		ignore := fl.Parent().FieldByName(params[2]).String()
		dField := fl.Field().String()
		col := strings.ToLower(params[1])
		col2 := strings.ToLower(params[2])
		count, err := database.Mongo.Collection(params[0]).CountDocuments(context.Background(), bson.M{col: dField, col2: bson.M{"$not": bson.M{"$gt": ignore}}})
		if err != nil {
			logger.New(err)
		}

		if count > 0 {
			return false
		} else {
			return true
		}
	default:
		return true
	}
}

func when(fl validator.FieldLevel) bool {
	param := fl.Param()

	params := strings.Split(param, ":")
	id := fl.Parent().FieldByName("ID").Interface().(primitive.ObjectID)
	blankId := primitive.ObjectID{}

	if len(params) < 2 {
		return false
	}

	switch params[0] {
	case "update":
		if params[1] == "required" && id != blankId && fl.Field().String() == "" {
			return false
		}
		return true
	case "create":
		if params[1] == "required" && id == blankId && fl.Field().String() == "" {
			return false
		}

		return true
	}

	return true
}
