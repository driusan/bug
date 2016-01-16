package scm

type UnsupportedType string

func (a UnsupportedType) Error() string {
	return string(a)
}

