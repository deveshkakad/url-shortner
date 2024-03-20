package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	urlMapping = make(map[string]string)
)

type urlRequest struct {
	Url string `json:"url"`
}
type ShortnedUrlResponse struct {
	Shortner_Url string `json:"shortned_url"`
}

func simpleHash(s string) uint32 {
	var hash uint32 = 0
	for _, c := range s {
		hash = 31*hash + uint32(c)
	}
	return hash
}

func shortenURL(url string) string {
	hashed := simpleHash(url)
	return fmt.Sprintf("%x", hashed)[:6]
}

func UrlShortnerHandler(writer http.ResponseWriter, request *http.Request) {
	body, _ := io.ReadAll(request.Body)
	req := &urlRequest{}
	err := json.Unmarshal([]byte(body), &req)
	if err != nil {
		log.Fatalf("Unmarshal failed,err:%v", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer request.Body.Close()
	originalURL := req.Url
	shortenedURL, ok := urlMapping[originalURL]
	if ok {
		response := &ShortnedUrlResponse{
			Shortner_Url: shortenedURL,
		}
		fmt.Println("Returning the existing shortned url for above domain")
		ndata, _ := json.Marshal(&response)
		ndata = append(ndata, '\n')
		writer.Write([]byte(ndata))
		return
	}
	shortenedURL = shortenURL(originalURL)
	urlMapping[shortenedURL] = originalURL
	response := &ShortnedUrlResponse{
		Shortner_Url: shortenedURL,
	}
	fmt.Println("Adding shortned url in the map and returning the result")
	ndata, _ := json.Marshal(&response)
	ndata = append(ndata, '\n')
	writer.Write([]byte(ndata))
}
func RedirectHandler(writer http.ResponseWriter, request *http.Request) {
	shortenedURL := request.URL.Path[1:]
	fmt.Printf("shortenedurl received:%v\n", shortenedURL)
	originalURL, ok := urlMapping[shortenedURL]
	fmt.Printf("urlMapping:%v", urlMapping)
	fmt.Printf("shortenedurl received:%v\n", originalURL)
	if !ok {
		http.NotFound(writer, request)
		return
	}
	http.Redirect(writer, request, originalURL, http.StatusFound)
}
func main() {
	http.HandleFunc("/shorten", UrlShortnerHandler)
	http.HandleFunc("/", RedirectHandler)
	http.ListenAndServe(":8080", nil)
}
