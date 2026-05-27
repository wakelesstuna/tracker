package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/wakelesstuna/tracker/internal/ical"
	"github.com/wakelesstuna/tracker/internal/pushups"
)

var tpl = template.Must(template.New("page").Parse(`
<html>
<body>
<h1>My Calendar</h1>

{{ range . }}
<div style="margin-bottom:20px">
	<h3>{{ .Total }}</h3>
	<p>{{ .Date }}</p>
</div>
{{ end }}

</body>
</html>
`))

func handler(w http.ResponseWriter, r *http.Request) {
	c := &ical.Client{}
	ups := pushups.Pushups{Ical: c}

	events, err := ups.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tpl.Execute(w, events.Items)
}

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("Running on :8080")
	http.ListenAndServe(":8080", nil)
}
