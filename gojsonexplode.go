package gojsonexplode

// TODO: cannot have delimiter be " or \

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Exploder struct {
	delimiter string
}

func (e *Exploder) explode_list(l []interface{}, parent string) (map[string]interface{}, error) {
	var err error
	var key string
	j := make(map[string]interface{})
	for k, i := range l {
		if len(parent) > 0 {
			key = parent + e.delimiter + strconv.Itoa(k)
		} else {
			key = strconv.Itoa(k)
		}
		switch v := i.(type) {
		case nil:
			j[key] = v
		case int:
			j[key] = v
		case float64:
			j[key] = v
		case string:
			j[key] = v
		case bool:
			j[key] = v
		case []interface{}:
			out := make(map[string]interface{})
			out, err = e.explode_list(v, key)
			if err != nil {
				return nil, err
			}
			for newkey, value := range out {
				j[newkey] = value
			}
		case map[string]interface{}:
			out := make(map[string]interface{})
			out, err = e.explode_map(v, key)
			if err != nil {
				return nil, err
			}
			for newkey, value := range out {
				j[newkey] = value
			}
		default:
			// do nothing
		}
	}
	return j, nil
}

func (e *Exploder) explode_map(m map[string]interface{}, parent string) (map[string]interface{}, error) {
	var err error
	j := make(map[string]interface{})
	for k, i := range m {
		if len(parent) > 0 {
			k = parent + e.delimiter + k
		}
		switch v := i.(type) {
		case nil:
			j[k] = v
		case int:
			j[k] = v
		case float64:
			j[k] = v
		case string:
			j[k] = v
		case bool:
			j[k] = v
		case []interface{}:
			out := make(map[string]interface{})
			out, err = e.explode_list(v, k)
			if err != nil {
				return nil, err
			}
			for key, value := range out {
				j[key] = value
			}
		case map[string]interface{}:
			out := make(map[string]interface{})
			out, err = e.explode_map(v, k)
			if err != nil {
				return nil, err
			}
			for key, value := range out {
				j[key] = value
			}
		default:
			//nothing
		}
	}
	return j, nil
}

func explodejson(s string, d string) (string, error) {
	var input interface{}
	var exploded map[string]interface{}
	var out []byte
	var err error
	b := []byte(s)
	err = json.Unmarshal(b, &input)
	if err != nil {
		return "", err
	}
	exploder := Exploder{d}
	switch t := input.(type) {
	case map[string]interface{}:
		exploded, err = exploder.explode_map(t, "")
		if err != nil {
			return "", err
		}
	case []interface{}:
		exploded, err = exploder.explode_list(t, "")
		if err != nil {
			return "", err
		}
	default:
		// How did we get here? It is impossible!!
		return "", errors.New("Possible error in JSON")
	}
	out, err = json.Marshal(exploded)
	if err != nil {
		return "", err
	}
	return string(out), nil

}
