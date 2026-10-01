package tasks_transport_http

import (
	"net/http"

	"github.com/DenisAstrakhan/api-server/internal/core/domain"
	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

// DTO для запроса пользователя
type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100" example:"Сделать домашку"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000" example:"Сделать Геометрию"`
	AuthorUserID int     `json:"author_user_id" validate:"required" example:"1"`
}

// DTO для ответа пользоваьтелю
type CreateTaskResponse TaskDTOResponse

// CreateTask godoc
// @Summary Создать задачу
// @Description Создать новую задачу в системе с привязкой к пользователю
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "CreateTask тело запроса"
// @Success 201 {object}  CreateTaskResponse "Успешно созданная задача"
// @Failure 404 {object} core_http_response.ErrorResponse "Author  not found"
// @Failure 400 {object} core_http_response.ErrorResponse "Bed request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke CreateUser handler")

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHendler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}
	taskDomain := domainFromDTO(request)
	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to create task")
		return
	}
	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))
	responseHendler.JSONResponse(response, http.StatusCreated)

}

func domainFromDTO(dto CreateTaskRequest) domain.Task {
	return domain.NewTaskUninitialized(dto.Title, dto.Description, dto.AuthorUserID)
}
