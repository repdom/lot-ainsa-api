package middleware

import (
	"be-lotsanmateo-api/internal/config"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
)

// KeycloakConfig contiene la configuración para conectar con Keycloak
type KeycloakConfig struct {
	RealmURL     string   // URL del realm de Keycloak (ej: https://keycloak.example.com/auth/realms/mi-realm)
	ClientID     string   // ID del cliente en Keycloak
	RequiredRole []string // Rol requerido (opcional)
}

// JWKSResponse estructura para la respuesta de JWKS de Keycloak
type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

// JWK estructura para una clave pública
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// KeycloakClaims estructura para los claims de Keycloak
type KeycloakClaims struct {
	jwt.RegisteredClaims
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	Name              string `json:"name"`
}

func (config KeycloakConfig) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización requerido"})
			c.Abort()
			return
		}

		// Verificar que el header tenga el formato correcto
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		// Validar y parsear el token
		claims, err := validateToken(tokenString, config)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		//// Verificar rol requerido si está configurado
		//if len(config.RequiredRole) > 0 && !hasRequiredRole(claims, config) {
		//	c.JSON(http.StatusForbidden, gin.H{"error": "Rol insuficiente"})
		//	c.Abort()
		//	return
		//}

		// Añadir claims al contexto
		c.Set("user_claims", claims)
		c.Set("user_id", claims.Subject)
		c.Set("username", claims.PreferredUsername)
		c.Set("email", claims.Email)

		c.Next()
	}
}

func validateToken(tokenString string, config KeycloakConfig) (*KeycloakClaims, error) {
	// Parsear el token para obtener el header
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &KeycloakClaims{})
	if err != nil {
		return nil, fmt.Errorf("error al parsear token: %v", err)
	}

	// Obtener el kid del header
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("kid no encontrado en el header del token")
	}

	// Obtener la clave pública correspondiente
	publicKey, err := getPublicKey(config.RealmURL, kid)
	if err != nil {
		return nil, fmt.Errorf("error al obtener clave pública: %v", err)
	}

	// Validar el token con la clave pública
	parsedToken, err := jwt.ParseWithClaims(tokenString, &KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token inválido: %v", err)
	}

	claims, ok := parsedToken.Claims.(*KeycloakClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("claims de token inválidos")
	}

	return claims, nil
}

func getPublicKey(realmURL, kid string) (*rsa.PublicKey, error) {
	jwksURL := fmt.Sprintf("%s/protocol/openid-connect/certs", realmURL)

	log.Print(jwksURL)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("error al obtener JWKS: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	var jwksResp JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwksResp); err != nil {
		return nil, fmt.Errorf("error al decodificar JWKS: %v", err)
	}

	// Buscar la clave con el kid correcto
	for _, key := range jwksResp.Keys {
		if key.Kid == kid {
			return jwkToRSAKey(key)
		}
	}

	return nil, fmt.Errorf("clave con kid %s no encontrada", kid)
}

func jwkToRSAKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decodificar N (módulo)
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar N: %v", err)
	}

	// Decodificar E (exponente)
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar E: %v", err)
	}

	// Convertir bytes a big.Int
	n := new(big.Int).SetBytes(nBytes)

	// Convertir E a int
	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

func hasRequiredRole(claims *KeycloakClaims, config KeycloakConfig) bool {
	// Verificar roles del cliente específico debe de contar con uno de los roles requeridos
	if clientAccess, exists := claims.ResourceAccess[config.ClientID]; exists {
		for _, role := range clientAccess.Roles {
			if contains(config.RequiredRole, role) {
				return true
			}
		}
	}

	return false
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v != "" && v == value {
			return true
		}
	}
	return false
}

func NewAuthMiddleware(env *config.Env) KeycloakConfig {
	return KeycloakConfig{
		RealmURL:     env.GetEnv("KEYCLOAK_URL", "https://keycloak.rca-dev.com/realms/master"),
		ClientID:     env.GetEnv("KEYCLOAK_CLIENT_ID", "san-mateo-user"),
		RequiredRole: env.GetEnvArr("KEYCLOAK_REQUIRED_ROLE", []string{"admin", "analista", "caja"}),
	}
}
