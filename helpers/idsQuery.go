package helpers

import (
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// ParseIdsFilter accepts a string of ids delimited by commas. It parses the ids and build mongo db query filters.
//
// Return an bson.M containing the query filters that will return results from db only for these ids
func ParseIdsFilter(ids string) bson.A {
	arrayIds := strings.Split(ids, ",")
	aBson := bson.A{}
	for _, id := range arrayIds {
		numId, err := strconv.Atoi(id)
		if err == nil {
			returnBson := bson.M{}
			returnBson["tokenid"] = numId
			aBson = append(aBson, returnBson)
		}
	}

	return aBson
}
