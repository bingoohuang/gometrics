package gometrics_test

import (
	"testing"

	"github.com/bingoohuang/gometrics"
	"github.com/stretchr/testify/assert"
)

func TestStripAny(t *testing.T) {
	str := "你好吗? 我好! 好我好!? 你好好!"
	stripped := gometrics.StripAny(str, "我好") // now with remove/strip a set of unicode characters
	assert.Equal(t, "你吗? ! !? 你!", stripped)

	str = "Happy Go Lucky!"
	stripped = gometrics.StripAny(str, "aGo") // will work with a set of characters
	assert.Equal(t, "Hppy  Lucky!", stripped)
}

func TestEsc(t *testing.T) {
	j := gometrics.Esc("\"\\\r\n")
	assert.Equal(t, `\"\\\r\n`, j)
}

func TestAbbr(t *testing.T) {
	assert.Equal(t, "abc", gometrics.Abbr("abc", 100, ""))
	assert.Equal(t, "a..", gometrics.Abbr("abcd", 3, ".."))
	assert.Equal(t, "...", gometrics.Abbr("abcd", 3, "..."))
}
