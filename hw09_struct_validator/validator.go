package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("Field: %s, Error: %v", err.Field, err.Err))
	}
	return sb.String()
}

type (
	resolveValidator  []string
	ErrValidateString string
)

var (
	resolveValidatorStr = resolveValidator{"len", "regexp", "in"}
	resolveValidatorInt = resolveValidator{"min", "max", "in"}
	ErrWrongArgument    = errors.New("get wrong argument")
	ErrWrongFormat      = errors.New("wrong format")
	ErrWrongParamType   = errors.New("wrong param for type")
	ErrValidateMin      = errors.New("min validate err")
	ErrValidateMax      = errors.New("max validate err")
	ErrValidateLen      = errors.New("len validate err")
	ErrValidateIn       = errors.New("in validate err")
	ErrValidateRegexp   = errors.New("regexp validate err")
)

func Validate(v interface{}) error {
	var resultErr ValidationErrors
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Kind() != reflect.Struct {
		return ErrWrongArgument
	}
	for num := 0; num < rt.NumField(); num++ {
		fieldT := rt.Field(num)
		fieldV := rv.Field(num)

		if fieldT.PkgPath == "" { // определяем является ли поле публичным
			fldTag := fieldT.Tag.Get("validate")
			if fldTag == "" {
				// пропускаем поля без тегов
				continue
			}
			// создаем мапу для валидации параметра
			paramMap, err := generateValidateMap(fldTag)
			if err != nil {
				resultErr = append(resultErr, ValidationError{fieldT.Name, err})
				continue
			}

			switch {
			case fieldT.Type.Kind() == reflect.String:
				err := validateField(fieldV, fieldT, paramMap, resolveValidatorStr)
				resultErr = append(resultErr, err...)
			case fieldT.Type.Kind() == reflect.Int:
				err := validateField(fieldV, fieldT, paramMap, resolveValidatorInt)
				resultErr = append(resultErr, err...)
			case fieldT.Type.Kind() == reflect.Slice:
				sliceValue := fieldV

				// Проходимся по каждому элементу слайса
				for i := 0; i < sliceValue.Len(); i++ {
					// Преобразуем элемент в reflect.Value
					elemValue := sliceValue.Index(i)
					elemType := elemValue.Type()

					resolve := resolveValidatorStr
					if elemType.Kind() == reflect.Int {
						resolve = resolveValidatorInt
					}
					// Вызываем validateField для каждого элемента
					err := validateField(elemValue, fieldT, paramMap, resolve)
					resultErr = append(resultErr, err...)
				}
			case fieldT.Type.Kind() == reflect.Struct:
				structErr := Validate(fieldV.Interface())
				if structErr != nil {
					// тут можно сделать преобразование чтобы каждый элемент добавлялся к текущему слайсу по отдельности с новым именем
					resultErr = append(resultErr, ValidationError{fieldT.Name, structErr})
				}
			}
		}
	}

	if len(resultErr) == 0 {
		return nil
	}
	return resultErr
}

func generateValidateMap(s string) (map[string]string, error) {
	parts := strings.Split(s, "|")
	resultMap := make(map[string]string, len(parts))
	for _, part := range parts {
		getPair := strings.Split(part, ":")
		if len(getPair) != 2 || getPair[0] == "" || getPair[1] == "" {
			return nil, ErrWrongFormat
		}
		resultMap[getPair[0]] = getPair[1]
	}
	return resultMap, nil
}

func validateField(
	v reflect.Value, t reflect.StructField, paramMap map[string]string, r resolveValidator,
) ValidationErrors {
	var resultErr ValidationErrors
	for paramName, paramValue := range paramMap {
		if !slices.Contains(r, paramName) {
			resultErr = append(resultErr, ValidationError{v.Type().Name(), ErrWrongParamType})
			continue
		}
		var err error
		// тут можно заменить мапой, типа map[string]func() но так мне кажется нагляднее
		switch paramName {
		case "len":
			err = lenValidator(v.String(), paramValue)
		case "regexp":
			err = regexpValidator(v.String(), paramValue)
		case "in":
			var inV string
			if t.Type.Kind() == reflect.Int {
				inV = strconv.FormatInt(v.Int(), 10)
			} else {
				inV = v.String()
			}
			err = inValidator(inV, paramValue)
		case "min":
			err = minValidator(v.Int(), paramValue)
		case "max":
			err = maxValidator(v.Int(), paramValue)
		}
		if err != nil {
			resultErr = append(resultErr, ValidationError{t.Name, err})
		}
	}
	return resultErr
}

func lenValidator(s string, l string) error {
	intLen, err := strconv.Atoi(l)
	if err != nil {
		return err
	}
	if len(s) > intLen {
		return ErrValidateLen
	}
	return nil
}

func regexpValidator(s string, val string) error {
	re, err := regexp.Compile(val)
	if err != nil {
		return err
	}
	res := re.MatchString(s)
	if !res {
		return ErrValidateRegexp
	}
	return nil
}

func minValidator(i int64, val string) error {
	intMin, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}
	if i < intMin {
		return ErrValidateMin
	}
	return nil
}

func maxValidator(i int64, val string) error {
	intMax, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}
	if i > intMax {
		return ErrValidateMax
	}
	return nil
}

func inValidator(i string, val string) error {
	valSlice := strings.Split(val, ",")

	if !slices.Contains(valSlice, i) {
		return ErrValidateIn
	}

	return nil
}
