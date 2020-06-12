package main

import (
	"errors"
// 	"io"
// 	"crypto/md5"
// 	"fmt"
// 	"strings"
)

// ErrorAvatarインスタンスがアバターのURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できませんでした。")

type Avatar interface { // ユーザのプロフィール画像を表す型
	GetAvatarURL(c *client) (string, error) // 指定されたクライアントのアバターのURLを返す。問題が発生：エラー、取得できない：ErrorAvatarURL
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}