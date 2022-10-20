package services

import (
	"insight-api/dbs"
)

func QueryDevelopers(keyword string, page, count int) *Developers {
	developerDao := dbs.DeveloperDao{}
	developerdbs, err := developerDao.QueryListByPage(keyword, page, count)
	developers := &Developers{
		PageInfo: &PageInfo{
			Page:  page,
			Count: count,
		},
	}
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
				AppCount:       developerdb.AppCount,
			})
		}
	}
	return developers
}

func GetDeveloperById(id, title string) *Developer {
	developerDao := dbs.DeveloperDao{}
	dev := &Developer{
		Id:    id,
		Title: title,
	}
	developerdb, err := developerDao.FindById(id)
	if err == nil {
		dev.Trade = developerdb.Industry
		dev.FoundedTime = developerdb.FoundedYear
		dev.AddressArea = developerdb.AddressArea
		dev.FinancingRound = developerdb.FinancingRound
		dev.LogoUrl = developerdb.LogoUrl
		dev.AppCount = developerdb.AppCount
	}
	return dev
}
