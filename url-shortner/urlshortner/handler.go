package urlshortner

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

func readYAMLFile() map[string]string {
	allLinks := []map[string]string{}

	yamlFile, err := ioutil.ReadFile("links.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &allLinks)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	mappedLinks := parseLinksIntoMap(allLinks)

	return mappedLinks
}

func readJSONFile() map[string]string {
	links := []map[string]string{}

	jsonFile, err := ioutil.ReadFile("links.json")

	if err != nil {
		log.Fatalf("jsonFile.Get err #%v ", err)
	}

	err = json.Unmarshal(jsonFile, &links)

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

func YAMLHandler(fallback http.HandlerFunc) http.HandlerFunc {
	mappedLinks := readYAMLFile()

	return MapHandler(mappedLinks, fallback)
}

func JSONHandler(fallback http.Handler) http.HandlerFunc {
	mappedLinks := readJSONFile()

	return MapHandler(mappedLinks, fallback)
}
