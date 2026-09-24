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
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
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
