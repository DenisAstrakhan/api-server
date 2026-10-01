package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
	core_http_types "github.com/DenisAstrakhan/api-server/internal/core/transport/http/types"
)

type PatchTaskRespons TaskDTOResponse

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" example:"Сделать домашку" swaggertype:"string"`
	Description core_http_types.Nullable[string] `json:"description" example:"Сделать Геометрия" swaggertype:"string"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" example:"true" swaggertype:"bool"`
}

func (p *PatchTaskRequest) Validate() error {
	if p.Title.Set {
		if p.Title.Value == nil {
			return fmt.Errorf("'Title' can't be NULL")
		}
		lenTitle := len([]rune(*p.Title.Value))
		if lenTitle < 1 || lenTitle > 100 {
			return fmt.Errorf("'Title' must be between 1 and 100 symbols")
		}
	}
	if p.Description.Set {
		if p.Description.Value != nil {
			lenDescription := len([]rune(*p.Description.Value))
			if lenDescription < 1 || lenDescription > 1000 {
				return fmt.Errorf("'Description' must be between 1 and 1000 symbols")
			}
		}
	}
	if p.Completed.Set {
		if p.Completed.Value == nil {
			return fmt.Errorf("'Complited' can't be NULL")
		}
	}
	return nil
}

// PatchTask godoc
// @Summary Изменение задачи
// @Description Изменение существующей задачи по id
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано** description игнорируется, значени в БД не меняется
// @Description 2. **Явно передано значение** "description":"Сделать Геометрия" устанавливает новое значение в БД
// @Description 3. **Выставленно null** "description":null значение удаляется из БД (set to null)
// @Description Ограничение title и completed не может быть null
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body PatchTaskRequest true "PatchTask тело запроса"
// @Param id path int true "id изменяемой задачи"
// @Success 200 {object} PatchTaskRespons "Успешно изменённая задача"
// @Failure 400 {object} core_http_response.ErrorResponse "Bed request"
// @Failure 404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke PatchTask handler")

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get taskID path values")
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHendler.ErrorResponse(err, "failed decode and validate request")
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to patch task")
		return
	}
	response := PatchTaskRespons(taskDTOFromDomain(taskDomain))
	responseHendler.JSONResponse(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(request.Title.ToDomain(), request.Description.ToDomain(), request.Completed.ToDomain())
}
