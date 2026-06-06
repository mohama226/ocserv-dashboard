package customer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	ocservUser "github.com/mmtaee/ocserv-dashboard/common/ocserv/user"
	"github.com/mmtaee/ocserv-dashboard/common/pkg/config"
)

const iosSetupCertificateTokenTTL = 10 * time.Minute

// IOSSetup creates Cisco Secure Client iOS setup URIs for the customer.
//
// @Summary      Create customer iOS setup links
// @Description  Create Cisco Secure Client iOS certificate import and connection creation URIs using ocserv username/password
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        request body SummaryData true "customer username and password (same ocserv account)."
// @Failure      400 {object} request.ErrorResponse
// @Failure      429 {object} middlewares.TooManyRequests
// @Success      200 {object} IOSSetupResponse
// @Router       /customers/setup/ios [post]
func (ctl *Controller) IOSSetup(c echo.Context) error {
	var data SummaryData

	if err := ctl.request.DoValidate(c, &data); err != nil {
		return ctl.request.BadRequest(c, err)
	}

	if data.Password == "Secret-Ocpasswd" {
		return ctl.request.BadRequest(c, errors.New("invalid username or password"))
	}

	user, err := ctl.ocservUserRepo.GetByUsername(c.Request().Context(), data.Username)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	if user.Password != data.Password {
		return ctl.request.BadRequest(c, errors.New("invalid username or password"))
	}

	if !user.CertificateAvailable {
		return ctl.request.BadRequest(c, errors.New("certificate is not available for this user"))
	}

	systemConfig, err := ctl.systemRepo.System(c.Request().Context())
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	connectionName, err := ocservUser.NormalizeProfileConnectionName(systemConfig.ClientProfileConnectionName)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	serverAddress, err := ocservUser.NormalizeProfileServerAddress(systemConfig.ClientProfileServerAddress)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	serverPort, err := ocservUser.NormalizeProfileServerPort(systemConfig.ClientProfileServerPort)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	expiresAt := time.Now().Add(iosSetupCertificateTokenTTL)
	token, err := createIOSSetupCertificateToken(user.Username, expiresAt)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	certificateURL := publicAPIBaseURL(c) + "/api/customers/setup/ios/certificate/" + url.PathEscape(token)

	certificateImportURI, err := ocservUser.BuildAnyConnectImportURI(certificateURL)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	connectionCreateURI, err := ocservUser.BuildAnyConnectCreateURI(connectionName, serverAddress, serverPort, user.Username)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	return c.JSON(http.StatusOK, IOSSetupResponse{
		CertificateImportURI: certificateImportURI,
		ConnectionCreateURI:  connectionCreateURI,
		CertificatePassword:  user.Password,
		ConnectionName:       connectionName,
		ServerAddress:        serverAddress,
		ServerPort:           serverPort,
		ExpiresAt:            expiresAt,
	})
}

// DownloadIOSSetupCertificate downloads the customer's certificate through a short-lived signed setup token.
//
// @Summary      Download customer iOS setup certificate
// @Description  Download customer's PKCS#12 certificate using a short-lived iOS setup token
// @Tags         Customers
// @Produce      application/x-pkcs12
// @Param        token path string true "iOS setup certificate token"
// @Failure      400 {object} request.ErrorResponse
// @Failure      429 {object} middlewares.TooManyRequests
// @Success      200 {file} file "user.p12"
// @Router       /customers/setup/ios/certificate/{token} [get]
func (ctl *Controller) DownloadIOSSetupCertificate(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return ctl.request.BadRequest(c, errors.New("token is required"))
	}

	username, err := parseIOSSetupCertificateToken(token)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	path, err := ctl.ocservUserRepo.CertificatePathByUsername(c.Request().Context(), username)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	c.Response().Header().Set(echo.HeaderContentType, "application/x-pkcs12")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-store")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")

	return c.Attachment(path, username+".p12")
}

func publicAPIBaseURL(c echo.Context) string {
	req := c.Request()

	scheme := strings.TrimSpace(req.Header.Get("X-Forwarded-Proto"))
	if scheme == "" {
		scheme = strings.TrimSpace(req.URL.Scheme)
	}
	if scheme == "" {
		if req.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	host := strings.TrimSpace(req.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = strings.TrimSpace(req.Host)
	}

	return scheme + "://" + host
}

func createIOSSetupCertificateToken(username string, expiresAt time.Time) (string, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return "", errors.New("username is required")
	}

	if strings.Contains(username, "|") {
		return "", errors.New("username contains invalid characters")
	}

	payload := username + "|" + strconv.FormatInt(expiresAt.Unix(), 10)

	signature, err := signIOSSetupCertificatePayload(payload)
	if err != nil {
		return "", err
	}

	rawToken := payload + "|" + signature

	return base64.RawURLEncoding.EncodeToString([]byte(rawToken)), nil
}

func parseIOSSetupCertificateToken(token string) (string, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return "", errors.New("token is required")
	}

	rawToken, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", errors.New("invalid token")
	}

	parts := strings.Split(string(rawToken), "|")
	if len(parts) != 3 {
		return "", errors.New("invalid token")
	}

	username := parts[0]
	expiresAtUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", errors.New("invalid token expiry")
	}

	if time.Now().After(time.Unix(expiresAtUnix, 0)) {
		return "", errors.New("token has expired")
	}

	payload := username + "|" + parts[1]

	expectedSignature, err := signIOSSetupCertificatePayload(payload)
	if err != nil {
		return "", err
	}

	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return "", errors.New("invalid token signature")
	}

	return username, nil
}

func signIOSSetupCertificatePayload(payload string) (string, error) {
	secretKey := strings.TrimSpace(config.Get().SecretKey)
	if secretKey == "" {
		return "", errors.New("secret key is not configured")
	}

	mac := hmac.New(sha256.New, []byte(secretKey))
	if _, err := mac.Write([]byte(payload)); err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}
