package httputil_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/envelope-zero/backend/v2/pkg/controllers"
	"github.com/envelope-zero/backend/v2/pkg/httputil"
	"github.com/envelope-zero/backend/v2/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetURLFields(t *testing.T) {
	url, _ := url.Parse("http://example.com/api/v1/accounts?budget=87645467-ad8a-4e16-ae7f-9d879b45f569&onBudget=false&name=")

	queryFields, setFields := httputil.GetURLFields(url, controllers.AccountQueryFilter{})

	assert.Equal(t, []interface{}{"BudgetID", "OnBudget"}, queryFields)
	assert.Equal(t, []string{"Name", "BudgetID", "OnBudget"}, setFields)
}

func TestGetBodyFields(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.PATCH("/", func(ctx *gin.Context) {
		fields, err := httputil.GetBodyFields(c, models.AccountCreate{})
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, fields)
	})

	json := []byte(`{ "name": "test account" }`)

	c.Request, _ = http.NewRequest(http.MethodPatch, "https://example.com/", bytes.NewBuffer(json))
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code, "Status is wrong, return body %#v", w.Body.String())
}

func TestGetBodyFieldsNull(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.PATCH("/", func(ctx *gin.Context) {
		fields, err := httputil.GetBodyFields(c, models.AccountCreate{})
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, fields)
	})

	json := []byte(`{ "name": null }`)

	c.Request, _ = http.NewRequest(http.MethodPatch, "https://example.com/", bytes.NewBuffer(json))
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusOK, w.Code, "Status is wrong, return body %#v", w.Body.String())
	assert.Equal(t, `["Name"]`, w.Body.String(), `Fields are not parsed correctly, should be ["Name"]`)
}

func TestGetBodyFieldsUnparseable(t *testing.T) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.PATCH("/", func(ctx *gin.Context) {
		fields, err := httputil.GetBodyFields(c, models.AccountCreate{})
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		c.JSON(http.StatusOK, fields)
	})

	json := []byte(`{ "name": "test account }`)

	c.Request, _ = http.NewRequest(http.MethodPatch, "https://example.com/", bytes.NewBuffer(json))
	r.ServeHTTP(w, c.Request)
	assert.Equal(t, http.StatusBadRequest, w.Code, "Status is wrong, return body %#v", w.Body.String())
}
