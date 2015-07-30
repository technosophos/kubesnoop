package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const tpl = `<!doctype html>
<html>
<head>
<title>KubeSnoop</title>
</head>
<body>

<h2>Environment Variables</h2>
<pre>{{.Env}}</pre>

</body>
</html>
`

var renderer *template.Template

func main() {
	renderer = template.Must(template.New("html").Parse(tpl))

	http.HandleFunc("/", snoop)

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Println(addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func snoop(w http.ResponseWriter, r *http.Request) {

	tvars := map[string]string{
		"Env": "",
	}

	out, _ := exec.Command("env").CombinedOutput()
	tvars["Env"] = string(out)

	renderer.Execute(w, tvars)
}
