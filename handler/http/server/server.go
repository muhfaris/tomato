package server

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/tomatool/tomato/resource"
)

type Resource interface {
	resource.Resource

	SetResponse(method string, requestPath string, responseCode int, responseBody []byte) error
	GetRequestsCount(method, path string) (int, error)
}

type Handler struct {
	r map[string]Resource
}

func New(r map[string]Resource) *Handler {
	return &Handler{r}
}

func (h *Handler) setResponse(resourceName, path string, code int, body *gherkin.DocString) error {
	return h.setResponseWithMethod(resourceName, "GET", path, code, body)
}

func (h *Handler) setResponseWithMethodAndNoBody(resourceName, method string, path string, code int) error {
	return h.setResponseWithMethod(resourceName, method, path, code, nil)
}

func (h *Handler) setResponseWithMethod(resourceName, method string, path string, code int, body *gherkin.DocString) error {
	r, ok := h.r[resourceName]
	if !ok {
		return fmt.Errorf("%s not found to set response", resourceName)
	}

	var bodyBytes []byte
	if body != nil {
		bodyBytes = []byte(body.Content)
	}

	return r.SetResponse(method, path, code, bodyBytes)
}

func (h *Handler) verifyRequestsCount(resourceName, target string, expectedRequestCount int) error {
	r, ok := h.r[resourceName]
	if !ok {
		return fmt.Errorf("%s not found to set response", resourceName)
	}

	tt := strings.Split(target, " ")
	if len(tt) != 2 {
		return errors.New("target format should be following `[METHOD] [PATH]`")
	}

	count, err := r.GetRequestsCount(tt[0], tt[1])
	if err != nil {
		return err
	}
	if count != expectedRequestCount {
		return fmt.Errorf("expecting request count to be %d, got %d", expectedRequestCount, count)
	}
	return nil
}
