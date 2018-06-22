package atcs

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// PostAPI PostAPI
func (a Atcs) PostAPI(api string, params string) map[string]interface{} {
	p := map[string]string{}
	err := json.Unmarshal([]byte(params), &p)
	if err != nil {
		log.Println(err)
		return nil
	}
	postParam := url.Values{}
	for key, value := range p {
		postParam[key] = []string{value}
	}
	resp, err := http.PostForm(api, postParam)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Println(err)
		return nil
	}
	return result
}

// GetAPI GetAPI
func (a Atcs) GetAPI(api string) map[string]interface{} {
	resp, err := http.Get(api)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Println(err)
		return nil
	}
	return result
}

// GetBody GetBody
func (a Atcs) GetBody(url string) template.HTML {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	return template.HTML(b)
}
