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
	ndata, _ := json.Marshal(&response)
	ndata = append(ndata, '\n')
	writer.Write([]byte(ndata))
}
func main() {
	http.HandleFunc("/shorten", UrlShortnerHandler)
	http.ListenAndServe(":8080", nil)
}
