package hiapns

type Counter chan uint32

func NewCounter() Counter {
	cnt := make(Counter, 4)
	go func() {
		i := uint32(0)
		for {
			i++
			cnt <- i
		}
	}()
	return cnt
}

func (c Counter) Next() uint32 {
	return <-c
}
