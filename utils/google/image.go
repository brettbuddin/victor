package google

import (
    "net/http"
    "net/url"
    "io/ioutil"
    "math/rand"
    "encoding/json"
)

type ImageResult struct {
    UnescapedUrl string
}

type ImageResults struct {
    Results []ImageResult
}

type ImageResponseDate struct {
    ResponseData ImageResults
}

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

    var result ImageResponseDate

    err = json.Unmarshal(buf, &result)

    if err != nil {
        return "", err
    }

    images := result.ResponseData.Results

    if len(images) > 0 {
        return images[rand.Intn(len(images))].UnescapedUrl, nil
    }

    return "", nil
}
