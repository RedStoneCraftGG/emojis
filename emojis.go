package emojis

import (
	"fmt"
	"log"
	"regexp"

	"github.com/df-mc/dragonfly/server/player"
)

// EmojiPattern is a regular expression compiled to match emoji shortcodes formatted as :alias:.
// For example, :smile:, :heart:, or :custom_emoji:.
var EmojiPattern = regexp.MustCompile(`:[a-zA-Z0-9_]+:`)

// EmojiMap maps emoji aliases (e.g., ":smile:") to their corresponding hexadecimal string representations.
// These hexadecimal strings are used to generate Unicode characters within a designated private use area.
var EmojiMap = map[string]string{}

// EmojiHandler is a player handler designed to intercept player chat messages.
// It inspects chat text and replaces any valid emoji shortcodes with the corresponding Unicode private-use glyph.
type EmojiHandler struct {
	player.NopHandler
}

// AddEmoji registers a new emoji mapping by associating the given alias (e.g., ":heart:")
// with the specified hex code (e.g., "E001"). This enables other packages to extend the available emoji set at runtime.
func AddEmoji(alias, hex string) {
	EmojiMap[alias] = hex
}

// AddEmojis registers multiple emoji mappings at once, accepting a map where each key is an emoji alias
// (e.g., ":admin:") and each value is the hexadecimal code representing the emoji.
func AddEmojis(emojis map[string]string) {
	for alias, hex := range emojis {
		AddEmoji(alias, hex)
	}
}

// HandleChat is a chat event handler that intercepts player messages. It searches the message for all
// emoji shortcodes that match the EmojiPattern regular expression. For each match, if a corresponding
// hex code is registered in EmojiMap, the function attempts to convert it to a Unicode glyph with EmojiFormat.
// If successful, the shortcode in the message is replaced with the glyph; otherwise, the original shortcode
// remains, and any conversion errors are logged.
func (h EmojiHandler) HandleChat(ctx *player.Context, msg *string) {
	modified := EmojiPattern.ReplaceAllStringFunc(*msg, func(match string) string {
		if hex, ok := EmojiMap[match]; ok {
			if glyph, err := EmojiFormat(hex); err == nil {
				return glyph
			} else {
				log.Println("Error converting", match, ":", err)
			}
		}
		return match
	})

	*msg = modified
}

// EmojiFormat converts a custom emoji code in the format "EXXX" (hexadecimal)
// into a corresponding Unicode character in the Private Use Area (U+E000 - U+EFFF).
// Returns the emoji as a string if successful, or an error if the code is invalid or out of range.
func EmojiFormat(code string) (string, error) {
	var hex int
	if _, err := fmt.Sscanf(code, "E%X", &hex); err != nil {
		return "", err
	}
	// Private Use Area (U+E000 - U+EFFF)
	r := rune(0xE000 + hex)
	if r < 0xE000 || r > 0xEFFF {
		return "", fmt.Errorf("out of private use area")
	}
	return string(r), nil
}
