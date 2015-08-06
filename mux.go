//package router implements a Mux for github.com/korbjt/relay. It uses the
//github.com/julienschmidt/httprouter package to provide the underlying http
//multiplexing.
package router

import (
	"net/http"

	"github.com/korbjt/httprouter"
	"github.com/korbjt/relay"
	"golang.org/x/net/context"
)

type key int

const pkey key = 0

//PathParam returns a path paramter corresponding to key. It uses the
//httprouter.Params ByName method.
func PathParam(ctx context.Context, key string) string {
	if params, ok := ctx.Value(pkey).(httprouter.Params); ok {
		return params.ByName(key)
	}
	return ""
}

type wrap struct {
	route relay.Route
}

func (w *wrap) Handle(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := context.WithValue(context.Background(), pkey, ps)
	w.route.Route(ctx, r).ServeHTTP(rw, r)
}

//Mux layers the funtionality of a relay.Mux on top of an httprouter.Router.
type Mux struct {
	Router *httprouter.Router
}

func (m *Mux) Add(method, pattern string, route relay.Route) {
	m.Router.Handle(method, pattern, &wrap{route: route})
}

func (m *Mux) Route(ctx context.Context, r *http.Request) http.Handler {
	handler, params, _ := m.Router.Lookup(r.Method, r.URL.Path)
	if handler == nil {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}
		return http.HandlerFunc(fn)
	}
	if w, ok := handler.(*wrap); ok {
		ctx = context.WithValue(ctx, pkey, params)
		return w.route.Route(ctx, r)
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		handler.Handle(w, r, params)
	}

	return http.HandlerFunc(fn)
}

//New creates a new router with the default httprouter.Router configuration
func New() *Mux {
	return &Mux{
		Router: httprouter.New(),
	}
}
