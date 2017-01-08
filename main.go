package main

import (
	"atc/base"
	"atc/funcs/atcs"
	"atc/funcs/ints"
	"atc/funcs/strs"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"git.oschina.net/yinshuwei/render"
)

var (
	htmlRender *render.Render
	s          *base.AtcStatic
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	base.ConfigPath()
	err := base.ReadConfig()
	if err != nil {
		log.Println(err)
		return
	}

	s = &base.AtcStatic{
		Dir:       http.Dir(base.Config.WebPath),
		Prefix:    "",
		IndexFile: "index.html",
	}

	htmlRender = render.New(render.Options{
		Directory:     base.Config.WebPath,
		IsDevelopment: base.Config.IsDev,
		Extensions:    []string{".html"},
	})
	http.HandleFunc("/", renderHTML)
	http.HandleFunc("/_help_/", help)
	http.HandleFunc("/_clear_/", clear)
	log.Println(base.Config.Port)
	log.Println(http.ListenAndServe(base.Config.Port, nil))
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
	base.ReadConfig()
	htmlRender.DoCompileTemplates()
	fmt.Fprintln(w, "清理成功^_^")
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
		serveContentWithTmpl(w, r, base.Config.Page404, http.StatusNotFound)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		serveContentWithTmpl(w, r, base.Config.Page404, http.StatusNotFound)
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
		serveContentWithTmpl(w, r, file, http.StatusOK)
	} else {
		http.ServeContent(w, r, file, fi.ModTime(), f)
	}
}

func serveContentWithTmpl(w http.ResponseWriter, r *http.Request, name string, status int) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	params := map[string]string{}
	for key, vals := range r.Form {
		if len(vals) > 0 {
			params[key] = vals[0]
		}
	}
	err = htmlRender.HTML(w, status, templateName(name), map[string]interface{}{
		"envs":   base.Config.Envs,
		"params": params,
		"ints":   ints.Ints{},
		"strs":   strs.Strs{},
		"atc":    atcs.Atcs{},
	})
	if err != nil {
		log.Println(err)
	}
}

func templateName(path string) string {
	log.Println(path)
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[1:i]
		}
	}
	return ""
}
