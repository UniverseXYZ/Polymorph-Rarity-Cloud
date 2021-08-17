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
	Value     string
	Join      string
	Operator2 string
	Value2    string
}

// ParseFilterQueryString accepts a special syntax string and tries to parse it to a mongo db filter query
//
// The syntax can accept many expressions separated with the expSeparator.
//
// Each part of the expression is separated by the paramSeparator
// Each expresion can be short and long one
//
// Short expressions consist of 3 parts: field, operator, value
//
// Example: mainsetname_eq_Spartan
//
// Long expressions consist of 6 parts: field, operator, value, join operator, operator 2, value 2
//
// Example: rarityscore_gte_10.2_and_lte_12.4
//
// Supported fields: any field from contants.PolymorphFields.go
//
// Supported operators: eq, lt, lte, gt, gte
//
// Supported join operators: and, or
func ParseFilterQueryString(filter string) bson.M {
	// rarityscore_gte_10.2_and_lte_12.4;mainsetname_eq_Spartan;isvirgin_eq_false
	expressions := strings.Split(filter, paramSeparator)
	expArray := []Expression{}
	filters := bson.M{}

	for _, expression := range expressions {
		exParts := strings.Split(expression, expSeparator)
		if len(exParts) == 3 {
			field, operator, value := strings.ToLower(exParts[0]), exParts[1], exParts[2]
			// TODO: Fix this ugly ass workaround here or in the frontend
			if value == "Bow and Arrow" {
				value = "Bow & Arrow"
			} else if value == "Bow Tie and Suit" {
				value = "Bow Tie & Suit"
			} else if value == "Suit and Tie" {
				value = "Suit & Tie"
			} else if value == "Tennis Socks and Shoes" {
				value = "Tennis Socks & Shoes"
			}

			expArray = append(expArray, Expression{
				Field:    field,
				Operator: operator,
				Value:    value,
			})
		} else if len(exParts) > 3 {
			field, _, values := strings.ToLower(exParts[0]), exParts[1], exParts[2:]
			aBson := bson.A{}
			for _, value := range values {
				// TODO: Fix this ugly ass workaround here or in the frontend
				if value == "Bow and Arrow" {
					value = "Bow & Arrow"
				} else if value == "Bow Tie and Suit" {
					value = "Bow Tie & Suit"
				} else if value == "Suit and Tie" {
					value = "Suit & Tie"
				} else if value == "Tennis Socks and Shoes" {
					value = "Tennis Socks & Shoes"
				}
				bson := createEqBson(field, value)
				aBson = append(aBson, bson)
			}
			filters["$or"] = aBson
		}
		// if len(exParts) == 6 {
		// 	field, operator, value, join, operator2, value2 := strings.ToLower(exParts[0]), exParts[1], exParts[2], exParts[3], exParts[4], exParts[5]

		// 	expArray = append(expArray, Expression{
		// 		Field:     field,
		// 		Operator:  operator,
		// 		Value:     value,
		// 		Join:      join,
		// 		Operator2: operator2,
		// 		Value2:    value2,
		// 	})

		// }
	}

	filters = buildFilter(expArray, filters)
	return filters
}

// buildFilter iterates over each parsed expression, parses it, creates a mongodb query and appends it to global filter query
func buildFilter(expressions []Expression, filter bson.M) bson.M {
	for _, exp := range expressions {
		switch exp.Join {
		case "":
			switch exp.Operator {
			case "eq":
				currBson := createEqBson(exp.Field, exp.Value)
				for k, v := range currBson {
					filter[k] = v
				}
			case "lt", "lte", "gt", "gte":
				currBson := createCompareBson(exp.Field, exp.Operator, exp.Value)
				for k, v := range currBson {
					filter[k] = v
				}
			}
		case "and", "or":
			var bson1, bson2 bson.M
			bson1 = createCompareBson(exp.Field, exp.Operator, exp.Value)
			bson2 = createCompareBson(exp.Field, exp.Operator2, exp.Value2)

			aBson := bson.A{bson1, bson2}
			filter["$"+exp.Join] = aBson
		}
	}

	return filter
}

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
