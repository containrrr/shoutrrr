//go:build conduit_e2e
// +build conduit_e2e

package matrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestConduitE2E(t *testing.T) {
	if os.Getenv("SHOUTRRR_MATRIX_CONDUIT_E2E") != "1" {
		t.Skip("set SHOUTRRR_MATRIX_CONDUIT_E2E=1 to run the Conduit e2e test")
	}

	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatalf("docker is required for Conduit e2e test: %v", err)
	}

	image := os.Getenv("SHOUTRRR_MATRIX_CONDUIT_IMAGE")
	if image == "" {
		image = "docker.io/matrixconduit/matrix-conduit:latest"
	}
	port := os.Getenv("SHOUTRRR_MATRIX_CONDUIT_PORT")
	if port == "" {
		port = "6167"
	}

	configPath := filepath.Join(t.TempDir(), "conduit.toml")
	config := `
[global]
server_name = "localhost"
address = "0.0.0.0"
port = 6167
database_backend = "rocksdb"
database_path = "/var/lib/matrix-conduit"
allow_registration = true
allow_federation = false
trusted_servers = []
max_request_size = 20000000
`
	if err := os.WriteFile(configPath, []byte(config), 0600); err != nil {
		t.Fatalf("write Conduit config: %v", err)
	}

	container := fmt.Sprintf("shoutrrr-conduit-e2e-%d", time.Now().UnixNano())
	runDocker(t, "run", "--rm", "-d",
		"--name", container,
		"-p", "127.0.0.1:"+port+":6167",
		"-v", configPath+":/etc/conduit.toml:ro",
		"-e", "CONDUIT_CONFIG=/etc/conduit.toml",
		image)
	t.Cleanup(func() {
		_ = exec.Command("docker", "rm", "-f", container).Run()
	})

	baseURL := "http://127.0.0.1:" + port
	waitForConduit(t, container, baseURL)
	user := "shoutrrr"
	password := "shoutrrr-password"
	deviceID := "shoutrrr-e2e"
	registration := conduitRegister(t, baseURL, user, password)
	roomID := conduitCreateRoom(t, baseURL, registration.AccessToken)

	values := url.Values{}
	values.Set("disableTLS", "yes")
	values.Set("rooms", roomID)
	values.Set("deviceID", deviceID)
	serviceURL := url.URL{
		Scheme:   Scheme,
		Host:     strings.TrimPrefix(baseURL, "http://"),
		User:     url.UserPassword(user, password),
		RawQuery: values.Encode(),
	}

	for i := 0; i < 3; i++ {
		service := &Service{}
		if err := service.Initialize(&serviceURL, nil); err != nil {
			t.Fatalf("initialize Matrix service: %v", err)
		}
		if err := service.Send(fmt.Sprintf("Conduit e2e test %d", i+1), nil); err != nil {
			t.Fatalf("send Matrix message: %v", err)
		}
	}

	devices := conduitDevices(t, baseURL, registration.AccessToken)
	count := 0
	for _, device := range devices.Devices {
		if device.DeviceID == deviceID {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected exactly one Matrix device with ID %q, got %d; devices: %+v", deviceID, count, devices.Devices)
	}
}

func runDocker(t *testing.T, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("docker %s failed: %v\n%s", strings.Join(args, " "), err, output)
	}
	return strings.TrimSpace(string(output))
}

func waitForConduit(t *testing.T, container string, baseURL string) {
	t.Helper()

	deadline := time.Now().Add(45 * time.Second)
	for time.Now().Before(deadline) {
		res, err := http.Get(baseURL + "/_matrix/client/versions")
		if err == nil {
			_ = res.Body.Close()
			if res.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for Conduit container %s at %s\n%s", container, baseURL, dockerLogs(container))
}

func dockerLogs(container string) string {
	cmd := exec.Command("docker", "logs", container)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("docker logs failed: %v\n%s", err, output)
	}
	return string(output)
}

type conduitRegisterResponse struct {
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
	UserID      string `json:"user_id"`
}

func conduitRegister(t *testing.T, baseURL string, user string, password string) conduitRegisterResponse {
	t.Helper()

	request := map[string]interface{}{
		"username": user,
		"password": password,
		"auth": map[string]string{
			"type": "m.login.dummy",
		},
	}

	response := conduitRegisterResponse{}
	doMatrixJSON(t, http.MethodPost, baseURL+"/_matrix/client/r0/register?kind=user", "", request, &response)
	if response.AccessToken == "" {
		t.Fatalf("Conduit registration did not return an access token: %+v", response)
	}
	return response
}

func conduitCreateRoom(t *testing.T, baseURL string, accessToken string) string {
	t.Helper()

	request := map[string]string{
		"name":   "Shoutrrr Conduit E2E",
		"preset": "private_chat",
	}

	response := struct {
		RoomID string `json:"room_id"`
	}{}
	doMatrixJSON(t, http.MethodPost, baseURL+"/_matrix/client/r0/createRoom", accessToken, request, &response)
	if response.RoomID == "" {
		t.Fatalf("Conduit createRoom did not return a room ID: %+v", response)
	}
	return response.RoomID
}

type conduitDevicesResponse struct {
	Devices []struct {
		DeviceID string `json:"device_id"`
	} `json:"devices"`
}

func conduitDevices(t *testing.T, baseURL string, accessToken string) conduitDevicesResponse {
	t.Helper()

	response := conduitDevicesResponse{}
	doMatrixJSON(t, http.MethodGet, baseURL+"/_matrix/client/r0/devices", accessToken, nil, &response)
	return response
}

func doMatrixJSON(t *testing.T, method string, endpoint string, accessToken string, request interface{}, response interface{}) {
	t.Helper()

	var body io.Reader
	if request != nil {
		data, err := json.Marshal(request)
		if err != nil {
			t.Fatalf("marshal Matrix request: %v", err)
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		t.Fatalf("create Matrix request: %v", err)
	}
	if request != nil {
		req.Header.Set("Content-Type", contentType)
	}
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("perform Matrix request %s %s: %v", method, endpoint, err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read Matrix response: %v", err)
	}
	if res.StatusCode >= 400 {
		t.Fatalf("Matrix request %s %s failed with %s: %s", method, endpoint, res.Status, data)
	}
	if response != nil && len(data) > 0 {
		if err := json.Unmarshal(data, response); err != nil {
			t.Fatalf("decode Matrix response: %v\n%s", err, data)
		}
	}
}
