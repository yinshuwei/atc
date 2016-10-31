package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/yinshuwei/render"
)

// AtcConfig AtcConfig
type AtcConfig struct {
	WebPath       string
	Port          string
	Envs          map[string]string
	Extensions    []string
	ExtensionsMap map[string]bool
	IsDev         bool
}

// AtcStatic Static
type AtcStatic struct {
	// Dir is the directory to serve static files from
	Dir http.FileSystem
	// Prefix is the optional prefix used to serve the static directory content
	Prefix string
	// IndexFile defines which file to serve as index if it exists.
	IndexFile string
}

var (
	configPath = ""
	config     = &AtcConfig{}
	s          *AtcStatic
	htmlRender *render.Render
)

func main() {
	configPath = *(flag.String("C", "./config/atc.json", "config path"))
	err := readConfig()
	if err != nil {
		log.Println(err)
		return
	}

	s = &AtcStatic{
		Dir:       http.Dir(config.WebPath),
		Prefix:    "",
		IndexFile: "index.html",
	}
	htmlRender = render.New(render.Options{
		Directory:     config.WebPath,
		IsDevelopment: config.IsDev,
		Extensions:    []string{".html"},
		Funcs: []template.FuncMap{
			template.FuncMap{"post": post},
			template.FuncMap{"get": get},
		},
	})
	http.HandleFunc("/", renderHTML)
	http.HandleFunc("/_help_/", help)
	http.HandleFunc("/_clear_/", clear)
	log.Println(config.Port)
	log.Println(http.ListenAndServe(config.Port, nil))
}

func help(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(w, err)
		}
	}()
	fmt.Fprintln(w, "HELP")
}

func clear(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(w, err)
		}
	}()
	readConfig()
	htmlRender.DoCompileTemplates()
	fmt.Fprintln(w, "清理成功^_^")
}

func readConfig() error {
	f, err := os.Open(configPath)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	f.Close()

	err = json.Unmarshal(b, config)
	if err != nil {
		return err
	}
	return nil
}

func renderHTML(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(w, err)
		}
	}()

	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "必须使用GET和HEAD方式的请求", 500)
		return
	}
	file := r.URL.Path
	// if we have a prefix, filter requests by stripping the prefix
	if s.Prefix != "" {
		if !strings.HasPrefix(file, s.Prefix) {
			http.Error(w, "错误", 500)
			return
		}
		file = file[len(s.Prefix):]
		if file != "" && file[0] != '/' {
			http.Error(w, "错误", 500)
			return
		}
	}
	f, err := s.Dir.Open(file)
	if err != nil {
		http.Error(w, "不存在", 400)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(w, "不存在", 400)
		return
	}

	// try to serve index file
	if fi.IsDir() {
		// redirect if missing trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusFound)
			return
		}

		file = path.Join(file, s.IndexFile)
		f, err = s.Dir.Open(file)
		if err != nil {
			return
		}
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			return
		}
	}

	if filepath.Ext(file) == ".html" {
		serveContentWithTmpl(w, r, file, fi.ModTime(), f)
	} else {
		http.ServeContent(w, r, file, fi.ModTime(), f)
	}
}

func serveContentWithTmpl(w http.ResponseWriter, r *http.Request, name string, modtime time.Time, content io.ReadSeeker) {
	err := htmlRender.HTML(w, http.StatusOK, templateName(name), Map{
		"envs": config.Envs,
	})
	if err != nil {
		log.Println(err)
	}
}

func templateName(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[1:i]
		}
	}
	return ""
}

// Object Object
type Object interface{}

// Array Array
type Array []Object

// Map Map
type Map map[string]Object

func post(api string, params string) (Map, error) {
	p := map[string]string{}
	err := json.Unmarshal([]byte(params), &p)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	postParam := url.Values{}
	for key, value := range p {
		postParam[key] = []string{value}
	}
	resp, err := http.PostForm(api, postParam)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result := Map{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

func get(api string) (Map, error) {
	resp, err := http.Get(api)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	result := Map{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}
