package webapp

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"ratatoskr/internal/config"
	"ratatoskr/internal/db"
	"ratatoskr/internal/logger"
	"ratatoskr/internal/models"
	"time"
)

//go:embed static
var content embed.FS

func NewServer(c *config.WepAppConfig, db db.DB, logger *logger.Logger) (http.Handler, error) {
	mux := http.NewServeMux()
	t, err := loadTemplate()
	if err != nil {
		return nil, logger.Error(err.Error())
	}
	addRoutes(mux, c, db, logger, t)

	var handler http.Handler = mux
	handler = LoggerMiddleware(logger, handler)
	return handler, nil
}

func addRoutes(
	mux *http.ServeMux,
	config *config.WepAppConfig,
	db db.DB,
	logger *logger.Logger,
	template *template.Template,
) {
	mux.Handle("/static/", http.FileServer(http.FS(content)))
	mux.HandleFunc("/", tokenOnly(config, logger, handleHome(config, db, logger, template)))
	mux.Handle("/ping", ping())
}

func handleHome(
	config *config.WepAppConfig,
	db db.DB,
	logger *logger.Logger,
	template *template.Template,
) http.HandlerFunc {
	type data struct {
		Version string
		Groups  []models.Group
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		group, err := db.GetAllGroupsWithTags(ctx)
		if err != nil {
			logger.Error(err.Error())
			fmt.Fprintf(w, "")
			return
		}
		err = template.ExecuteTemplate(w, "webapp", data{
			Version: config.Version,
			Groups:  *group})
		if err != nil {
			logger.Error(err.Error())
			fmt.Fprintf(w, "")
			return
		}
	}
}

func ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}
}

func tokenOnly(c *config.WepAppConfig, l *logger.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != fmt.Sprintf("/%s", c.Token) {
			l.Error("requested server, but not bot")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next(w, r)
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
