package base

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

var (
	configPath    = ""
	Config        = &AtcConfig{}
	RefCache      = map[string]*string{}
	RefCacheMutex sync.Mutex
)

// AtcConfig AtcConfig
type AtcConfig struct {
	WebPath       string
	Port          string
	Envs          map[string]string
	Extensions    []string
	ExtensionsMap map[string]bool
	IsDev         bool
	Page404       string
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

func ReadConfig() error {
	RefCacheMutex.Lock()
	defer RefCacheMutex.Unlock()
	RefCache = map[string]*string{}

	f, err := os.Open(configPath)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	f.Close()

	err = json.Unmarshal(b, Config)
	if err != nil {
		return err
	}
	return nil
}

func ConfigPath() {
	configPath = *(flag.String("C", "./config/atc.json", "config path"))
}
