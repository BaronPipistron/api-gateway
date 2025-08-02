package proxy

import (
	"github.com/BaronPipistron/api-gateway/internal/config"
	"github.com/BaronPipistron/api-gateway/internal/telemetry/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ProxyHandler struct {
	Rules []config.RuleConfig
}

func NewProxyHandler(cfg *config.Config) *ProxyHandler {
	return &ProxyHandler{
		Rules: cfg.Rules,
	}
}

func (ph *ProxyHandler) Handle(c *gin.Context) {
	rule := ph.matchRule(c.Request.URL.Path)
	if rule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching proxy rule"})
		return
	}

	for _, header := range rule.HeadersRequired {
		if c.GetHeader(header) == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Missing required header: " + header})
			return
		}
	}

	newHeaders := http.Header{}
	for _, h := range rule.AllowedHeaders {
		if val := c.GetHeader(h); val != "" {
			newHeaders.Set(h, val)
		}
	}

	targetURL, err := url.Parse(rule.RedirectTo)
	if err != nil {
		logging.Error("Invalid redirect URL:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	c.Request.URL.Scheme = targetURL.Scheme
	c.Request.URL.Host = targetURL.Host
	c.Request.Host = targetURL.Host
	c.Request.Header = newHeaders

	proxy.ServeHTTP(c.Writer, c.Request)
}

func (ph *ProxyHandler) matchRule(path string) *config.RuleConfig {
	for _, rule := range ph.Rules {
		if strings.HasPrefix(path, rule.From) {
			return &rule
		}
	}
	return nil
}
