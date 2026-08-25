package core_http_request

import (
	"fmt"
	"time"

	"github.com/kupr666/Orange_Team/internal/core/domain"
	core_http_types "github.com/kupr666/Orange_Team/internal/core/transport/http/types"
)

func ToDomainNullableTime(nullable core_http_types.Nullable[string]) (domain.Nullable[time.Time], error) {
	if !nullable.Set {
		return domain.Nullable[time.Time]{Set: false}, nil
	}
	if nullable.Value == nil {
		return domain.Nullable[time.Time]{Value: nil, Set: true}, nil
	}
	t, err := time.Parse(time.RFC3339, *nullable.Value)
	if err != nil {
		return domain.Nullable[time.Time]{}, fmt.Errorf("invalid RFC3339 time: %w", err)
	}
	return domain.Nullable[time.Time]{Value: &t, Set: true}, nil
}
