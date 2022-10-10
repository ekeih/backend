package controllers_test

import (
	"net/http"

	"github.com/envelope-zero/backend/pkg/models"
	"github.com/envelope-zero/backend/test"
	"github.com/shopspring/decimal"
)

func (suite *TestSuiteStandard) TestCleanup() {
	_ = suite.createTestBudget(models.BudgetCreate{})
	_ = suite.createTestAccount(models.AccountCreate{})
	_ = suite.createTestCategory(models.CategoryCreate{})
	_ = suite.createTestEnvelope(models.EnvelopeCreate{})
	_ = suite.createTestAllocation(models.AllocationCreate{})
	_ = suite.createTestTransaction(models.TransactionCreate{Amount: decimal.NewFromFloat(17.32)})

	tests := []string{
		"http://example.com/v1/budgets",
		"http://example.com/v1/accounts",
		"http://example.com/v1/categories",
		"http://example.com/v1/transactions",
		"http://example.com/v1/envelopes",
		"http://example.com/v1/allocations",
	}

	// Delete
	recorder := test.Request(suite.controller, suite.T(), http.MethodDelete, "http://example.com/v1", "")
	suite.assertHTTPStatus(&recorder, http.StatusNoContent)

	// Verify
	for _, tt := range tests {
		suite.Run(tt, func() {
			recorder := test.Request(suite.controller, suite.T(), http.MethodGet, tt, "")
			suite.assertHTTPStatus(&recorder, http.StatusOK)

			var response struct {
				Data []any `json:"data"`
			}

			suite.decodeResponse(&recorder, &response)
			suite.Assert().Len(response.Data, 0, "There are resources left for type %s", tt)
		})
	}
}
