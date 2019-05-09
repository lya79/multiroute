package multiroute

import (
	"net/http"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	handle  http.HandlerFunc
}

type MultiRouter struct {
	rootPattern     string
	routes          []route
	notFoundHandler http.HandlerFunc
}

func NewMultiRouter(rootPattern string, notFoundHandler http.HandlerFunc) *MultiRouter {
	serv := new(MultiRouter)
	serv.rootPattern = rootPattern
	serv.notFoundHandler = notFoundHandler
	return serv
}

func (m *MultiRouter) AddRoute(regular string, handler http.HandlerFunc) {
	r := route{
		pattern: regexp.MustCompile(regular),
		handle:  handler,
	}
	m.routes = append(m.routes, r)
}

func (m *MultiRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lenURL := len(r.URL.Path)

	x := len(m.rootPattern)
	y := lenURL // len=y-x
	z := lenURL // cap=z-x
	sub := []byte(r.URL.Path)[x:y:z]

	for i := range m.routes {
		if m.routes[i].pattern.Match(sub) {
			m.routes[i].handle(w, r)
			return
		}
	}

	m.notFoundHandler(w, r)
}
