package chat

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	imAppID  string
	imSecret string
}

func NewHandler(imAppID, imSecret string) *Handler {
	return &Handler{imAppID: imAppID, imSecret: imSecret}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/chat/user-sig", h.GetUserSig)
}

// @Summary      获取IM UserSig
// @Tags         即时通讯
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /chat/user-sig [get]
func (h *Handler) GetUserSig(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if h.imAppID == "" || h.imSecret == "" {
		c.JSON(http.StatusOK, gin.H{"user_sig": "", "warning": "IM not configured, set IM_APP_ID and IM_SECRET env vars"})
		return
	}

	appID, _ := strconv.Atoi(h.imAppID)
	sig, err := genUserSig(appID, h.imSecret, userID.(string), 86400*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate user sig"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_sig": sig, "im_app_id": appID})
}

func genUserSig(appID int, secret, userID string, expire int) (string, error) {
	now := time.Now().Unix()
	sigDoc := map[string]interface{}{
		"TLS.ver":        "2.0",
		"TLS.identifier": userID,
		"TLS.sdkappid":   appID,
		"TLS.expire":     expire,
		"TLS.time":       now,
	}
	data, err := json.Marshal(sigDoc)
	if err != nil {
		return "", err
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(b64))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	result := fmt.Sprintf("%s.%s.%s.%s",
		trimPad(b64),
		trimPad(sig),
		"", "",
	)
	return result, nil
}

func trimPad(s string) string {
	for len(s) > 0 && s[len(s)-1] == '=' {
		s = s[:len(s)-1]
	}
	return s
}
