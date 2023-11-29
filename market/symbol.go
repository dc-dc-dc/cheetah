package market

type Symbol string

func (s Symbol) String() string {
	return string(s)
}
