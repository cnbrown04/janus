package ui

import (
	"strings"
	"sync"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/lipgloss"
)

// Gruvbox-adjacent palette for the query editor.
var (
	sqlComment    = lipgloss.NewStyle().Foreground(lipgloss.Color("#928374")).Italic(true)
	sqlKeyword    = lipgloss.NewStyle().Foreground(lipgloss.Color("#fb4934")) // red
	sqlString     = lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26")) // bright green
	sqlNumber     = lipgloss.NewStyle().Foreground(lipgloss.Color("#d3869b")) // purple
	sqlName       = lipgloss.NewStyle().Foreground(lipgloss.Color("#83a598")) // blue
	sqlOperator   = lipgloss.NewStyle().Foreground(lipgloss.Color("#fe8019")) // orange
	sqlPunct      = lipgloss.NewStyle().Foreground(lipgloss.Color("#a89984")) // gray
	sqlDefault    = lipgloss.NewStyle()
)

var (
	sqlLexerOnce sync.Once
	sqlLexer     chroma.Lexer
)

func initSQLLexer() {
	sqlLexerOnce.Do(func() {
		l := lexers.Get("sql")
		if l == nil {
			l = lexers.Fallback
		}
		sqlLexer = chroma.Coalesce(l)
	})
}

func styleForSQLToken(t chroma.TokenType) lipgloss.Style {
	switch t {
	case chroma.Keyword, chroma.KeywordConstant, chroma.KeywordDeclaration,
		chroma.KeywordPseudo, chroma.KeywordReserved, chroma.KeywordType:
		return sqlKeyword
	case chroma.String, chroma.StringChar, chroma.StringDouble, chroma.StringBacktick,
		chroma.StringHeredoc, chroma.StringSingle, chroma.StringInterpol, chroma.StringAffix:
		return sqlString
	case chroma.Number, chroma.NumberInteger, chroma.NumberIntegerLong, chroma.NumberFloat, chroma.NumberHex,
		chroma.NumberOct, chroma.NumberBin:
		return sqlNumber
	case chroma.Name, chroma.NameAttribute, chroma.NameBuiltin, chroma.NameClass, chroma.NameDecorator,
		chroma.NameEntity, chroma.NameException, chroma.NameFunction, chroma.NameProperty, chroma.NameLabel,
		chroma.NameNamespace, chroma.NameOther, chroma.NameTag, chroma.NameVariable,
		chroma.NameConstant, chroma.NameBuiltinPseudo, chroma.NameFunctionMagic:
		return sqlName
	case chroma.Comment, chroma.CommentSingle, chroma.CommentMultiline, chroma.CommentPreproc,
		chroma.CommentSpecial, chroma.CommentHashbang:
		return sqlComment
	case chroma.Operator, chroma.OperatorWord:
		return sqlOperator
	case chroma.Punctuation:
		return sqlPunct
	default:
		return sqlDefault
	}
}

// highlightSQLLine returns ANSI-styled text for one logical line (no newlines).
func highlightSQLLine(line string) string {
	initSQLLexer()
	it, err := sqlLexer.Tokenise(nil, line)
	if err != nil {
		return line
	}
	var b strings.Builder
	for t := it(); t != chroma.EOF; t = it() {
		if t.Value == "" {
			continue
		}
		b.WriteString(styleForSQLToken(t.Type).Render(t.Value))
	}
	return b.String()
}
