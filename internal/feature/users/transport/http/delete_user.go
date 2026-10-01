package users_transport_http

import (
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

// DeleteUser godoc
// @Summary Удоляет пользователя
// @Description Удоляет пользователя по id
// @Tags users
// @Param id path int true "id удоляемого пользователя"
// @Success 204 "Успешное удоление пользователя"
// @Failure 400 {object} core_http_response.ErrorResponse "Bed request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke DeleteUser handler")
	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get userID path values")
		return
	}
	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHendler.ErrorResponse(err, "failed to delete user")
		return
	}
	responseHendler.NoContentResponse()
}
