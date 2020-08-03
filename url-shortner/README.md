# url-shortner

A url shortner application. The app uses key value pairs in maps or JSON/YAML files to match shortened URLs to the full URL.

## Install

```bash
go get github.com/supercede/go-exercises/url-shortner
cd $GOPATH/src/github.com/supercede/go-exercises/url-shortner
```

## Build & Run

- Create a file named `links.json` or `links.yaml` in the root folder containing your short paths and corresponding full URL in the format:

YAML:

```yaml
- path: '/googl'
  url: 'https://google.com'
```

JSON:

```json
[
  {
    "path": "/googl",
    "url": "https://google.com"
  }
]
```

- The app is set up to use json files by default and then falls back to YAML files. To change this, add a format flag to your run command and choose the file type like so:

  ```bash
  go run main.go --format yaml
  ```

- `http://localhost:8080/path` redirects to your full URL if it exists in your file.
