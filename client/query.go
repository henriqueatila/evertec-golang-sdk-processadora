package client

import (
	"net/url"
	"strconv"
)

// pathParam escapes a path parameter for safe URL interpolation.
// This prevents path traversal and query injection from crafted IDs.
func pathParam(s string) string { return url.PathEscape(s) }

// PaginationParams defines common pagination parameters used across list endpoints.
// These fields are shared by ListAccountsParams, ListCardsParams, ListTransactionsRequest, etc.
type PaginationParams struct {
	Limit         int
	StartingAfter string
	EndingBefore  string
}

// addPaginationParams adds common pagination query parameters to url.Values.
// This helper reduces duplication across query builder functions.
func addPaginationParams(q url.Values, limit int, startingAfter, endingBefore string) {
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if startingAfter != "" {
		q.Set("starting_after", startingAfter)
	}
	if endingBefore != "" {
		q.Set("ending_before", endingBefore)
	}
}
