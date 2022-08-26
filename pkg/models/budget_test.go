package models_test

import (
	"time"

	"github.com/envelope-zero/backend/internal/database"
	"github.com/envelope-zero/backend/pkg/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func (suite *TestSuiteEnv) TestBudgetCalculations() {
	marchFifteenthTwentyTwentyTwo := time.Date(2022, 3, 15, 0, 0, 0, 0, time.UTC)

	budget := models.Budget{}
	err := database.DB.Save(&budget).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	emptyBudget := models.Budget{}
	err = database.DB.Save(&emptyBudget).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	bankAccount := models.Account{
		AccountCreate: models.AccountCreate{
			BudgetID: budget.ID,
			OnBudget: true,
			External: false,
		},
	}
	err = database.DB.Save(&bankAccount).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	cashAccount := models.Account{
		AccountCreate: models.AccountCreate{
			BudgetID: budget.ID,
			OnBudget: true,
			External: false,
		},
	}
	err = database.DB.Save(&cashAccount).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	employerAccount := models.Account{
		AccountCreate: models.AccountCreate{
			BudgetID: budget.ID,
			External: true,
		},
	}
	err = database.DB.Save(&employerAccount).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	groceryAccount := models.Account{
		AccountCreate: models.AccountCreate{
			BudgetID: budget.ID,
			External: true,
		},
	}
	err = database.DB.Save(&groceryAccount).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	category := models.Category{
		CategoryCreate: models.CategoryCreate{
			BudgetID: budget.ID,
		},
	}
	err = database.DB.Save(&category).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	envelope := models.Envelope{
		EnvelopeCreate: models.EnvelopeCreate{
			CategoryID: category.ID,
		},
	}
	err = database.DB.Save(&envelope).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	allocation1 := models.Allocation{
		AllocationCreate: models.AllocationCreate{
			EnvelopeID: envelope.ID,
			Amount:     decimal.NewFromFloat(17.42),
			Month:      time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	err = database.DB.Save(&allocation1).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	allocation2 := models.Allocation{
		AllocationCreate: models.AllocationCreate{
			EnvelopeID: envelope.ID,
			Amount:     decimal.NewFromFloat(24.58),
			Month:      time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	err = database.DB.Save(&allocation2).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	salaryTransaction := models.Transaction{
		TransactionCreate: models.TransactionCreate{
			Date:                 marchFifteenthTwentyTwentyTwo,
			BudgetID:             budget.ID,
			EnvelopeID:           nil,
			SourceAccountID:      employerAccount.ID,
			DestinationAccountID: bankAccount.ID,
			Reconciled:           true,
			Amount:               decimal.NewFromFloat(2800),
		},
	}
	err = database.DB.Save(&salaryTransaction).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	salaryTransactionApril := models.Transaction{
		TransactionCreate: models.TransactionCreate{
			Date:                 marchFifteenthTwentyTwentyTwo.AddDate(0, 1, 0),
			BudgetID:             budget.ID,
			EnvelopeID:           nil,
			SourceAccountID:      employerAccount.ID,
			DestinationAccountID: bankAccount.ID,
			Reconciled:           true,
			Amount:               decimal.NewFromFloat(2800),
		},
	}
	err = database.DB.Save(&salaryTransactionApril).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	outgoingTransactionBank := models.Transaction{
		TransactionCreate: models.TransactionCreate{
			BudgetID:             budget.ID,
			EnvelopeID:           &envelope.ID,
			SourceAccountID:      bankAccount.ID,
			DestinationAccountID: groceryAccount.ID,
			Amount:               decimal.NewFromFloat(87.45),
		},
	}
	err = database.DB.Save(&outgoingTransactionBank).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	outgoingTransactionCash := models.Transaction{
		TransactionCreate: models.TransactionCreate{
			BudgetID:             budget.ID,
			EnvelopeID:           &envelope.ID,
			SourceAccountID:      cashAccount.ID,
			DestinationAccountID: groceryAccount.ID,
			Amount:               decimal.NewFromFloat(23.17),
		},
	}
	err = database.DB.Save(&outgoingTransactionCash).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	budget = budget.WithCalculations()

	shouldBalance := decimal.NewFromFloat(5489.38)
	assert.True(suite.T(), budget.Balance.Equal(shouldBalance), "Balance for budget is not correct. Should be %s, is %s", shouldBalance, budget.Balance)

	// Verify income for used budget in March
	shouldIncome := decimal.NewFromFloat(2800)
	income, err := budget.Income(marchFifteenthTwentyTwentyTwo)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), income.Equal(shouldIncome), "Income is %s, should be %s", income, shouldIncome)

	// Verify total income for used budget
	shouldIncomeTotal := decimal.NewFromFloat(5600)
	income, err = budget.TotalIncome()
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), income.Equal(shouldIncomeTotal), "Income is %s, should be %s", income, shouldIncomeTotal)

	// Verify income for empty budget in March
	income, err = emptyBudget.Income(marchFifteenthTwentyTwentyTwo)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), income.IsZero(), "Income is %s, should be 0", income)

	// Verify total income for empty budget
	income, err = emptyBudget.TotalIncome()
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), income.IsZero(), "Income is %s, should be 0", income)

	// Verify total budgeted for used budget
	budgeted, err := budget.TotalBudgeted()
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), budgeted.Equal(decimal.NewFromFloat(42)), "Income is %s, should be 42", income)

	// Verify total budgeted for empty budget
	budgeted, err = emptyBudget.TotalBudgeted()
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), budgeted.IsZero(), "Income is %s, should be 0", income)
}

func (suite *TestSuiteEnv) TestMonthIncomeNoTransactions() {
	budget := models.Budget{}
	err := database.DB.Save(&budget).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	income, err := budget.Income(time.Date(2022, 3, 15, 0, 0, 0, 0, time.UTC))
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), income.IsZero(), "Income is %s, should be 0", income)
}

func (suite *TestSuiteEnv) TestTotalIncomeNoTransactions() {
	budget := models.Budget{}
	err := database.DB.Save(&budget).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	income, err := budget.TotalIncome()
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), income.IsZero(), "Income is %s, should be 0", income)
}

func (suite *TestSuiteEnv) TestTotalBudgetedNoTransactions() {
	budget := models.Budget{}
	err := database.DB.Save(&budget).Error
	if err != nil {
		suite.Assert().Fail("Resource could not be saved", err)
	}

	budgeted, err := budget.TotalBudgeted()
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), budgeted.IsZero(), "Income is %s, should be 0", budgeted)
}
