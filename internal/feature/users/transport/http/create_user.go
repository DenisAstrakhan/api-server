package users_transport_http

import (
	"net/http"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

// DTO для запроса пользователя
type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

// DTO для ответа пользоваьтелю
type CreateUserResponse UserDTOResponse

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke CreateUser handler")

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHendler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}
	userDomain := domainFromDTO(request)
	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to create user")
		return
	}
	response := CreateUserResponse(userDTOFromDomain(userDomain))
	responseHendler.JSONResponse(response, http.StatusCreated)
}
func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
