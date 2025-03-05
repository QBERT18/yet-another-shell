package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func (t Token) String() string {
	tokenNames := map[Token]string{
		ILLEGAL:      "ILLEGAL",
		EOF:          "EOF",
		WS:           "WS",
		IDENT:        "IDENT",
		PIPE:         "PIPE",
		AMPERSAND:    "AMPERSAND",
		SEMICOLON:    "SEMICOLON",
		REDIR_IN:     "REDIR_IN",
		REDIR_OUT:    "REDIR_OUT",
		DOLLAR:       "DOLLAR",
		BACKTICK:     "BACKTICK",
		DASH:         "DASH",
		SINGLE_QUOTE: "SINGLE_QUOTE",
		DOUBLE_QUOTE: "DOUBLE_QUOTE",
		ECHO:         "ECHO",
		CD:           "CD",
		LS:           "LS",
		CAT:          "CAT",
		TOUCH:        "TOUCH",
		MKDIR:        "MKDIR",
		RM:           "RM",
		CP:           "CP",
		MV:           "MV",
		CHMOD:        "CHMOD",
		EXPORT:       "EXPORT",
		UNSET:        "UNSET",
		PWD:          "PWD",
		WHICH:        "WHICH",
	}

	if name, ok := tokenNames[t]; ok {
		return name
	}

	return fmt.Sprintf("UNKNOWN(%d)", t)
}

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	switch ch {
	case eof:
		return EOF, ""
	case '|':
		return PIPE, string(ch)
	case '&':
		return AMPERSAND, string(ch)
	case ';':
		return SEMICOLON, string(ch)
	case '<':
		return REDIR_IN, string(ch)
	case '>':
		return REDIR_OUT, string(ch)
	case '$':
		return DOLLAR, string(ch)
	case '`':
		return BACKTICK, string(ch)
	case '-':
		return DASH, string(ch)
	case '\'':
		return SINGLE_QUOTE, string(ch)
	case '"':
		return DOUBLE_QUOTE, string(ch)
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	// Check if the identifier is a keyword
	switch strings.ToUpper(buf.String()) {
	case "ECHO", "WRITE-OUTPUT":
		return ECHO, buf.String()
	case "CD", "SET-LOCATION":
		return CD, buf.String()
	case "LS", "GET-CHILDITEM":
		return LS, buf.String()
	case "CAT", "GET-CONTENT":
		return CAT, buf.String()
	case "TOUCH", "NEW-ITEM":
		return TOUCH, buf.String()
	case "MKDIR", "NEW-ITEM_DIR":
		return MKDIR, buf.String()
	case "RM", "REMOVE-ITEM":
		return RM, buf.String()
	case "CP", "COPY-ITEM":
		return CP, buf.String()
	case "MV", "MOVE-ITEM":
		return MV, buf.String()
	case "CHMOD", "SET-ACL":
		return CHMOD, buf.String()
	case "EXPORT", "SET-ENV":
		return EXPORT, buf.String()
	case "UNSET", "REMOVE-ENV":
		return UNSET, buf.String()
	case "PWD", "GET-LOCATION":
		return PWD, buf.String()
	case "WHICH", "GET-COMMAND":
		return WHICH, buf.String()
	}

	return IDENT, buf.String()
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

var eof = rune(0)
