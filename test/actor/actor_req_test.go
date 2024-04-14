package actor

type MethodRequest int

const (
	Incr MethodRequest = iota
	Decr
)

type Service struct {
	queue chan MethodRequest
	v     int
}

func NewService(buffer int) *Service {
	s := &Service{
		queue: make(chan MethodRequest, buffer),
	}
	go s.schedule()
	return s
}

func (s *Service) schedule() {
	for r := range s.queue {
		if r == Incr {
			s.v++
		} else {
			s.v--
		}
	}
}

func (s *Service) Incr() {
	s.queue <- Incr
}

func (s *Service) Decr() {
	s.queue <- Decr
}
