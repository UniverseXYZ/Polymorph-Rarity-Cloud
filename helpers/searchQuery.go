package helpers

import (
	"math"
	"strconv"

	"github.com/vikinatora/rarity-cloud-function/config"
	"github.com/vikinatora/rarity-cloud-function/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParseSearchQueryString accepts search string as parameter and build a mongodb filter for each filter from the config
//
// Config can be found in config/apiConfig.go
func ParseSearchQueryString(search string) bson.M {
	queries := bson.A{}

	for _, field := range config.SEARCH_QUERY_FIELDS {
		pattern := ""
		switch field {
		case constants.MorphFieldNames.TokenId,
			constants.MorphFieldNames.Rank:
			parsed, err := strconv.Atoi(search)
			if err == nil {
				queries = append(queries, bson.M{field: parsed})
			}
		case constants.MorphFieldNames.RarityScore:
			parsed, err := strconv.ParseFloat(search, 64)
			parsed = math.Floor(parsed*100) / 100
			if err == nil {
				queries = append(queries, bson.M{field: parsed})
			}

		default:
			pattern = search
			regex := primitive.Regex{Pattern: pattern, Options: "i"}
			regexFilter := bson.M{"$regex": regex}
			queries = append(queries, bson.M{field: regexFilter})
		}
	}
	orQuery := bson.M{"$or": queries}
	return orQuery
}

// getExactTokenPattern build a regex expression that will make an exact match with the specified stirng
func getExactTokenPattern(number string) string {
	return "(^|\\D)" + number + "(?!\\d)"
}
