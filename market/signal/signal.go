package signal

type Signal int

const (
	BuySignal Signal = iota + 1
	SellSignal
)
