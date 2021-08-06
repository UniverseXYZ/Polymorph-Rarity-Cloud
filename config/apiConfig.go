package config

import "github.com/vikinatora/rarity-cloud-function/constants"

var SEARCH_QUERY_FIELDS []string = []string{
	constants.MorphFieldNames.Rank,
	constants.MorphFieldNames.TokenId,
	constants.MorphFieldNames.RarityScore,
	constants.MorphFieldNames.Headwear,
	constants.MorphFieldNames.Eyewear,
	constants.MorphFieldNames.Torso,
	constants.MorphFieldNames.Pants,
	constants.MorphFieldNames.Footwear,
	constants.MorphFieldNames.LeftHand,
	constants.MorphFieldNames.RightHand,
	constants.MorphFieldNames.Character,
	constants.MorphFieldNames.MainSetName,
	constants.MorphFieldNames.SecSetName,
}

var MORPHS_NO_PROJECTION_FIELDS []string = []string{
	constants.MorphFieldNames.ObjId,
	constants.MorphFieldNames.OldGenes,
}

const RESULTS_LIMIT int64 = 10000
