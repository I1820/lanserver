package actions

import (
	"net/http"
	"net/http/httptest"
)

func (suite *LSTestSuite) Test_HealthzHandler() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/healthz", nil)
	suite.NoError(err)
	suite.engine.ServeHTTP(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
}
