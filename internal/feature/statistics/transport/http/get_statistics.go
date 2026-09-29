package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	core_logger "github.com/DenisAstrakhan/api-server/internal/core/logger"
	core_http_request "github.com/DenisAstrakhan/api-server/internal/core/transport/http/request"
	core_http_response "github.com/DenisAstrakhan/api-server/internal/core/transport/http/response"
)

type GetStatisticsResponse StatisticsDTOResponse

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke GetStatistics handler")

	userID, from, to, err := getUserIDFromToQueryParam(r)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get user_id/from/to query param")
		return
	}

	statisticsDomain, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHendler.ErrorResponse(err, "failed to get statistics")
		return
	}
	response := GetStatisticsResponse(statisticsDTOFromDomain(statisticsDomain))
	responseHendler.JSONResponse(response, http.StatusOK)

}

func getUserIDFromToQueryParam(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
		userIDQueryParamKey = "user_id"
	)
	from, err := core_http_request.GetDataQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDataQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	return userID, from, to, nil
}
