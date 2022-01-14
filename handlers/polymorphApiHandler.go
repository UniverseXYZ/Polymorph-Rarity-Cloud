package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"github.com/vikinatora/rarity-cloud-function/config"
	"github.com/vikinatora/rarity-cloud-function/constants"
	"github.com/vikinatora/rarity-cloud-function/db"
	"github.com/vikinatora/rarity-cloud-function/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func GetPolymorphs(polymorphDBName string, rarityCollectionName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		collection, err := db.GetMongoDbCollection(polymorphDBName, rarityCollectionName)
		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, err)
			log.Errorln(err)
			return
		}

		defer db.DisconnectDB()

		page := r.URL.Query().Get("page")
		take := r.URL.Query().Get("take")
		sortField := r.URL.Query().Get("sortField")
		sortDir := r.URL.Query().Get("sortDir")
		search := r.URL.Query().Get("search")
		filter := r.URL.Query().Get("filter")
		ids := r.URL.Query().Get("ids")

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
			render.Status(r, 500)
			render.JSON(w, r, err)
			log.Println(err)
		}

		defer curr.Close(context.Background())

		var results []bson.M
		// var jsonResults []byte
		curr.All(context.Background(), &results)
		if results != nil {
			render.JSON(w, r, results)
		} else {
			render.JSON(w, r, []bson.M{})
		}

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
