package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/function61/gokit/net/http/httputils"
	"github.com/gorilla/mux"
)

var ui, _ = template.New("_").Parse(`<!doctype html>

<html>
<head>
	<title>Screen server</title>
</head>
<body>

{{range .}}
<h3>{{.Opts.Description}}</h3>

<img src="/api/screen/{{.Id}}/screenshot" alt="Screenshot from screen" style="width: 30%;" />

<hr />
{{end}}

</body>
</html>
`)

func newServerHandler(screens []*Screen) http.Handler {
	routes := mux.NewRouter()

	routes.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := ui.Execute(w, screens)

		logIfError("frontpage", err)
	})

	routes.HandleFunc("/api/screen/{id}/osd/notify", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		screen := screenById(id, screens)
		if screen == nil {
			http.NotFound(w, r)
			return
		}

		msg := r.FormValue("msg")
		if msg == "" {
			http.Error(w, "empty message", http.StatusBadRequest)
			return
		}

		go showOsdMessage(context.Background(), screen, string(msg))
	})

	routes.HandleFunc("/api/screen/{id}/screenshot", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		screen := screenById(id, screens)
		if screen == nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		logIfError("screenshot", screen.Screenshot(w))
	})

	return routes
}

func runServer(ctx context.Context, handler http.Handler, logger *log.Logger) error {
	srv := &http.Server{
		Addr:    ":80",
		Handler: handler,
	}

	return httputils.CancelableServer(ctx, srv, func() error { return srv.ListenAndServe() })
}

func logIfError(origin string, err error) {
	if err != nil {
		log.Printf("%s: %v", origin, err)
	}
}

func screenById(id int, screens []*Screen) *Screen {
	for _, screen := range screens {
		if screen.Id == id {
			return screen
		}
	}

	return nil
}
