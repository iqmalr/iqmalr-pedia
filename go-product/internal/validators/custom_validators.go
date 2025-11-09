package validators

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iqmalr-pedia/go-product/pkg/database"
)

func ValidateExists(fl validator.FieldLevel) bool {
	params := fl.Param()
	if params == "" {
		return false
	}
	tableName := strings.Split(params, ",")[0]

	fieldValue := fl.Field()
	if fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			return true
		}
		fieldValue = fieldValue.Elem()
	}

	db := database.GetDB()
	var count int64
	err := db.Table(tableName).Where("id = ?", fieldValue.Interface()).Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}
