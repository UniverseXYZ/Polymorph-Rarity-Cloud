package handlers

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"strconv"

	"github.com/vikinatora/rarity-cloud-function/config"
	"github.com/vikinatora/rarity-cloud-function/constants"
	"github.com/vikinatora/rarity-cloud-function/db"
	"github.com/vikinatora/rarity-cloud-function/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/aws/aws-lambda-go/events"
)

// GetPolymorphs endpoints returns polymorphs based on different filters that can be applied.
//
//If no polymorph is found returns empty response
//
//	Accepted query parameters:
//
// 		Take - int - Sets the number of results that should be returned
//
// 		Page - int - skips ((page - 1) * take) results
//
// 		SortField - string - sets field on which the results will be sorted. Default is polymorph id
//
// 		SortDir  - asc/desc - sets the sort direction of the results. Default is ascending
//
// 		Search - string - the string will be searched in different fields.
//
//		Searchable fields can be found in "apiConfig.go".
//
//		Filter - string - this query param requires special syntax in order to work.
//
//		See helpers.ParseFilterQueryString() for more information.
//
//		Example filter query: "rarityscore_gte_13.2_and_lte_20;isvirgin_eq_true;"
func GetPolymorphs(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		polymorphDBName := os.Getenv("POLYMORPH_DB")
		rarityCollectionName := os.Getenv("RARITY_COLLECTION")

		// Connect to DB
		collection, err := db.GetMongoDbCollection(polymorphDBName, rarityCollectionName)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}

		// Close connection to DB at the end of the function execution
		defer db.DisconnectDB()
		
		// Parse query params
		page, err := parseQueryParam("page", request)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}

		take, err := parseQueryParam("take", request)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}

		sortField, err := parseQueryParam("sortField", request)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}

		sortDir, err := parseQueryParam("sortDir", request)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}

		search, err := parseQueryParam("search", request)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}

		filter, err := parseQueryParam("filter", request)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}

		ids, err := parseQueryParam("ids", request)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}

		var filters = bson.M{}
		if filter != "" || ids != "" || search != "" {
			filters = helpers.ParseFilterQueryString(filter, ids, search)
		}

		var findOptions options.FindOptions

		removePrivateFields(&findOptions)

		takeInt, err := strconv.ParseInt(take, 10, 64)
		if err != nil || takeInt > config.RESULTS_LIMIT {
			takeInt = config.RESULTS_LIMIT
		}
		findOptions.SetLimit(takeInt)

		pageInt, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			pageInt = 1
		}

		findOptions.SetSkip((pageInt - 1) * takeInt)

		sortDirInt := 1

		if sortDir == "desc" {
			sortDirInt = -1
		}

		if sortField != "" {
			findOptions.SetSort(bson.D{{sortField, sortDirInt}, {constants.MorphFieldNames.TokenId, 1}})
		} else {
			findOptions.SetSort(bson.M{constants.MorphFieldNames.TokenId: sortDirInt})
		}

		curr, err := collection.Find(context.Background(), filters, &findOptions)
		if err != nil {
			return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
		}

		defer curr.Close(context.Background())

		var results []bson.M
		curr.All(context.Background(), &results)
	
		if results != nil {
			content, err := json.Marshal(results)
			if err != nil {
				return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 503}, nil
			}
			jsonResponse := string(content)
		
			return events.APIGatewayProxyResponse{ Headers: map[string]string{"Content-Type": "application/json"}, Body: jsonResponse, StatusCode: 200}, nil
		} else {
			return events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}, Body: "", StatusCode: 200}, nil
		}

}

// removePrivateFields removes internal fields that are of no interest to the users of the API.
//
// Configuration of these fields can be found in helpers.apiConfig.go
func removePrivateFields(findOptions *options.FindOptions) {
	noProjectionFields := bson.M{}
	for _, field := range config.MORPHS_NO_PROJECTION_FIELDS {
		noProjectionFields[field] = 0
	}
	findOptions.SetProjection(noProjectionFields)
}

func parseQueryParam(queryParam string, request events.APIGatewayProxyRequest ) (string, error) {
	pageValue, found := request.QueryStringParameters[queryParam]

	if found {
		// query parameters are typically URL encoded so to get the value
		value, err := url.QueryUnescape(pageValue)
		if err != nil {
			return "", err
		}
		return value, nil
	}

	return "", nil
}