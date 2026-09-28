package core_postgres_pool

import "errors"

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violates foreign key") //ошибка нарушения ограничения внешнего ключа
	ErrUnknown            = errors.New("unknown")              //если не удалось смапить ошибку
)
