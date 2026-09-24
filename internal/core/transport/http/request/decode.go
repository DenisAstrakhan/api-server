package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/DenisAstrakhan/api-server/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

// создаём валидатор
var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v:%w", err, core_errors.ErrInvalidArgument)
	}
	var err error
	// проверяем имеет ли переданный тип кастомные правила валидации
	v, ok := dest.(validatable) // проверяем подходит ли dest под интерфейс validatable
	if ok {
		//проводим кастомную валидацию
		err = v.Validate()
	} else {
		// проводим обычную валидацию
		err = requestValidator.Struct(dest)
	}
	if err != nil {
		return fmt.Errorf("request validation: %v:%w", err, core_errors.ErrInvalidArgument)
	}
	return nil
}
