package helper

import (
	"reflect"
	"sort"
	logger "walls-notification-service/internal/core/helper/log-helper"
)

type Pair struct {
	Key   string
	Value string
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func Sort(data map[string]string) ([]interface{}, error) {
	logger.LogEvent("INFO", "Sorting "+reflect.TypeOf(data).String()+" Data...")
	p := make(PairList, len(data))

	i := 0
	for k, v := range data {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	var a []interface{}

	for _, k := range p {
		a = append(a, struct {
			Key   string
			Value string
		}{Key: k.Key, Value: k.Value})
	}

	logger.LogEvent("INFO", reflect.TypeOf(data).String()+" Data Validated Successfully...")
	return a, nil
}
