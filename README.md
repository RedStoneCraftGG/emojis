# emojis

A minimal and extensible emoji glyph system for [dragonfly-mc](https://github.com/df-mc/dragonfly).

---

## Features

- Replace emoji shortcodes like `:smile:` or `:heart:` in player chat messages with custom Unicode Private Use Area (PUA) glyphs.
- Simple runtime registration of new emojis using Go code.
- Plug-and-play player handler integration with Dragonfly.

---

## Getting Started

### 1. Installing

Just import the package anywhere you need emoji support:

```go
import emojis "github.com/redstonecraftgg/emojis"
```

### 2. Register Emojis

Register emoji shortcodes and their corresponding PUA hex codes:

```go
emojis.AddEmojis(map[string]string{
    ":heart:": "E001",
    ":smile:": "E002",
    ":laugh:": "E003",
})
```

You can also add a single emoji at runtime:

```go
emojis.AddEmoji(":star:", "E010")
```

Only codes in the range `E000`-`EFFF` (PUA) are supported.

---

### 3. Enable Emoji Chat Handler 

Register the handler for each player (example using Dragonfly's `srv.Accept()`):

```go
for p := range srv.Accept() {
    p.Handle(emojis.EmojiHandler{})
}
```

---

## How it works

- Players type messages in chat using shortcodes (e.g., `I :heart: Dragonfly!`).
- The `EmojiHandler` detects and replaces matching shortcodes based on a regular expression.
- The corresponding shortcode is mapped to a PUA code (e.g., `E001`), and that Unicode glyph is inserted.
- The actual display of emoji glyphs relies on the client/resource pack supporting those PUA codepoints.

---

## Limitations / Notes

- **Resource pack support is required:** Emoji glyphs in the Private Use Area (PUA — U+E000 ... U+EFFF) will NOT be visible unless you provide a Minecraft resource pack that maps your emoji to these codepoints.
- Shortcodes not registered remain unchanged in chat.
- Adding, updating, or removing emoji mappings is safe at runtime.

---

## Example

```go
emojis.AddEmojis(map[string]string{
    ":yes:": "E123",
    ":no:":  "E124",
})

// When a player types:
Hello :yes: :smile:

// ...they see (if resourcepack installed):
Hello <custom glyph for yes> <custom glyph for smile>
```
