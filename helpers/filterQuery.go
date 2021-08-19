package helpers

import (
	"log"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	paramSeparator = "."
	expSeparator   = "_"
)

type Expression struct {
	Field     string
	Operator  string
	Values    []string
	Join      string
	Operator2 string
	Value2    string
}

// ParseFilterQueryString accepts a special syntax string {@param filter}, ids string and search string
// and tries to parse them to a single mongo db filter query
//
// The filter syntax can accept many expressions separated with the {expSeparator}.
//
// Each part of the expression is separated by the {paramSeparator}
// Each expresion has the following format {attribute}_eq_{trait}_{trait}_{trait}...
//
// Example: headwear_eq_Amish Hat_Copter Hat
func ParseFilterQueryString(filter string, ids string, search string) bson.M {
	// rarityscore_gte_10.2_and_lte_12.4;mainsetname_eq_Spartan;isvirgin_eq_false
	expressions := strings.Split(filter, paramSeparator)
	expArray := []Expression{}

	for _, expression := range expressions {
		exParts := strings.Split(expression, expSeparator)
		if len(exParts) >= 3 {
			field, operator, values := strings.ToLower(exParts[0]), exParts[1], exParts[2:]

			expArray = append(expArray, Expression{
				Field:    field,
				Operator: operator,
				Values:   values,
			})
		}
	}

	filters := buildFilters(expArray, ids, search)
	return filters
}

// buildFilter iterates over each parsed expression, parses it, creates a mongodb query and appends it to a global "$and" query
func buildFilters(expressions []Expression, ids string, search string) bson.M {
	searchFilters, finalFilter, orIdsFilter := bson.M{}, bson.M{}, bson.M{}
	idsFilterArray, finalFiltersArray := bson.A{}, bson.A{}

	//Parse and build search filter query
	if search != "" {
		searchFilters = ParseSearchQueryString(search)
		finalFiltersArray = append(finalFiltersArray, searchFilters)
	}

	//Parse and build ids filter query
	if ids != "" {
		idsFilters := ParseIdsFilter(ids)
		idsFilterArray = append(idsFilterArray, idsFilters...)
		orIdsFilter["$or"] = idsFilterArray
		finalFiltersArray = append(finalFiltersArray, orIdsFilter)
	}

	for _, exp := range expressions {
		if len(exp.Values) == 1 {
			trait := fixTraitNames(exp.Values[0])
			traitFilter := createEqBson(exp.Field, trait)
			finalFiltersArray = append(finalFiltersArray, traitFilter)
		} else if len(exp.Values) > 1 {
			orTraitsFilterArray, orFinalFilter := bson.A{}, bson.M{}
			for _, trait := range exp.Values {
				trait = fixTraitNames(trait)
				traitFilter := createEqBson(exp.Field, trait)
				orTraitsFilterArray = append(orTraitsFilterArray, traitFilter)
			}
			orFinalFilter["$or"] = orTraitsFilterArray
			finalFiltersArray = append(finalFiltersArray, orFinalFilter)

		}
	}
	finalFilter["$and"] = finalFiltersArray
	return finalFilter
}

//	TODO: Fix this ugly ass workaround here or in the frontend
func fixTraitNames(value string) string {
	if value == "Bow and Arrow" {
		value = "Bow & Arrow"
	} else if value == "Bow Tie and Suit" {
		value = "Bow Tie & Suit"
	} else if value == "Suit and Tie" {
		value = "Suit & Tie"
	} else if value == "Tennis Socks and Shoes" {
		value = "Tennis Socks & Shoes"
	}
	return value
}

// buildFilter iterates over each parsed expression, parses it, creates a mongodb query and appends it to global filter query
// func buildFilter(expressions []Expression, filter bson.M) bson.M {
// 	for _, exp := range expressions {
// 		switch exp.Join {
// 		case "":
// 			switch exp.Operator {
// 			case "eq":
// 				currBson := createEqBson(exp.Field, exp.Value)
// 				for k, v := range currBson {
// 					filter[k] = v
// 				}
// 			case "lt", "lte", "gt", "gte":
// 				currBson := createCompareBson(exp.Field, exp.Operator, exp.Value)
// 				for k, v := range currBson {
// 					filter[k] = v
// 				}
// 			}
// 		case "and", "or":
// 			var bson1, bson2 bson.M
// 			bson1 = createCompareBson(exp.Field, exp.Operator, exp.Value)
// 			bson2 = createCompareBson(exp.Field, exp.Operator2, exp.Value2)

// 			aBson := bson.A{bson1, bson2}
// 			filter["$"+exp.Join] = aBson
// 		}
// 	}

// 	return filter
// }

// createEqBson creates a mongodb filter if the operator is "eq"
func createEqBson(field string, value string) bson.M {
	returnBson := bson.M{}
	log.Println(value)
	if value == "true" || value == "false" {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			log.Println(err)
		} else {
			returnBson[field] = boolValue
		}
	} else {
		returnBson[field] = value
	}
	return returnBson
}

// createCompareBson creates a mongodb filter if the operator is lt, lte, gt, gte
func createCompareBson(field string, operator string, value string) bson.M {
	returnBson := bson.M{}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Println(err)
	} else {
		nestedBson := bson.M{"$" + operator: floatValue}
		returnBson[field] = nestedBson
	}
	return returnBson
}
