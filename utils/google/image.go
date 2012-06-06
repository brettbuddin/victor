package google

import (
    "net/http"
    "net/url"
    "io/ioutil"
    "math/rand"

    "victor/json"
)

func ImageSearch(term string) (string, error) {
    search, err := url.Parse("http://ajax.googleapis.com/ajax/services/search/images")

    if err != nil {
        return "", err
    }
    
    q := search.Query()
    q.Add("v", "1.0")
    q.Add("rsz", "8")
    q.Add("q", term)
    q.Add("safe", "active")
    search.RawQuery = q.Encode()

    resp, err := http.Get(search.String())

    if err != nil {
        return "", err
    }

    buf, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return "", err
    }

    result, err := json.Unmarshal(buf)

    if err != nil {
        return "", err
    }

    images, err := result.Get("responseData").Get("results").Array()

    if err != nil {
        return "", err
    }

    if len(images) > 0 {
        image    := images[rand.Intn(len(images))]
        imageMap := &json.Json{image}

        return imageMap.Get("unescapedUrl").MustString(), nil
    }

    return "", nil
}
