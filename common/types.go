package common

type Size struct {
	Width  int
	Height int
}

func (s *Size) Add(size2 Size) {
	s.Width += size2.Width
	s.Height += size2.Height
}
