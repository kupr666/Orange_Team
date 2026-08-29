package habbits_http_transport

type HabbitsHTTPHandler struct {
	habbitsService HabbitsService
}

type HabbitsService interface {
	
}

func NewHabbitsHTTPHandler(
	habbitsService HabbitsService,
) *HabbitsHTTPHandler {
	return &HabbitsHTTPHandler{
		habbitsService: habbitsService,
	}
}
