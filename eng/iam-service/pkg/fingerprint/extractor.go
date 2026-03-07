package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/mssola/user_agent"
)

type DeviceInfo struct {
	IPAddress  string
	UserAgent  string
	DeviceType string
	Browser    string
	OS         string
	Hash       string
}

func Extract(req *http.Request) DeviceInfo {
	userAgentStr := req.UserAgent()
	ua := user_agent.New(userAgentStr)

	deviceType := "Desktop"
	if ua.Mobile() {
		deviceType = "Mobile"
	}
	if ua.Bot() {
		deviceType = "Bot"
	}

	browserName, browserVersion := ua.Browser()
	osInfo := ua.OS()

	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.RemoteAddr
	}
	ip = strings.Split(ip, ":")[0]

	rawSignature := fmt.Sprintf("%s|%s|%s", ip, osInfo, browserName)
	hashBytes := sha256.Sum256([]byte(rawSignature))
	deviceHash := hex.EncodeToString(hashBytes[:])

	return DeviceInfo{
		IPAddress:  ip,
		UserAgent:  userAgentStr,
		DeviceType: deviceType,
		Browser:    fmt.Sprintf("%s %s", browserName, browserVersion),
		OS:         osInfo,
		Hash:       deviceHash,
	}
}
