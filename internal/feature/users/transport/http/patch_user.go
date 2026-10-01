package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
	core_http_types "github.com/DenisAstrakhan/api-server/internal/core/transport/http/types"
)

type PatchUserRespons UserDTOResponse

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" example:"Ivanov Ivan" swaggertype:"string"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" example:"+71112223344" swaggertype:"string"`
}

// кастомный валидатор для Nullable структуры
func (p *PatchUserRequest) Validate() error {
	if p.FullName.Set {
		if p.FullName.Value == nil {
			return fmt.Errorf("'FullName' can't be NULL")
		}
		lenFullName := len([]rune(*p.FullName.Value))
		if lenFullName < 3 || lenFullName > 100 {
			return fmt.Errorf("'FullName' must be between 3 and 100 symbols")
		}
	}
	if p.PhoneNumber.Set {
		if p.PhoneNumber.Value != nil {
			lenPhoneNumber := len([]byte(*p.PhoneNumber.Value))
			if lenPhoneNumber < 10 || lenPhoneNumber > 15 {
				return fmt.Errorf("'PhoneNumber' must be between 10 and 15 symbols")
			}
			//проверяем начинается ли PhoneNumber с '+'
			if !strings.HasPrefix(*p.PhoneNumber.Value, "+") {
				return fmt.Errorf("'PhoneNumber' must startswith '+' symbol")
			}
		}
	}
	return nil
}

// PatchUser godoc
// @Summary Изменение пользователя
// @Description Изменение существующего пользователя по id
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано** phone_number игнорируется, значени в БД не меняется
// @Description 2. **Явно передано значение** "phone_number":"+78889997766" устанавливает новое значение в БД
// @Description 3. **Выставленно null** "phone_number":null значение удаляется из БД (set to null)
// @Description Ограничение full_name не может быть null
// @Tags users
// @Accept json
// @Produce json
// @Param request body PatchUserRequest true "PatchUser тело запроса"
// @Param id path int true "id изменяемого пользователя"
// @Success 200 {object} PatchUserRespons "Успешно изменённый пользователей"
// @Failure 400 {object} core_http_response.ErrorResponse "Bed request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke PatchUser handler")

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get userID path values")
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHendler.ErrorResponse(err, "failed decode and validate request")
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to patch user")
		return
	}
	respons := PatchUserRespons(userDTOFromDomain(userDomain))
	responseHendler.JSONResponse(respons, http.StatusOK)

}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(request.FullName.ToDomain(), request.PhoneNumber.ToDomain())
}
