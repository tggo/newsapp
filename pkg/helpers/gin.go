package helpers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	legacyRouter "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// GetRouters Helper function to create a router during testing
func GetRouters(swaggerFileName string) (*gin.Engine, routers.Router, error) {
	gin.SetMode(gin.TestMode)
	ginRouter := gin.Default()

	ctx := context.Background()

	loader := openapi3.Loader{Context: ctx}
	doc, err := loader.LoadFromFile(swaggerFileName)
	if err != nil {
		return nil, nil, err
	}

	err = doc.Validate(ctx)
	if err != nil {
		return nil, nil, err
	}

	swaggerRouter, errNewRouter := legacyRouter.NewRouter(doc)
	if errNewRouter != nil {
		return nil, nil, errNewRouter
	}

	return ginRouter, swaggerRouter, nil
}

type checkFunc func(t *testing.T, w *httptest.ResponseRecorder, data TestData) bool

// TestHTTPResponse Helper function to process a request and test its response
func TestHTTPResponse(t *testing.T, r http.Handler, req *http.Request, data TestData, f checkFunc) {
	// CreatePriceAlert a response recorder
	w := httptest.NewRecorder()

	// CreatePriceAlert the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(t, w, data) {
		t.Fail()
	}
}

type TestData struct {
	Expected interface{}
	Status   int
	Response interface{}
}

func Check(t *testing.T, w *httptest.ResponseRecorder, data TestData) bool {
	statusOK := w.Code == data.Status
	p, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	err = json.Unmarshal(p, data.Response)
	require.NoError(t, err)
	jsonString, err1 := json.Marshal(data.Expected)
	require.NoError(t, err1)
	err1 = json.Unmarshal(jsonString, data.Expected)
	require.NoError(t, err1)
	require.Equal(t, data.Expected, data.Response)
	return statusOK
}
