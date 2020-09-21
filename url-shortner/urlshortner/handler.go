package urlshortner

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

func readYAMLFile(path string) map[string]string {
	allLinks := []map[string]string{}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &allLinks)
	if err != nil {
		log.Fatalf("Failed to parse yaml file: %v", err)
	}

	mappedLinks := parseLinksIntoMap(allLinks)
	return mappedLinks
}

func readJSONFile(path string) map[string]string {
	links := []map[string]string{}

	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("jsonFile.Get err #%v ", err)
	}

	err = json.Unmarshal(jsonFile, &links)
	if err != nil {
		log.Fatalf("Failed to parse file data into json format: #%v ", err)
	}

	mappedLinks := parseLinksIntoMap(links)
	return mappedLinks
}

func parseLinksIntoMap(allLinks []map[string]string) map[string]string {
	links := make(map[string]string)

	for _, entry := range allLinks {
		key := entry["path"]
		links[key] = entry["url"]
	}
	return links
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(fallback http.HandlerFunc, path string) http.HandlerFunc {
	mappedLinks := readYAMLFile(path)
	return MapHandler(mappedLinks, fallback)
}

func JSONHandler(fallback http.Handler, path string) http.HandlerFunc {
	mappedLinks := readJSONFile(path)
	return MapHandler(mappedLinks, fallback)
}
