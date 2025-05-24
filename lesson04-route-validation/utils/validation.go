package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func HandleValidationErrors(err error) gin.H {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, e := range validationError {
			switch e.Tag() {
			case "gt":
				errors[e.Field()] = e.Field() + " phải lớn hơn giá trị tối thiểu"
			case "lt":
				errors[e.Field()] = e.Field() + " phải nhỏ hơn giá trị tối thiểu"
			case "gte":
				errors[e.Field()] = e.Field() + " phải lớn hơn hoặc bằng giá trị tối thiểu"
			case "lte":
				errors[e.Field()] = e.Field() + " phải nhỏ hơn hoặc bằng giá trị tối thiểu"
			case "uuid":
				errors[e.Field()] = e.Field() + " phải là UUID hợp lệ"
			case "slug":
				errors[e.Field()] = e.Field() + " chỉ được chứa chữ thường, số, dấu gạch ngang hoặc dấu chấm"
			case "min":
				errors[e.Field()] = fmt.Sprintf("%s phải nhiều hơn %s ký tự", e.Field(), e.Param())
			case "max":
				errors[e.Field()] = fmt.Sprintf("%s phải ít hơn %s ký tự", e.Field(), e.Param())
			case "min_int":
				errors[e.Field()] = fmt.Sprintf("%s phải có giá trị lớn hơn %s", e.Field(), e.Param())
			case "max_int":
				errors[e.Field()] = fmt.Sprintf("%s phải có giá trị bé hơn %s", e.Field(), e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[e.Field()] = fmt.Sprintf("%s phải là một trong các giá trị: %s", e.Field(), allowedValues)
			case "required":
				errors[e.Field()] = e.Field() + " là bắt buộc"
			case "search":
				errors[e.Field()] = e.Field() + " chỉ được chứa chữ thường, in hoa, số và khoảng trắng"
			case "email":
				errors[e.Field()] = e.Field() + " phải đúng định dạng là email"
			case "datetime":
				errors[e.Field()] = e.Field() + " phải theo đúng định dạng YYYY-MM-DD"
			case "file_ext":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ",")
				errors[e.Field()] = fmt.Sprintf("%s chỉ cho phép những file có extension: %s", e.Field(), allowedValues)
			}
		}

		return gin.H{"error": errors}

	}

	return gin.H{"error": "Yêu cầu không hợp lệ" + err.Error()}
}

func RegisterValidators() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}

	var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	})

	var searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})

	v.RegisterValidation("min_int", func(fl validator.FieldLevel) bool {
		minStr := fl.Param()
		minVal, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			return false
		}

		return fl.Field().Int() >= minVal
	})

	v.RegisterValidation("max_int", func(fl validator.FieldLevel) bool {
		maxStr := fl.Param()
		maxVal, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			return false
		}

		return fl.Field().Int() <= maxVal
	})

	v.RegisterValidation("file_ext", func(fl validator.FieldLevel) bool {
		filename := fl.Field().String()

		allowedStr := fl.Param()
		if allowedStr == "" {
			return false
		}

		allowedExt := strings.Fields(allowedStr)
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")

		for _, allowed := range allowedExt {
			if ext == strings.ToLower(allowed) {
				return true
			}
		}

		return false
	})

	return nil
}
