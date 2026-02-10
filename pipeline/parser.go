package pipeline

import "context"

type Parser struct {
	prefix string
}

func NewParser() *Parser {
	return &Parser{
		prefix: "parsed - ",
	}
}

func (p *Parser) Parse(ctx context.Context, in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case s, ok := <-in:
				if !ok {
					return
				}

				parsed := p.parse(s)
				select {
				case <-ctx.Done():
					return
				case out <- parsed:
				}
			}
		}
	}()

	return out
}

func (p *Parser) parse(s string) string {
	return p.prefix + s
}
