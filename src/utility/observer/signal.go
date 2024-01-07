package observer

type TraceId struct {
	id int
}

var createTraceId = func() func() TraceId {
	var counter = -1
	return func() TraceId {
		counter += 1
		return TraceId{
			id: counter,
		}
	}
}()

func New() *Signal {
	return &Signal{
		observers:    nil,
		eventHistory: nil,
		handler:      nil,
	}
}

type Signal struct {
	observers    []Observer
	eventHistory []TraceId
	handler      func(id TraceId)
}

func (s *Signal) notify(traceId TraceId) {
	// when already have event with TraceId in eventHistory return
	for _, historic := range s.eventHistory {
		if historic.id == traceId.id {
			return
		}
	}

	s.eventHistory = append(s.eventHistory, traceId)

	if s.handler != nil {
		s.handler(traceId)
	}

	for _, o := range s.observers {
		o.notify(traceId)
	}

	// todo clean up history
}

func (s *Signal) Notify(id *TraceId) {
	if id != nil {
		s.notify(*id)
		return
	}
	s.notify(createTraceId())
}

func (s *Signal) SetHandler(fn func(id TraceId)) {
	s.handler = fn
}

func (s *Signal) Register(observer Observer) {
	s.observers = append(s.observers, observer)
}
