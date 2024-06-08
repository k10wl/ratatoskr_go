package webapp

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed static
var content embed.FS

func Start() error {
	t, err := loadTemplate()
	if err != nil {
		return err
	}
	http.Handle("/static/", http.FileServer(http.FS(content)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := t.ExecuteTemplate(w, "/", nil)
		fmt.Printf("err: %v\n", err)
	})
	fmt.Println("staring server")
	return http.ListenAndServe(":8080", nil)
}

func loadTemplate() (*template.Template, error) {
	tmpl := template.New("main")
	t, err := tmpl.ParseFS(content, "static/view.html")
	return t, err
}
