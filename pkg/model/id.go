package model

import (
	"crypto/sha1"
	"errors"
	"strings"

	"github.com/google/uuid"
)

func GenerateId(keys ...string) (string, error) {
	if len(keys) == 0 {
		return "", errors.New("keys are empty")
	}
	key := strings.Join(keys, "-")
	data := []byte(key)
	sha1 := sha1.Sum(data)
	return uuid.NewSHA1(uuid.Nil, sha1[:]).String(), nil
}

func GenerateArticleId(article Article) (string, error) {
	return GenerateId(article.title, article.link)
}
