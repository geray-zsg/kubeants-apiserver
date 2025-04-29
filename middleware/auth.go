// /middleware/auth.go
package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubeants.io/config"
	"kubeants.io/models"
	"kubeants.io/response"
	"kubeants.io/service"
	"kubeants.io/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// 动态获取 JWT Secret 和 Expiration 配置
func getJWTKey() []byte {
	return []byte(config.CONF.JWT.Secret)
}
func getJWTExpiration() int {
	return config.CONF.JWT.Expiration
}

// 定义user资源的GVR
var userGVR = schema.GroupVersionResource{
	Group:    "user.kubeants.io",
	Version:  "v1beta1",
	Resource: "users",
}

// User represents a simple user model
type User struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// LoginHandler 用户登录接口
// @Summary 用户登录接口: 验证用户登录信息并返回 JWT。
// @Description 验证用户输入的用户名和密码，如果正确，则返回 JWT。
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body User true "用户登录信息"
// @Success 200 {object} map[string]interface{} "登陆成功"
// @Router /gapi/login [post]
func LoginHandler(c *gin.Context) {
	ctx := context.TODO()
	logger := log.FromContext(ctx)
	// 解析用户请求
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error(err, "Failed to parse login request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证用户身份
	isValid, saToken, err := validateLocalUser(user.Username, user.Password)
	if err != nil {
		logger.Error(err, "Failed to validate user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !isValid {
		logger.Error(err, "Invalid username or password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	// 生成JWT Token
	token, err := generateJWTToken(user.Username, user.Password, saToken, getJWTExpiration())
	if err != nil {
		logger.Error(err, "Failed to generate JWT token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回token
	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		// "sa_token": saToken,

		"token": token,
	})
	logger.Info("User login", "current login user", user.Username)
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	logger := log.FromContext(context.TODO())
	logger.Info("====> AuthMiddleware is called")
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error(gin.Error{Err: fmt.Errorf("authorization header is empty")}, "[认证失败] Authorization header is empty")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// 处理 Basic 认证
		if strings.HasPrefix(authHeader, "Basic ") {
			// Base64 解码
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

			// 验证用户名密码
			isValid, saToken, err := validateLocalUser(username, password)
			if err != nil {
				logger.Error(err, "Failed to validate user")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
				return
			}
			if !isValid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
				return
			}

			// 生成token
			token, err := generateJWTToken(username, password, saToken, getJWTExpiration())
			if err != nil {
				logger.Error(err, "Failed to generate JWT token")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
				return
			}
			// 认证通过，设置上下文并添加 SA Token
			c.Set("username", username)
			c.Set("sa_token", saToken) // 存储 sa_token 以便后续检查权限
			c.Set("Authorization", "Bearer "+token)
			c.Next()
			return
		}

		// 处理 Bearer Token 认证（原有逻辑）
		if strings.HasPrefix(authHeader, "Bearer ") {
			// 确保 Authorization 头是 Bearer 方式
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Error(gin.Error{Err: fmt.Errorf("authorization header is not Bearer")}, "[认证失败] 请检查token是否是 Bearer 方式: Invalid token, please check if the token is Bearer format")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token, please check if the token is Bearer format"})
				return
			}

			tokenString := parts[1] // 获取实际的 JWT
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return getJWTKey(), nil
			})

			if err != nil || !token.Valid {
				logger.Error(err, "[认证失败] Invalid token,无效或者过期token")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token,无效或者过期token"})
				c.Abort()
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				logger.Error(gin.Error{Err: fmt.Errorf("invalid token claims")}, "[认证失败] Invalid token claims")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				return
			}

			username := claims["username"].(string)
			saToken := claims["sa_token"].(string) // 从 JWT 解析出 ServiceAccount Token

			c.Set("username", username)
			c.Set("sa_token", saToken) // 存储 sa_token 以便后续检查权限

			c.Next()
			return
		}
		// 其他认证方式拒绝
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unsupported authentication method"})
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
		fmt.Printf("✅ Password not found in user spec,k8s资源user[%v]中没有对应password字段", user)
		return false, "", fmt.Errorf("❌ password not found in user spec,k8s资源user[%v]中没有对应password字段: %w", username, err)
	}

	// 使用 bcrypt 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		fmt.Println("✅ Password verification failed[密码不正确，请检查您的密码]:", err)
		return false, "", fmt.Errorf("❌ password verification failed[密码不正确，请检查您的密码]: %w", err)
	}

	// 获取 ServiceAccount 名称
	saName, found, _ := unstructured.NestedString(user.UnstructuredContent(), "status", "serviceAccount")
	if !found || saName == "" {
		return false, "", fmt.Errorf("❌ serviceAccount not found in user status")
	}
	// 获取 ServiceAccount 所在namespace
	saNamespace, found, _ := unstructured.NestedString(user.UnstructuredContent(), "status", "serviceAccount")
	if !found || saNamespace == "" {
		fmt.Printf("✅ user[%v] status 中没有对应serviceAccount namespace字段，默认使用kubeants-system", username)
		saNamespace = "kubeants-system"
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
	sa, err := config.KubeClientSet.CoreV1().ServiceAccounts("kubeants-system").Get(context.TODO(), saName, v1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get ServiceAccount: %w", err)
	}

	// 获取关联的Secret
	for _, secret := range sa.Secrets {
		secretObj, err := config.KubeClientSet.CoreV1().Secrets("kubeants-system").Get(context.TODO(), secret.Name, v1.GetOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to get secret: %w", err)
		}
		if token, exists := secretObj.Data["token"]; exists {
			return string(token), nil
		}
	}
	return "", fmt.Errorf("failed to find token for ServiceAccount %s", saName)
}

// 生成 JWT Token
func generateJWTToken(username, password, saToken string, expiration int) (string, error) {
	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"password": password,
		"sa_token": saToken, // 存入 ServiceAccount Token
		"exp":      time.Now().Add(time.Duration(expiration) * time.Second).Unix(),
	})

	// 使用密钥签名
	return token.SignedString(getJWTKey())
}

// GetUserInfo 用户信息
// @Summary 认证通过后返回用户信息。
// @Description 验证用户输入的用户名和密码，如果正确，则返回 JWT。
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "登陆成功"
// @Router /gapi/user/info/{username} [get]
func GetUserInfo(c *gin.Context) {
	ctx := context.TODO()
	logger := log.FromContext(ctx)
	// 从请求载荷中提取username
	username := c.Param("username")

	gvr := schema.GroupVersionResource{
		Group:    "user.kubeants.io",
		Version:  "v1beta1",
		Resource: "users",
	}

	logger.Info("获取用户信息", "用户名", username)
	userUnstructured, err := service.ServiceGroupApp.ResourceServiceGroup.GetResource(ctx, "", gvr.Group, gvr.Version, gvr.Resource, "", username)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	// 将 userunstructured 转换为 User 结构体
	var user models.User
	if err := util.UnstructuredToStruct(userUnstructured, &user); err != nil {
		logger.Error(err, "结构转换失败")
		return
	}

	// 移除user中的password
	user.Spec.Password = ""

	response.SuccessWithDetailed(c, "用户信息获取成功", user)
}
