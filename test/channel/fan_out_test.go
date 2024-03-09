package channel

func fanOut[T any](ch <-chan T, out []chan T, async bool) {
	go func() {
		defer func() {
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()

		for v := range ch {
			v := v
			for i := 0; i < len(out); i++ {
				i := i
				if async {
					go func() {
						out[i] <- v
					}()
				} else {
					out[i] <- v
				}
			}
		}
	}()
}
