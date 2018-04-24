package tester

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	MSG = iota
	ERR
	TST
)

var (
	StdOut     chan MessageCh
	LineLength = 70
)

type MessageCh struct {
	Type   int
	Format string
	Args   []interface{}
}

func PrinterWatch(msg <-chan MessageCh) {
	for m := range msg {
		switch m.Type {
		case ERR:
			msg := fmt.Sprintf(m.Format, m.Args...)
			fmt.Println(padding(msg, "fail", true))
		case MSG:
			msg := fmt.Sprintf(m.Format, m.Args...)
			fmt.Println(strings.Join(wordWrap(msg, LineLength), "\n"))
		default:
			msg := fmt.Sprintf(m.Format, m.Args...)
			fmt.Println(padding(msg, "ok", true))
		}
	}
}

func Message(format string, a ...interface{}) {
	StdOut <- MessageCh{
		Type:   MSG,
		Format: format,
		Args:   a,
	}
}

func MessageTest(format string, a ...interface{}) {
	StdOut <- MessageCh{
		Type:   TST,
		Format: format,
		Args:   a,
	}
}

func MessageError(format string, a ...interface{}) {
	StdOut <- MessageCh{
		Type:   ERR,
		Format: format,
		Args:   a,
	}
}

func padding(left, right string, wrap bool) string {
	var (
		dots    string
		strSum  = utf8.RuneCountInString(left) + utf8.RuneCountInString(right)
		dotsQty int
	)

	dotsQty = LineLength - strSum
	if strSum >= LineLength {
		if wrap {
			var res string
			ww := wordWrap(left, LineLength-utf8.RuneCountInString(right)-3)
			for _, s := range ww {
				res += s + "\n"
			}
			left = strings.Trim(res, "\n")
			dotsQty = LineLength - utf8.RuneCountInString(ww[len(ww)-1]) - utf8.RuneCountInString(right)
		} else {
			left = left[:LineLength-utf8.RuneCountInString(right)-3]
			strSum = utf8.RuneCountInString(left) + utf8.RuneCountInString(right)
			dotsQty = LineLength - strSum
		}
	}
	dots = strings.Repeat(".", dotsQty)
	return fmt.Sprintf("%s%s%s", left, dots, right)
}

func wordWrap(text string, lineWidth int) []string {
	var chunks []string
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return []string{text}
	}
	wrapped := words[0]
	spaceLeft := lineWidth - utf8.RuneCountInString(wrapped)
	for _, word := range words[1:] {
		if utf8.RuneCountInString(word)+1 > spaceLeft {
			chunks = append(chunks, wrapped)
			wrapped = word
			spaceLeft = lineWidth - utf8.RuneCountInString(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + utf8.RuneCountInString(word)
		}
	}
	chunks = append(chunks, wrapped)

	return chunks
}
