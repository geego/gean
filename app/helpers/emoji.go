package helpers

import (
	"bytes"
	"sync"

	"github.com/gostores/emoji"
)

var (
	emojiInit sync.Once

	emojis = make(map[string][]byte)

	emojiDelim     = []byte(":")
	emojiWordDelim = []byte(" ")
	emojiMaxSize   int
)

// Emojify "emojifies" the input source.
// Note that the input byte slice will be modified if needed.
// See http://www.emoji-cheat-sheet.com/
func Emojify(source []byte) []byte {
	emojiInit.Do(initEmoji)

	start := 0
	k := bytes.Index(source[start:], emojiDelim)

	for k != -1 {

		j := start + k

		upper := j + emojiMaxSize

		if upper > len(source) {
			upper = len(source)
		}

		endEmoji := bytes.Index(source[j+1:upper], emojiDelim)
		nextWordDelim := bytes.Index(source[j:upper], emojiWordDelim)

		if endEmoji < 0 {
			start++
		} else if endEmoji == 0 || (nextWordDelim != -1 && nextWordDelim < endEmoji) {
			start += endEmoji + 1
		} else {
			endKey := endEmoji + j + 2
			emojiKey := source[j:endKey]

			if emoji, ok := emojis[string(emojiKey)]; ok {
				source = append(source[:j], append(emoji, source[endKey:]...)...)
			}

			start += endEmoji
		}

		if start >= len(source) {
			break
		}

		k = bytes.Index(source[start:], emojiDelim)
	}

	return source
}

func initEmoji() {
	emojiMap := emoji.CodeMap()

	for k, v := range emojiMap {
		emojis[k] = []byte(v)

		if len(k) > emojiMaxSize {
			emojiMaxSize = len(k)
		}
	}

}
