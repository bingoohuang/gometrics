package metric

import (
	"github.com/bingoohuang/gometrics/util"
	"github.com/sirupsen/logrus"
)

// Keys defines a slice of keys
type Keys struct {
	Keys    []string
	Checked bool
}

// NewKeys create Keys
func NewKeys(keys []string) Keys {
	ks := Keys{Keys: keys}
	ks.Check()

	return ks
}

// Check checks the validation of keys
func (k *Keys) Check() {
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

func (k *Keys) validateKey(i int, key string) bool {
	if key == "" {
		logrus.Warn("key can not be empty")

		return false
	}

	key = util.StripAny(key, strippedChars)
	if key == "" {
		logrus.Warnf("invalid key %s", key)

		return false
	}

	k.Keys[i] = util.Abbr(key, 100, "...")

	return true
}
