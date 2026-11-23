package lexer

type BaseLexer struct {
	Sp  int
	Len int
	Src []rune
}

func New(src string) BaseLexer {
	return BaseLexer{
		Sp:  0,
		Len: len([]rune(src)),
		Src: []rune(src),
	}
}

func (l *BaseLexer) Peek() rune {
	return l.Src[l.Sp]
}
func (l *BaseLexer) Adv() {
	l.Sp++
}
func (l *BaseLexer) IsEnd() bool {
	return l.Sp >= l.Len
}
