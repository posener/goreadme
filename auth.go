package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var githubSecret = os.Getenv("GITHUB_SECRET")

func auth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sig := r.Header.Get("X-Hub-Signature")
		if sig == "" {
			logrus.Errorf("401 - Empty signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		b := bytes.NewBuffer(nil)
		_, err := b.ReadFrom(r.Body)
		if err != nil {
			logrus.Errorf("Failed reading body")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		r.Body.Close()
		r.Body = ioutil.NopCloser(b)
		if !verifySignature([]byte(githubSecret), sig, b.Bytes()) {
			logrus.Infof("401 - Signature did not matched")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		inner.ServeHTTP(w, r)
	})
}

func verifySignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}
