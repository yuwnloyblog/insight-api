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
			developers.Items = append(developers.Items, &Developer{
				Id:             developerdb.ID,
				Title:          developerdb.Title,
				Trade:          developerdb.Industry,
				FoundedTime:    developerdb.FoundedYear,
				AddressArea:    developerdb.AddressArea,
				FinancingRound: developerdb.FinancingRound,
				LogoUrl:        developerdb.LogoUrl,
			})
		}
	}
	return developers
}

func GetDeveloperById(id string) *Developer {
	developerDao := dbs.DeveloperDao{}
	developerdb, err := developerDao.FindById(id)
	if err == nil {
		return &Developer{
			Id:             developerdb.ID,
			Title:          developerdb.Title,
			Trade:          developerdb.Industry,
			FoundedTime:    developerdb.FoundedYear,
			AddressArea:    developerdb.AddressArea,
			FinancingRound: developerdb.FinancingRound,
			LogoUrl:        developerdb.LogoUrl,
		}
	}
	return nil
}

func GetDeveloperByIdStr(idStr string) *Developer {
	return GetDeveloperById(idStr)
}
