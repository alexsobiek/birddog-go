package birddog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type API struct {
	Host string
}

func NewAPI(host string) *API {
	// check if host starts with http(s)://, if not, add it
	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}

	if !strings.HasSuffix(host, ":8080") {
		host = host + ":8080"
	}

	// check if host ends with /, if not, add it
	if !strings.HasSuffix(host, "/") {
		host = host + "/"
	}

	return &API{
		Host: host,
	}
}

func (a *API) decode(res *http.Response, target interface{}) error {
	defer res.Body.Close()

	if target == nil {
		return nil
	}

	// Check if target is a pointer
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	elemValue := targetValue.Elem()

	switch elemValue.Kind() {
	case reflect.String:
		// Handle plain text response
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		elemValue.SetString(string(body))
	case reflect.Map:
		// Ensure the map is of type map[string]string
		if elemValue.Type().Key().Kind() == reflect.String && elemValue.Type().Elem().Kind() == reflect.String {
			if err := json.NewDecoder(res.Body).Decode(target); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported map type: %s", elemValue.Type())
		}
	case reflect.Struct:
		// Handle JSON response
		if err := json.NewDecoder(res.Body).Decode(target); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported target type: %s", elemValue.Kind())
	}

	return nil
}

func (a *API) encode(body interface{}) (io.Reader, error) {
	if body == nil {
		return bytes.NewBuffer([]byte{}), nil
	}

	// Check if target is a pointer
	targetValue := reflect.ValueOf(body)
	if targetValue.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("body must be a pointer")
	}

	elemValue := targetValue.Elem()

	switch elemValue.Kind() {
	case reflect.String:
		// Handle plain text response
		return strings.NewReader(elemValue.String()), nil
	case reflect.Struct:
		// Handle JSON response
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(b), nil
	default:
		return nil, fmt.Errorf("unsupported body type: %s", elemValue.Kind())
	}
}

func (a *API) get(endpoint string, target interface{}) (interface{}, error) {
	res, err := http.Get(fmt.Sprintf("%s%s", a.Host, endpoint))

	if err != nil {
		return nil, err
	}

	if err = a.decode(res, target); err != nil {
		return nil, err
	}

	return target, nil
}

func (a *API) post(endpoint string, body interface{}, target interface{}) (interface{}, error) {
	// create new reader from body

	r, err := a.encode(body)

	if err != nil {
		return nil, err
	}

	res, err := http.Post(fmt.Sprintf("%s%s", a.Host, endpoint), "application/json", r)

	if err != nil {
		return nil, err
	}

	if err = a.decode(res, target); err != nil {
		return nil, err
	}

	return target, nil
}
