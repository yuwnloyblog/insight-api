package services

import (
	"insight-api/dbs"
	"insight-api/utils"
	"math"
)

func QueryDevelopers(keyword string, startStr string, count int) *Developers {
	developerDao := dbs.DeveloperDao{}
	start, err := utils.Decode(startStr)
	if err != nil {
		start = math.MaxInt32
	}
	developerdbs, err := developerDao.QueryList(keyword, start, count)
	developers := &Developers{}
	if err == nil {
		l := len(developerdbs)
		if l >= count {
			developers.HasMore = true
		}
		developers.Items = make([]*Developer, 0)
		for _, developerdb := range developerdbs {
			isStr, _ := utils.Encode(developerdb.ID)
			developers.Items = append(developers.Items, &Developer{
				Id:             isStr,
				Title:          developerdb.Title,
				Trade:          developerdb.Trade,
				FoundedTime:    developerdb.FoundedTime.String(),
				AddressArea:    developerdb.AddressArea,
				FinancingRound: developerdb.FinancingRound,
				LogoUrl:        developerdb.LogoUrl,
			})
		}
	}
	return developers
}

func GetDeveloperById(id int64) *Developer {
	developerDao := dbs.DeveloperDao{}
	developerdb, err := developerDao.FindById(id)
	if err == nil {
		idStr, _ := utils.Encode(developerdb.ID)
		return &Developer{
			Id:             idStr,
			Title:          developerdb.Title,
			Trade:          developerdb.Trade,
			FoundedTime:    developerdb.FoundedTime.String(),
			AddressArea:    developerdb.AddressArea,
			FinancingRound: developerdb.FinancingRound,
			LogoUrl:        developerdb.LogoUrl,
		}
	}
	return nil
}

func GetDeveloperByIdStr(idStr string) *Developer {
	id, err := utils.Decode(idStr)
	if err == nil {
		return GetDeveloperById(id)
	}
	return nil
}
