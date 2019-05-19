package plugins

import (
    "errors"
    "regexp"
    "strings"
)

// ExtractArguments extracts the arguments from a notification url, i.e everything following the initial ://
func ExtractArguments(url string) ([]string, error) {
    regex, err := regexp.Compile("^[a-zA-Z]+://(.*)$")
    if err != nil {
        return nil, errors.New("could not compile regex")
    }
    match := regex.FindStringSubmatch(url)
    if len(match[1]) <= 0 {
        return nil, errors.New("could not extract any arguments")
    }
    return strings.Split(match[1], "/"), nil
}