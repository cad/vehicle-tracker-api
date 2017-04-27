package event

var eventKinds = map[string]Kind{}

var stream = make(chan *Event)
var stop = make(chan bool)

type Kind struct {
	Name string
	handlers *[]*func(e *Event)
	stream *chan *Event
}

type Event struct {
	kind *Kind
	Payload interface{}
}

func (k *Kind) Emit(payload interface{}) {
	go func(){
		event := Event{
			kind: k,
			Payload: payload,
		}
		*k.stream <- &event
	}()
}

func (k *Kind) Register(handler *func(event *Event)){
	*k.handlers = append(*k.handlers, handler)
}

func (k *Kind) UnRegister(handler *func(event *Event)){
	found := -1
	for i, v := range *k.handlers {
		if v == handler {
			found = i
			break
		}
	}
	(*k.handlers)[found] = (*k.handlers)[len(*k.handlers)-1]
	*k.handlers = (*k.handlers)[:len(*k.handlers)-1]
}

func MakeKind(name string) *Kind {
	if _, ok := eventKinds[name]; !ok {
		eventKinds[name] = Kind{
			Name: name,
			stream: &stream,
			handlers: &[]*func(e *Event){},
		}
	}
	kind := eventKinds[name]
	return &kind
}

func Run() {
	go func(){

		for {
			var e *Event
			select {
			case e = <-stream:
				// Handle
				for _, handler := range *e.kind.handlers {
					(*handler)(e)
				}

			case <- stop:
				// Quit
				return

			}
		}
	}()
}

func Shutdown() {
	stop <- true
}
