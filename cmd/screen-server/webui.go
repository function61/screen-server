package main

import (
	"context"
	"html/template"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/function61/gokit/io/bidipipe"
	"github.com/function61/gokit/net/http/httputils"
	"github.com/function61/holepunch-server/pkg/wsconnadapter"
	"github.com/gorilla/websocket"
)

func AnyOriginIsOk(r *http.Request) bool { return true }

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     AnyOriginIsOk,
}

var ui, _ = template.New("_").Parse(`<!doctype html>

<html>
<head>
	<title>Screen server</title>
</head>
<body>

{{range .}}
<h3>{{.Opts.Description}}</h3>

<a href="/static/vnc/?path={{.VncWebsocketPath}}&autoconnect=1">
	<img src="/api/screen/{{.Id}}/screenshot" alt="Screenshot from screen" style="width: 30%;" />
</a>

<hr />
{{end}}

</body>
</html>
`)

func newServerHandler(screens []*Screen, logger *slog.Logger) http.Handler {
	routes := http.NewServeMux()

	// serves VNC client etc.
	routes.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./www"))))

	// websocket-to-VNC proxy
	routes.HandleFunc("/api/screen/{id}/ws", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			httputils.Error(w, http.StatusBadRequest)
			return
		}

		screen := screenByID(id, screens)
		if screen == nil {
			http.NotFound(w, r)
			return
		}

		// internally also checks for 'Upgrade: websocket' header
		clientWs, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("failure upgrading", "err", err)
			// "If the upgrade fails, then Upgrade replies to the client with an HTTP error response. "
			return
		}

		// adapt to regular net.Conn
		client := wsconnadapter.New(clientWs)
		defer client.Close()

		vncServer, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(screen.Opts.vncPort)))
		if err != nil {
			logger.Error("proxy dial to screen", "err", err)
			// can't reply error to client, because we're raw sockets now
			return
		}

		if err := bidipipe.Pipe(bidipipe.WithName("client", client), bidipipe.WithName("server", vncServer)); err != nil {
			logger.Error("bidipipe", "err", err)
			return
		}
	})

	routes.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := ui.Execute(w, screens)

		logIfError("frontpage", err)
	})

	routes.HandleFunc("/api/screen/{id}/osd/notify", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		screen := screenByID(id, screens)
		if screen == nil {
			http.NotFound(w, r)
			return
		}

		msg := r.FormValue("msg")
		if msg == "" {
			http.Error(w, "empty message", http.StatusBadRequest)
			return
		}

		go showOsdMessage(context.Background(), screen, msg)
	})

	routes.HandleFunc("/api/screen/{id}/screenshot", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		screen := screenByID(id, screens)
		if screen == nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		logIfError("screenshot", screen.Screenshot(w))
	})

	return routes
}

func runServer(ctx context.Context, handler http.Handler) error {
	srv := &http.Server{
		Addr:              ":80",
		Handler:           handler,
		ReadHeaderTimeout: httputils.DefaultReadHeaderTimeout,
	}

	return httputils.CancelableServer(ctx, srv, srv.ListenAndServe)
}

func logIfError(origin string, err error) {
	if err != nil {
		log.Printf("%s: %v", origin, err)
	}
}

func screenByID(id int, screens []*Screen) *Screen {
	for _, screen := range screens {
		if screen.ID == id {
			return screen
		}
	}

	return nil
}
