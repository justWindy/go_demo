package channel

import "reflect"

func fanInReflect[T any](chans ...<-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv, // case语句的方向是接受
				Chan: reflect.ValueOf(c),
			})
		}

		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok {
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface().(T)
		}
	}()

	return out
}

func mergeTwo[T any](a, b <-chan T) <-chan T {
	c := make(chan T)

	go func() {
		defer close(c)

		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				c <- v
			}
		}
	}()

	return c
}

func fanInRec[T any](chans ...<-chan T) <-chan T {
	switch len(chans) {
	case 0:
		c := make(chan T)
		close(c)
		return c
	case 1:
		return chans[0]
	case 2:
		return mergeTwo(chans[0], chans[1])
	default:
		m := len(chans) / 2
		return mergeTwo(fanInRec(chans[:m]...), fanInRec(chans[m:]...))
	}
}
