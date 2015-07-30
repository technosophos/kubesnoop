package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const tpl = `<!doctype html>
<html>
<head>
<title>KubeSnoop</title>
</head>
<body>
<h1>KubeSnoop</h1>
<p>This utility shows you what Kubernetes is exposing to your container/pod.</p>

<h2>Basics</h2>
<ul>
	<li>CPU Count: {{.NumCPU}}</li>
</ul>

<h2>Service Endpoint</h2>
<pre>{{.Service}}</pre>

<h2>Environment Variables</h2>
<pre>{{.Env}}</pre>

<h2>Mount Points</h2>
<pre>{{.Mount}}</pre>

<h2>Secrets (Guessed)</h2>
<pre>{{.Secrets}}</pre>

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
		"Env":     vars(),
		"Mount":   mounts(),
		"NumCPU":  strconv.Itoa(runtime.NumCPU()),
		"Secrets": secrets(),
		"Service": service(),
	}

	renderer.Execute(w, tvars)
}

func vars() string {
	out, err := exec.Command("env").CombinedOutput()
	if err != nil {
		return string(out)
	}
	var lines sort.StringSlice
	lines = strings.Split(string(out), "\n")
	lines.Sort()

	return strings.Join(lines, "\n")
}

func mounts() string {
	out, err := exec.Command("mount").CombinedOutput()
	if err != nil {
		return string(out)
	}
	var lines sort.StringSlice
	lines = strings.Split(string(out), "\n")
	lines.Sort()

	return strings.Join(lines, "\n")
}

func secrets() string {
	return rsecrets("/var/run/secrets")
}

func rsecrets(dir string) string {
	f, err := os.Open(dir)
	if err != nil {
		return err.Error()
	}

	infos, err := f.Readdir(0)
	if err != nil {
		return err.Error()
	}

	var b bytes.Buffer
	for _, info := range infos {
		if info.IsDir() {
			b.WriteString(fmt.Sprintf("Dir: %s\n", path.Join(dir, info.Name())))
			b.WriteString(rsecrets(path.Join(dir, info.Name())))
		} else {
			b.Write([]byte(info.Name() + ": "))
			d, err := ioutil.ReadFile(path.Join(dir, info.Name()))
			if err != nil {
				return b.String() + "\n"
			}
			b.Write(d)
			b.Write([]byte("\n"))
		}
	}
	return b.String()
}

func service() string {
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")
	url := fmt.Sprintf("https://%s:%s/api", host, port)

	r, err := http.Get(url)
	if err != nil {
		return err.Error()
	}
	defer r.Body.Close()

	out, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
