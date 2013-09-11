// Copyright 2013 beebbs authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package routers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/beego/beebbs/models"
)

func registerUser(form RegisterForm) error {
	salt := getRandomString(10)
	pwd := encodePassword(form.Password, salt)

	user := &models.User{}
	user.UserName = form.UserName
	user.Email = form.Email
	user.Password = fmt.Sprintf("%s$%s", salt, pwd)
	user.GrEmail = encodeMd5(form.Email)

	return models.NewUser(user)
}

func verifyUser(username, password string, user *models.User) bool {
	qs := orm.NewOrm().QueryTable("user")
	if strings.Index(username, "@") == -1 {
		qs = qs.Filter("UserName", username)
	} else {
		qs = qs.Filter("Email", username)
	}
	err := qs.One(user)
	if err != nil {
		return false
	}
	if verifyPassword(password, user.Password) {
		return true
	}
	return false
}

func verifyPassword(password, encoded string) bool {
	salt := encoded[:10]
	return encodePassword(password, salt) == encoded[11:]
}

func encodePassword(password string, salt string) string {
	pwd := PBKDF2([]byte(password), []byte(salt), 10000, 50, sha256.New)
	return hex.EncodeToString(pwd)
}
