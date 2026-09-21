package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
	core_http_utils "github.com/DenisAstrakhan/api-server/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetUsers handler")

	limit, offset, err := getlimitOffsetQueryParam(r)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get 'limit'/'offset query param'")
		return
	}
	userDomain, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get users")
		return
	}
	response := GetUsersResponse(usersDTOFromDomains(userDomain))
	responseHendler.JSONResponse(response, http.StatusOK)

}

func getlimitOffsetQueryParam(r *http.Request) (*int, *int, error) {
	limit, err := core_http_utils.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
