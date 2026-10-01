package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

type GetTaskRespons TaskDTOResponse

// GetTask godoc
// @Summary Вщзвращает задачу
// @Description Возвращает задачу по id
// @Tags tasks
// @Produce json
// @Param id path int true "id возвращаемй звдвчи"
// @Success 200 {object} GetTaskRespons "Успешно вернул задачу"
// @Failure 400 {object} core_http_response.ErrorResponse "Bed request"
// @Failure 404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetTask handler")

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get taskID path values")
		return
	}
	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get task")
		return
	}
	respons := GetTaskRespons(taskDTOFromDomain(taskDomain))
	responseHendler.JSONResponse(respons, http.StatusOK)
}
