package template

import "log"

type Def struct {
	kind DefKind
	left Token
}

type DefKind int

const (
	DefInvalid DefKind = iota
	DefTagPackage
	DefListPackage
	DefHtmlPackage
	DefRecord
	DefAlias
	DefTemplate
	DefDocline
	DefDocblock
	DefImport
	DefUsing
	DefMetatable
)

func NewParser(t *Tokenizer) Parser {
	if t.source == nil {
		panic("source must not be nil")
	}
	return Parser{t, []Def{}}
}

type Parser struct {
	tokenizer *Tokenizer
	ast       []Def
}

func (p *Parser) Parse(start int) ([]Def, error) {
	var next int
	var err error

	next, err = p.docPackage(start)
	if err == ErrNoMatch {
		log.Println("no package declaration was found")
	} else if err == ErrInvalid {
		log.Println("invalid package statement")
	}

	next, err = p.repeat(next, p.docImport)
	if err == ErrInvalid {
		log.Println("invalid import statement")
	}

	next, err = p.repeat(next, p.docUsing)
	if err == ErrInvalid {
		log.Println("invalid using statement")
	}

	_, err = p.repeat(next, p.docDef)
	if err == ErrInvalid {
		log.Println("invalid definition")
	}

	ast := p.ast
	p.ast = nil

	// TODO: return next
	return ast, nil
}

func (p *Parser) defTypeOrTempl(start int) (next int, err error) {
	var ast Def

	if next, err = p.defType(start); err != nil {

		if err == EOF {
			return start, EOF
		}

		if next, err = p.defTempl(start); err != nil {
			log.Println(err)
			return start, ErrInvalid
		}
	}

	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) decl(start int, kind TokenKind) (Token, int, bool) {
	var ident Token
	next := start

	if token, n, ok := p.match(next, TokenIdent); ok {
		ident = token
		next = n
	} else {
		return token, start, false
	}

	if token, n, ok := p.match(next, TokenColon); ok {
		next = n
	} else {
		return token, start, false
	}

	if _, n, ok := p.match(next, kind); ok {
		next = n
	}
	return ident, next, true
}

type matchhandler func(int) (Token, int, bool)

type matcher func(matchhandler) matchhandler

func (p *Parser) match(start int, kind TokenKind, matchers ...matcher) (Token, int, bool) {
	var main matchhandler = func(s int) (Token, int, bool) {
		token, next, err := p.tokenizer.next(s)
		if err != nil {
			log.Println(err)
			return token, s, false
		}
		if token.kind != kind {
			return token, s, false
		}
		return token, next, true
	}

	var handler matchhandler = main
	for _, matcher := range matchers {
		handler = matcher(handler)
	}

	handler = p.skipBefore(TokenSpace, TokenComment, TokenEOL)(handler)

	return handler(start)
}

func (p *Parser) skipBefore(kind TokenKind, more ...TokenKind) matcher {
	return func(h matchhandler) matchhandler {
		return func(start int) (Token, int, bool) {
			kinds := make([]TokenKind, 0, len(more)+1)
			kinds = append(kinds, kind)
			kinds = append(kinds, more...)
			var beforeToken = start
			for {
				toSkip, next, err := p.tokenizer.next(beforeToken)
				if err != nil {
					log.Println(err)
					return toSkip, start, false
				}
				for _, kind := range kinds {
					if toSkip.kind == kind {
						goto skip
					}
				}
				break
			skip:
				beforeToken = next
			}
			return h(beforeToken)
		}
	}
}

type repeatFunc func(int) (int, error)

func (p *Parser) repeat(start int, fn repeatFunc) (int, error) {
	var err error
	next := start
	for {
		if next, err = fn(next); err != nil {
			break
		}
	}
	return next, err
}
