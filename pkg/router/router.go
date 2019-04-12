package router

import (
    "errors"
    "regexp"
    "strings"
)

type ServiceRouter struct {
}

func (router *ServiceRouter) ExtractServiceName(url string) (string, error) {
    regex, err := regexp.Compile("^([a-zA-Z]+)://")
    if err != nil {
        return "", errors.New("could not compile regex")
    }
    match := regex.FindStringSubmatch(url)
    if len(match) <= 1 {
        return "", errors.New("could not find any service part")
    }
    return match[1], nil
}


func (router *ServiceRouter) ExtractArguments(url string) ([]string, error) {
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