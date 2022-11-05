package golinks

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/belljustin/golinks/pkg/golinks"
)

type LambdaHandler struct {
	service golinks.Service
	m       *htmlMarshaller
}

func NewLambdaHandler() LambdaHandler {
	return LambdaHandler{
		service: defaultService(),
		m:       &htmlMarshaller{},
	}
}

func (h LambdaHandler) Handle(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "/":
		return h.home()
	case "/health":
		return h.health()
	case "/links":
		return h.createLink(req)
	default:
		return h.followLink(req)
	}
}

func (h LambdaHandler) health() (*events.APIGatewayProxyResponse, error) {
	healthchecks := h.service.Health()
	content, err := h.m.healthChecks(healthchecks)
	if err != nil {
		return h.apiError(http.StatusInternalServerError)
	}

	if healthchecks.Error() {
		return h.apiResponse(http.StatusInternalServerError, content)
	}
	return h.apiResponse(http.StatusOK, content)
}

func (h LambdaHandler) home() (*events.APIGatewayProxyResponse, error) {
	content, err := h.m.home()
	if err != nil {
		return h.apiError(http.StatusInternalServerError)
	}
	return h.apiResponse(http.StatusOK, content)
}

func (h LambdaHandler) createLink(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	values, err := url.ParseQuery(req.Body)
	if err != nil {
		log.Printf("[INFO] failed to parse link data: %s", err)
		return h.apiError(http.StatusBadRequest)
	}

	link, err := parseLinkValues(values)
	if err != nil {
		return h.apiError(http.StatusBadRequest)
	}

	if err := h.service.SetLink(link.Name, link.URL); err != nil {
		return h.apiError(http.StatusInternalServerError)
	}

	content, err := h.m.setLink(link.Name, link.URL)
	if err != nil {
		return h.apiError(http.StatusInternalServerError)
	}

	return h.apiResponse(http.StatusCreated, content)
}

func (h LambdaHandler) followLink(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	name := strings.TrimLeft(req.Path, "/")
	l, err := h.service.GetLink(name)
	if err != nil {
		return h.apiError(http.StatusInternalServerError)
	}

	if l == nil {
		host := req.Headers["Host"]
		log.Printf("[INFO] name '%s' does not exist. Redirecting to %s", name, "https://"+host)
		return h.apiRedirect("https://" + host)
	}
	return h.apiRedirect(l.String())
}

func (h LambdaHandler) apiRedirect(location string) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Location": location}}
	resp.StatusCode = http.StatusTemporaryRedirect
	return &resp, nil
}

func (h LambdaHandler) apiResponse(status int, body []byte) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "text/html"}}
	resp.StatusCode = status
	resp.Body = string(body)
	return &resp, nil
}

func (h LambdaHandler) apiError(status int) (*events.APIGatewayProxyResponse, error) {
	body := []byte(http.StatusText(status))
	return h.apiResponse(status, body)
}
