package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertStringSlice(slice []string) []types.String {
	result := make([]types.String, len(slice))
	for i, s := range slice {
		result[i] = types.StringValue(s)
	}
	return result
}

func GenerateStringID(namespace, name string) types.String {
	return types.StringValue(fmt.Sprintf("%s/%s", namespace, name))
}

func ParseStringID(ID types.String) (namespace, name string, err error) {
	chunks := strings.Split(ID.String(), "/")
	if len(chunks) != 2 {
		err = errors.New("wrong ID")
		return
	}

	return chunks[0], chunks[1], nil
}
