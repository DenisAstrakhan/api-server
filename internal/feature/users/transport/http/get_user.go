package users_transport_http

import (
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

type GetUserRespons UserDTOResponse

// GetUser godoc
// @Summary Вщзвращает пользователя
// @Description Возвращает пользователя по id
// @Tags users
// @Produce json
// @Param id path int true "id возвращаемого пользователя"
// @Success 200 {object} GetUserRespons "Успешно вернул пользователя"
// @Failure 400 {object} core_http_response.ErrorResponse "Bed request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [get]
func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetUser handler")

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get userID path values")
		return
	}
	userDomain, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get user")
		return
	}
	respons := GetUserRespons(userDTOFromDomain(userDomain))
	responseHendler.JSONResponse(respons, http.StatusOK)

}
