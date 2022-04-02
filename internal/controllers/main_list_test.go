package controllers_test

import (
	"net/http"
	"testing"

	"github.com/envelope-zero/backend/internal/test"
	"github.com/stretchr/testify/assert"
)

var getOverviewTests = []struct {
	path     string
	expected string
}{
	{"/", `{ "links": { "v1": "http:///v1", "version": "http:///version" }}`},
	{"/v1", `{ "links": { "budgets": "http:///v1/budgets" }}`},
	{"/version", `{"data": { "version": "0.0.0" }}`},
}

func TestGetOverview(t *testing.T) {
	for _, tt := range getOverviewTests {
		recorder := test.Request(t, "GET", tt.path, "")

		test.AssertHTTPStatus(t, http.StatusOK, &recorder)
		assert.JSONEq(t, tt.expected, recorder.Body.String())
	}
}
