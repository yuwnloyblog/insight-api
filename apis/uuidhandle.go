package apis

import "insight-api/utils"

func EncodeUuid(uuid string) string {
	ret, err := utils.PruneUuid(uuid)
	if err == nil && ret != "" {
		return ret
	}
	return uuid
}

func DecodeUuid(pruneUuid string) string {
	//ret,err := utils.get
	return ""
}
