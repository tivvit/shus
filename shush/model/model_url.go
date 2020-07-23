/*
 * Shush API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package model

import (
	"encoding/json"
	"time"
)

type Url struct {
	ShortUrl string `json:"short_url"`
	Target string `json:"target"`
	Owners *[]string `json:"owners,omitempty"`
	Expiration *time.Time `json:"expiration,omitempty"`
}

// todo NewUrlFromJsonString

func UrlDeserialize(v string) (Url, error) {
	d := Url{}
	err := json.Unmarshal([]byte(v), &d)
	return d, err
}

// todo url method

func UrlSerialize(u Url) (string, error) {
	r, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

// todo
// func (u Url) validate()
// shorturl has to be in some format ?
// url should be valid?


// todo
// func (u *Url) generateShort(g shortner)


// todo
// func (u *Url) store()
