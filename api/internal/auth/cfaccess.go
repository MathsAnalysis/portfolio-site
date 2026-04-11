// Package auth implements Cloudflare Access JWT verification.
//
// When a user accesses an app protected by Cloudflare Access, CF puts a signed
// JWT in the header Cf-Access-Jwt-Assertion and also CF-Authorization cookie.
// The JWT is signed with RS256 using one of the keys published at
//   https://<team>.cloudflareaccess.com/cdn-cgi/access/certs
//
// We verify:
//   - signature against the JWKS (cached with TTL)
//   - issuer matches https://<team>.cloudflareaccess.com
//   - audience matches the app AUD tag (from CF Access application settings)
//   - exp / iat / nbf
//
// On success we return the authenticated email. On failure, the caller rejects
// the request with 401.
package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CFAccess struct {
	teamDomain string // e.g. "mathsanalysis.cloudflareaccess.com"
	audience   string // app AUD tag
	issuer     string // "https://mathsanalysis.cloudflareaccess.com"

	mu        sync.RWMutex
	keys      map[string]*rsa.PublicKey
	keysAt    time.Time
	keysTTL   time.Duration
	httpClient *http.Client
}

type Config struct {
	TeamDomain string // e.g. "mathsanalysis.cloudflareaccess.com"
	Audience   string
}

func New(cfg Config) *CFAccess {
	return &CFAccess{
		teamDomain: cfg.TeamDomain,
		audience:   cfg.Audience,
		issuer:     "https://" + cfg.TeamDomain,
		keys:       map[string]*rsa.PublicKey{},
		keysTTL:    15 * time.Minute,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Verify returns the authenticated email on success.
func (c *CFAccess) Verify(ctx context.Context, tokenStr string) (string, error) {
	if tokenStr == "" {
		return "", errors.New("empty token")
	}

	token, err := jwt.Parse(tokenStr, c.keyFunc(ctx),
		jwt.WithAudience(c.audience),
		jwt.WithIssuer(c.issuer),
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}
	email, _ := claims["email"].(string)
	if email == "" {
		return "", errors.New("missing email claim")
	}
	return email, nil
}

func (c *CFAccess) keyFunc(ctx context.Context) jwt.Keyfunc {
	return func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("missing kid")
		}
		key, err := c.keyByID(ctx, kid)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
}

func (c *CFAccess) keyByID(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	c.mu.RLock()
	if k, ok := c.keys[kid]; ok && time.Since(c.keysAt) < c.keysTTL {
		c.mu.RUnlock()
		return k, nil
	}
	c.mu.RUnlock()

	if err := c.refresh(ctx); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	if k, ok := c.keys[kid]; ok {
		return k, nil
	}
	return nil, fmt.Errorf("unknown kid %q", kid)
}

type jwksResponse struct {
	Keys []struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		Use string `json:"use"`
		Alg string `json:"alg"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

func (c *CFAccess) refresh(ctx context.Context) error {
	url := c.issuer + "/cdn-cgi/access/certs"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch jwks: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jwks status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	if err != nil {
		return err
	}
	var jr jwksResponse
	if err := json.Unmarshal(body, &jr); err != nil {
		return fmt.Errorf("parse jwks: %w", err)
	}

	next := map[string]*rsa.PublicKey{}
	for _, k := range jr.Keys {
		if k.Kty != "RSA" {
			continue
		}
		nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
		if err != nil {
			continue
		}
		pub := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: int(new(big.Int).SetBytes(eBytes).Int64()),
		}
		next[k.Kid] = pub
	}
	if len(next) == 0 {
		return errors.New("jwks: no RSA keys")
	}

	c.mu.Lock()
	c.keys = next
	c.keysAt = time.Now()
	c.mu.Unlock()
	return nil
}
