package lexer

import "fmt"

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // variable names, commands

	// Operators and special characters
	PIPE         // |
	AMPERSAND    // &
	SEMICOLON    // ;
	REDIR_IN     // <
	REDIR_OUT    // >
	DOLLAR       // $
	BACKTICK     // `
	DASH         // -
	SINGLE_QUOTE // '
	DOUBLE_QUOTE // "
	ASTERISK     // *
	COMMA        // ,

	// Built-in shell commands
	ECHO
	CD
	LS
	CAT
	TOUCH
	MKDIR
	RM
	CP
	MV
	CHMOD
	EXPORT
	UNSET
	PWD
	WHICH

	// PowerShell equivalents
	WRITE_OUTPUT  // Equivalent to echo
	SET_LOCATION  // Equivalent to cd
	GET_CHILDITEM // Equivalent to ls
	GET_CONTENT   // Equivalent to cat
	NEW_ITEM      // Equivalent to touch
	NEW_ITEM_DIR  // Equivalent to mkdir
	REMOVE_ITEM   // Equivalent to rm
	COPY_ITEM     // Equivalent to cp
	MOVE_ITEM     // Equivalent to mv
	SET_ACL       // Equivalent to chmod
	SET_ENV       // Equivalent to export
	REMOVE_ENV    // Equivalent to unset
	GET_LOCATION  // Equivalent to pwd
	GET_COMMAND   // Equivalent to which
)

func String(t Token) string {
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
