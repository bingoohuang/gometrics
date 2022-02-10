package metric

import (
	"github.com/bingoohuang/gometrics/pkg/util"
	"github.com/sirupsen/logrus"
)

// Key defines a slice of keys.
type Key struct {
	Keys    []string
	Checked bool
}

// NewKey create Keys.
func NewKey(keys []string) Key {
	ks := Key{Keys: keys, Checked: false}
	ks.Check()

	return ks
}

// Check checks the validation of keys.
func (k *Key) Check() {
	k.Checked = true

	if len(k.Keys) == 0 {
		k.Checked = false
		logrus.Warn("Keys required")
		return
	}

	for i, key := range k.Keys {
		if !k.validateKey(i, key) {
			k.Checked = false
		}
	}
}

const strippedChars = `" .,|#\` + "\t\r\n"

func (k *Key) validateKey(i int, key string) bool {
	if key == "" {
		logrus.Warn("Key can not be empty")
		return false
	}

	key = util.StripAny(key, strippedChars)
	if key == "" {
		logrus.Warnf("invalid Key %s", key)
		return false
	}

	k.Keys[i] = util.Abbr(key, 100, "...")
	return true
}
