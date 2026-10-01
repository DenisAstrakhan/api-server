package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

// DeleteTask godoc
// @Summary Удоляет задачу
// @Description Удоляет задачу по id
// @Tags tasks
// @Param id path int true "id удоляемой задачи"
// @Success 204 "Успешное удоление задачи"
// @Failure 400 {object} core_http_response.ErrorResponse "Bed request"
// @Failure 404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke DeleteTask handler")

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get taskID path values")
		return
	}
	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHendler.ErrorResponse(err, "failed to delete task")
		return
	}
	responseHendler.NoContentResponse()
}
