package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers godoc
// @Summary Список пользователeй
// @Description Возвращает список пользователя с опцианальной пагинацией
// @Tags users
// @Produce json
// @Param limit query int false "Размер списка пользователей"
// @Param offset query int false "Смещение списка пользователей"
// @Success 200 {object} GetUsersResponse "Успешно вернул список пользователей"
// @Failure 400 {object} core_http_response.ErrorResponse "Bed request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users [get]
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
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
