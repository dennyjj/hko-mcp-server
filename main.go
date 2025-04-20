package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const baseUrl = "https://data.weather.gov.hk/weatherAPI/opendata/weather.php"
const (
	currentWeatherReport = "rhrread"
)
const (
	langTc = "tc"
	langEn = "en"
)

func main() {
	s := server.NewMCPServer(
		"HKO MCP Server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	currentWeatherTool := mcp.NewTool("current-weather",
		mcp.WithDescription("Fetch and display the current weather"),
	)

	s.AddTool(currentWeatherTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		apiUrl := baseUrl + "?dataType=" + currentWeatherReport + "&lang=" + langTc

		var req *http.Request
		var err error
		req, err = http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to create request", err), nil
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to execute request", err), nil
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("unable to read request response", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Status: %d\nBody: %s", resp.StatusCode, string(respBody))), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
