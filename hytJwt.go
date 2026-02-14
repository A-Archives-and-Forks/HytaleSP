package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)
var jwtPublic, jwtPrivate, _ = ed25519.GenerateKey(rand.Reader);


type jwtHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}

type identityToken struct {
	Exp     int     `json:"exp"`
	Iat     int     `json:"iat"`
	Iss     string  `json:"iss"`
	Jti     string  `json:"jti"`
	Scope   string  `json:"scope"`
	Sub     string  `json:"sub"`
	Profile profileInfo `json:"profile"`
}


type sessionToken struct {
	Exp   int    `json:"exp"`
	Iat   int    `json:"iat"`
	Iss   string `json:"iss"`
	Jti   string `json:"jti"`
	Scope string `json:"scope"`
	Sub   string `json:"sub"`
}

type profileInfo struct {
	Username     string   `json:"username"`
	Entitlements []string `json:"entitlements"`
	Skin         string   `json:"skin"`
}

type jwkKeyList struct {
	Keys []jwkKey `json:"keys"`
}
type jwkKey struct {
	Alg string `json:"alg"`
	Crv string `json:"crv"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	X   string `json:"x"`
}

// jwt related functions ...

func signJwt(j string) string {

	sig := ed25519.Sign(jwtPrivate, []byte(j));
	return base64.RawURLEncoding.EncodeToString(sig);
}

func makeJwt(body any) string {
	head := jwtHeader{
		Alg: "EdDSA",
		Kid: "2025-10-01",
		Typ: "JWT",
	};

	jHead, _ := json.Marshal(head);
	jBody, _ := json.Marshal(body);


	jwt := base64.RawURLEncoding.EncodeToString(jHead) + "." + base64.RawURLEncoding.EncodeToString(jBody)
	jwt += "." + signJwt(jwt);
	return jwt;
}

func unmakeJwt(jwt string) (jwtHeader, sessionToken, error) {
	jwtParts := strings.Split(jwt, ".");
	if len(jwtParts) < 2 {
		return jwtHeader{}, sessionToken{}, fmt.Errorf("does not contain atleast 2 parts.");
	}

	jHead, err:= base64.RawURLEncoding.DecodeString(jwtParts[0]);
	if err != nil {
		return jwtHeader{}, sessionToken{}, err;
	}
	jPayload, err := base64.RawURLEncoding.DecodeString(jwtParts[0]);
	if err != nil {
		return jwtHeader{}, sessionToken{}, err;
	}


	head := jwtHeader{};
	session := sessionToken{};

	err = json.Unmarshal(jHead, &head);
	if err != nil {
		return jwtHeader{}, sessionToken{}, err;
	}

	err = json.Unmarshal(jPayload, &session);
	if err != nil {
		return jwtHeader{}, sessionToken{}, err;
	}


	return head, session, nil;
}
