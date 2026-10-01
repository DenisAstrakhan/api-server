package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary Список задач
// @Description Возвращает список задачь. Опчианально по user_id пользователя, с пагинацией
// @Tags tasks
// @Produce json
// @Param user_id path int false "id пользователя"
// @Param limit query int false "Размер списка пользователей"
// @Param offset query int false "Смещение списка пользователей"
// @Success 200 {object} GetTasksResponse "Успешно вернул список задач"
// @Failure 400 {object} core_http_response.ErrorResponse "Bed request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /task [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetTasks handler")

	userID, limit, offset, err := getUserIDLimitOffsetQueryParam(r)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get user_id/limit/offset query param")
		return
	}
	tasksDomain, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get tasks")
		return
	}
	response := GetTasksResponse(tasksDTOFromDomains(tasksDomain))
	responseHendler.JSONResponse(response, http.StatusOK)
}
func getUserIDLimitOffsetQueryParam(r *http.Request) (*int, *int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
		userIDQueryParamKey = "user_id"
	)
	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	return userID, limit, offset, nil
}
