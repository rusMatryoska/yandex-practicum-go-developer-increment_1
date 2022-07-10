package main

import (
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

//separate to packeges

func addNewKey(url string, j *int, ui *map[string]int, iu *map[int]string, mutex *sync.Mutex) {

	mutex.Lock()
	defer mutex.Unlock()

	currentMap := *ui
	currentMapMirror := *iu

	if currentMap[url] == 0 {
		*j = *j + 1
		currentMap[url] = *j
		currentMapMirror[*j] = url
	}
}

var idUrl, urlId, i = make(map[int]string), make(map[string]int), 1000
var mutex sync.Mutex

func commonHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		fullUrl, err := io.ReadAll(r.Body)
		url := string(fullUrl)

		//we won't add empty-string
		if err != nil || url == "" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		addNewKey(url, &i, &urlId, &idUrl, &mutex)
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(urlId[url])))

	case "GET":
		id := r.URL.Query().Get("id")

		x, err := strconv.Atoi(id)
		gettingUrl := idUrl[x]

		if reflect.TypeOf(x).Kind() != reflect.Int || err != nil {
			http.Error(w, "ID parameter must be Integer type", http.StatusBadRequest)
			return
		}

		//checking is this url exists
		if gettingUrl == "" {
			http.Error(w, "There is not url with this id", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusTemporaryRedirect)
		io.WriteString(w, gettingUrl)
	}

}

func main() {
	http.HandleFunc("/", commonHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
