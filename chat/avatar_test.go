package chat

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("..", "avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer os.Remove(filename)
	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.userData = map[string]interface{}{"userId": "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
