package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubeants.io/config"
)

// 定义User资源的GVR
var userGVR = schema.GroupVersionResource{
	Group:    "user.kubeants.io",
	Version:  "v1beta1",
	Resource: "users",
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// 检测是否为 Bearer Token（SSO 认证）
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if validateSSOToken(token) { // 这里实现 SSO 认证逻辑
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid SSO token"})
			return
		}

		// 检测是否为 Basic 认证
		if strings.HasPrefix(authHeader, "Basic ") {
			decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid basic auth encoding"})
				return
			}

			// 提取用户名和密码
			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) != 2 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid basic auth format"})
				return
			}

			username, password := credentials[0], credentials[1]

			// 校验用户名密码
			isTrue, saToken, err := validateLocalUser(username, password)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if isTrue {
				c.Set("username", username) // 存入 Context 供后续使用
				// ✅ 代理用户请求，使用 `SA Token`
				c.Set("Authorization", "Bearer "+saToken)
				c.Next()
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// 其他情况拒绝访问
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication method"})
	}
}

// validateLocalUser 本地用户认证，校验 bcrypt 哈希密码，并返回 ServiceAccount Token
func validateLocalUser(username, password string) (isTrue bool, token string, err error) {
	// 获取用户资源
	user, err := config.KubeDynamicClient.Resource(userGVR).Get(context.TODO(), username, v1.GetOptions{})
	if err != nil {
		fmt.Println("Failed to get user for k8s:", err)
		return false, "", fmt.Errorf("failed to get user for k8s: %w", err)
	}

	// 获取 bcrypt 加密后的密码哈希
	passwordHash, found, _ := unstructured.NestedString(user.UnstructuredContent(), "spec", "password")
	if !found {
		fmt.Printf("Password not found in user spec,k8s资源user[%v]中没有对应password字段", user)
		return false, "", fmt.Errorf("password not found in user spec,k8s资源user[%v]中没有对应password字段: %w", username, err)
	}

	// 使用 bcrypt 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		fmt.Println("Password verification failed[密码不正确，请检查您的密码]:", err)
		return false, "", fmt.Errorf("password verification failed[密码不正确，请检查您的密码]: %w", err)
	}

	// 获取 ServiceAccount 名称
	saName, found, _ := unstructured.NestedString(user.UnstructuredContent(), "status", "serviceAccount")
	if !found || saName == "" {
		return false, "", fmt.Errorf("serviceAccount not found in user status")
	}

	// 查询 ServiceAccount 的 Secret，获取 Token
	token, err = getSAToken(saName)
	if err != nil {
		return false, "", err
	}

	return true, token, nil
}

// getSAToken 获取 ServiceAccount 的 Token
func getSAToken(saName string) (token string, err error) {
	// 获取sa
	sa, err := config.KubeClientSet.CoreV1().ServiceAccounts("default").Get(context.TODO(), saName, v1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get ServiceAccount: %w", err)
	}

	// 获取关联的Secret
	for _, secret := range sa.Secrets {
		secretObj, err := config.KubeClientSet.CoreV1().Secrets("default").Get(context.TODO(), secret.Name, v1.GetOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to get secret: %w", err)
		}
		if token, exists := secretObj.Data["token"]; exists {
			return string(token), nil
		}
	}
	return "", fmt.Errorf("failed to find token for ServiceAccount %s", saName)
}

// validateSSOToken 这里实现你的 SSO 认证逻辑（如 OIDC 或 OAuth2）
func validateSSOToken(token string) bool {
	// 假设这里是调用 OIDC 或 OAuth2 服务器的逻辑
	// 这里只是一个示例，实际中需要向 SSO 服务器发送请求验证 token
	if token == "valid_token_example" { // 这里需要替换成实际的 Token 验证逻辑
		return true
	}
	return false
}
