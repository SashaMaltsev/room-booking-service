package httptransport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	scheduledomain "github.com/SashaMaltsev/room-booking-service/internal/domain/schedule"
)

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func parseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func parsePagination(r *http.Request) (commondomain.PageRequest, error) {
	query := r.URL.Query()

	page := 0
	if raw := strings.TrimSpace(query.Get("page")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return commondomain.PageRequest{}, err
		}
		page = value
	}

	pageSize := 0
	if raw := strings.TrimSpace(query.Get("pageSize")); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil {
			return commondomain.PageRequest{}, err
		}
		pageSize = value
	}

	return commondomain.PageRequest{
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func parseWeekdays(values []int) ([]scheduledomain.Weekday, error) {
	days := make([]scheduledomain.Weekday, 0, len(values))
	for _, value := range values {
		day, err := scheduledomain.ParseWeekday(value)
		if err != nil {
			return nil, err
		}

		days = append(days, day)
	}

	return days, nil
}

func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "

	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}

	return token, true
}
