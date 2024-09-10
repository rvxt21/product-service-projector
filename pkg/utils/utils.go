package utils

import (
	"products/internal/enteties"
	"strconv"
	"strings"
)

func ProcessIfIdsInt(ids string) ([]int, error) {
	res := []int{}
	s := strings.Split(ids, ",")
	for _, elem := range s {
		intVal, err := strconv.Atoi(elem)
		if err != nil {
			return nil, err
		}
		res = append(res, intVal)
	}

	return res, nil
}

func ConvertIntSliceToString(intSlice []int) string {
	valuesText := []string{}

	for _, val := range intSlice {
		text := strconv.Itoa(val)
		valuesText = append(valuesText, text)
	}
	result := strings.Join(valuesText, ",")
	return result
}

func IsEmptyFullProduct(structToCheck enteties.FullProductInfo) bool {
	return structToCheck == enteties.FullProductInfo{}
}
