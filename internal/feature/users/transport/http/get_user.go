package users_transport_http

import (
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

type GetUserRespons UserDTOResponse

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
