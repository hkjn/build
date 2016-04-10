package main

import (
	"path"
	"strings"

	"github.com/kelseyhightower/envconfig"

	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type config struct {
	Addr, Key, Cert, Tasks string
	Paths                  []string
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to %s\n", r.RequestURI)
	if r.Method != http.MethodPost {
		http.Error(w, "Bad method, want POST.", http.StatusBadRequest)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request: %v\n", err)
		http.Error(w, "Well, that's embarrasing.", http.StatusInternalServerError)
		return
	}
	log.Printf("The POST was from %q: \n%q\n", r.RemoteAddr, string(b))
	parts := strings.Split(r.requestURI, "/")
	if len(parts) != 4 {
		log.Printf("Can't understand request; want /some-domain.com/some-repo/abcd; got %q\n", parts)
		http.Error(w, "Well, that's embarrasing.", http.StatusInternalServerError)
		return
	}
	repo := path.Join(parts[1], parts[2])
	log.Printf("This was request to build repo %q\n", repo)
	// TODO(hkjn): Write file to c.Tasks directory for each build
	// request that makes it here with info on repo, commit, branch.

	// TODO(hkjn): Set up separate watcher for c.Tasks directory which
	// does the docker build && docker push, then removes the task file.

	// TODO(hkjn): Set up separate watcher for c.Tasks directory which
	// notifies on Slack when build task lands / is finished.
	fmt.Fprintf(w, "hi!\n")
}

func run(c config) error {
	if c.Addr == "" {
		return errors.New("no BUILD_ADDR")
	}
	if c.Cert == "" {
		return errors.New("no BUILD_CERT")
	}
	if c.Key == "" {
		return errors.New("no BUILD_KEY")
	}
	if len(c.Paths) == 0 {
		return errors.New("no BUILD_PATHS")
	}
	log.Printf("build serving HTTPS on %s, using cert %q and key %q..\n", c.Addr, c.Cert, c.Key)
	for _, p := range c.Paths {
		http.HandleFunc(p, handler)
		log.Printf("Serving callbacks at %q..\n", p)
	}
	return http.ListenAndServeTLS(c.Addr, c.Cert, c.Key, nil)
}

func main() {
	var c config
	if err := envconfig.Process("build", &c); err != nil {
		panic(err)
	}
	panic(run(c))
}
