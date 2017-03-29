package chat

import (
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var avatar AuthAvatar
	client := new(client)

	url, err := avatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURLwhen no value present")
	}

	testUrl := "http://url-to-gravatar/"
	client.userData = map[string]interface{}{"avatar": testUrl}

	url, err = avatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no errorwhen value present")
	}

	if url != testUrl {
		t.Error("AuthAvatar.GetAvatarURL should return correct URL")
	}
}
