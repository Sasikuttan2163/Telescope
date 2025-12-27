package transport

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Sasikuttan2163/Telescope/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type httpTripper struct {
	mainTripper http.RoundTripper
	headers     map[string]string
}

func (ht *httpTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range ht.headers {
		req.Header.Set(k, v)
	}
	return ht.mainTripper.RoundTrip(req)
}

func FetchToolsOfStar(star config.StarConfig) {
	httpClient := &http.Client{
		Transport: &httpTripper{
			headers:     star.Transport.HTTP.Headers,
			mainTripper: http.DefaultTransport,
		},
	}

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "telescope-client",
		Version: "v0.0.1",
	}, nil)
	transport := &mcp.StreamableClientTransport{
		Endpoint:   star.Transport.HTTP.BaseURL,
		HTTPClient: httpClient,
	}

	session, err := client.Connect(ctx, transport, nil)

	if err != nil {
		log.Fatalf("Failed to connect to %s\nError: %v", star.ID, err)
	}

	defer session.Close()
	res, err := session.ListTools(ctx, nil)

	if err != nil {
		log.Fatalf("Calltool failed: %v", err)
	}
	fmt.Println("Found ", len(res.Tools), " tools in ", star.Name)
}

func FetchTools(mc config.MainConfig) {
	stars := mc.Stars
	for _, star := range stars {
		FetchToolsOfStar(star)
	}
}
