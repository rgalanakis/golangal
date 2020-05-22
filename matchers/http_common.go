package matchers

import (
	"errors"
	"net/http/httptest"
)

func requireRespRec(actual interface{}) (*httptest.ResponseRecorder, error) {
	rr, ok := actual.(*httptest.ResponseRecorder)
	if !ok {
		return nil, errors.New("actual must be a *httptest.ResponseRecorder")
	}
	return rr, nil
}

func mustRespSec(actual interface{}) *httptest.ResponseRecorder {
	return actual.(*httptest.ResponseRecorder)
}

const noNegate = "do not negate this matcher"
