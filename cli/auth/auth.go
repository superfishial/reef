package auth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/cli/oauth/device"
	"github.com/superfishial/reef/cli/config"
)

func performLoginFlow(conf config.Config) (string, error) {
	code, err := device.RequestCode(
		http.DefaultClient,
		conf.OAuth.DeviceEndpoint,
		conf.OAuth.ClientID,
		[]string{"openid", "profile", "email"},
	)
	if err != nil {
		// TODO: error info
		return "", err
	}

	fmt.Printf("Go to: %s\n", code.VerificationURIComplete)
	fmt.Printf("Enter code (if necessary): %s\n", code.UserCode)

	// Try to open in browser but ignore if it fails
	_ = openInBrowser(code.VerificationURIComplete)

	// Wait for the user to login
	accessToken, err := device.Wait(context.TODO(), http.DefaultClient, conf.OAuth.TokenEndpoint, device.WaitOptions{
		ClientID:   conf.OAuth.ClientID,
		DeviceCode: code,
	})
	if err != nil {
		// TODO: error info
		return "", err
	}

	// Create URL for fetching the signed token
	serverURL, err := url.Parse(conf.Server.RootURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse server url: %w", err)
	}
	serverURL.Path = "/v1/auth/sign"
	serverURLQuery := serverURL.Query()
	serverURLQuery.Add("access_token", accessToken.Token)
	serverURL.RawQuery = serverURLQuery.Encode()

	// Fetch signed token
	resp, err := http.Get(serverURL.String())
	if err != nil {
		return "", fmt.Errorf("failed to fetch signed token: %w", err)
	}
	defer resp.Body.Close()
	signedTokenBytesOrError, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch signed token with status and message: %s \"%s\"", resp.Status, signedTokenBytesOrError)
	}
	if err != nil {
		return "", fmt.Errorf("failed to read body while fetching signed token: %w", err)
	}
	return string(signedTokenBytesOrError), nil
}

func openInBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
