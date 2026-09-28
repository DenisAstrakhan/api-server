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
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
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
