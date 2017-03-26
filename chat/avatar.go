package chat

import (
	"errors"
	"fmt"
)

var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar  URL.")

type Avatar interface {
	GetAvatarURL(self *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar"]
	if !ok {
		return "", ErrNoAvatarURL
	}

	urlStr, ok := url.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}

	return urlStr, nil
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	userId, ok := c.userData["userId"]
	if !ok {
		return "", ErrNoAvatarURL
	}

	emailStr, _ := userId.(string) // I no check if it ok because it should be 100% string if it in userData map
	return fmt.Sprintf("//www.gravatar.com/avatar/%s", emailStr), nil
}
