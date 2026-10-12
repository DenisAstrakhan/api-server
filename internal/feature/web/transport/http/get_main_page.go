package web_transport_http

import (
	"net/http"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetMainPage(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetMainPage handler")

	html, err := h.webService.GetMainPage()
	if err != nil {
		responseHendler.ErrorResponse(err, "failet to get index.html for main page")
		return
	}
	responseHendler.HTMLResponse(html)
}
