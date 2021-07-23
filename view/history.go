package view

type History struct {
	stack       []func()
	HistorySize int
}

func (h *History) push(back func()) {
	h.stack = append(h.stack, back)
	if len(h.stack) > h.HistorySize {
		h.stack = h.stack[1:]
	}
}

func (h *History) pop() {
	if len(h.stack) > 1 {
		last := h.stack[len(h.stack)-2]
		last()
		h.stack = h.stack[:len(h.stack)-2]
	}
}
