package util_test

import (
	"testing"

	"github.com/bingoohuang/gometrics/util"

	"github.com/stretchr/testify/assert"
)

func TestStripAny(t *testing.T) {
	str := "你好吗? 我好! 好我好!? 你好好!"
	stripped := util.StripAny(str, "我好") // now with remove/strip a set of unicode characters
	assert.Equal(t, "你吗? ! !? 你!", stripped)

	str = "Happy Go Lucky!"
	stripped = util.StripAny(str, "aGo") // will work with a set of characters
	assert.Equal(t, "Hppy  Lucky!", stripped)
}

func TestEsc(t *testing.T) {
	j := util.Esc("\"\\\r\n")
	assert.Equal(t, `\"\\\r\n`, j)
}

func TestAbbr(t *testing.T) {
	assert.Equal(t, "abc", util.Abbr("abc", 100, ""))
	assert.Equal(t, "a..", util.Abbr("abcd", 3, ".."))
	assert.Equal(t, "...", util.Abbr("abcd", 3, "..."))
}
