package web_file_system_repository

import (
	"fmt"
	"os"

	core_errors "github.com/DenisAstrakhan/api-server/internal/core/errors"
)

func (r *WebRepository) GetFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			//полученная ошибка вызвана тем что файл не найден
			return nil, fmt.Errorf("file: %s: %w", filePath, core_errors.ErrNotFound)
		}
		return nil, fmt.Errorf("get file: %s: %w", filePath, err)
	}
	return file, nil
}
