package webapp

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"ratatoskr/internal/logger"
	"time"
)

//go:embed static
var content embed.FS

func NewServer(logger *logger.Logger) (http.Handler, error) {
	mux := http.NewServeMux()
	t, err := loadTemplate()
	if err != nil {
		return nil, logger.Error(err.Error())
	}
	addRoutes(mux, logger, t)

	var handler http.Handler = mux
	handler = LoggerMiddleware(logger, handler)
	return handler, nil
}

func addRoutes(
	mux *http.ServeMux,
	logger *logger.Logger,
	t *template.Template,
) {
	mux.Handle("/static/", http.FileServer(http.FS(content)))
	mux.HandleFunc("/", handleHome(logger, t))
}

func handleHome(logger *logger.Logger, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := t.ExecuteTemplate(w, "/", nil)
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func loadTemplate() (*template.Template, error) {
	tmpl := template.New("main")
	t, err := tmpl.ParseFS(content, "static/view.html")
	return t, err
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func LoggerMiddleware(l *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		l.Info(fmt.Sprintf("Started %s %s", r.Method, r.URL.Path))
		responseRecorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(responseRecorder, r)
		l.Info(fmt.Sprintf(
			"Completed in %v %s %s %d",
			time.Since(startTime),
			r.Method,
			r.URL.Path,
			responseRecorder.statusCode,
		))
	})
}
