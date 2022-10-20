package utils

import (
	"encoding/base32"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func ParseInt64(value string) (int64, error) {
	ret, err := strconv.ParseInt(value, 10, 64)
	return ret, err
}

func ParseInt(value string) (int, error) {
	ret, err := strconv.Atoi(value)
	return ret, err
}

func GetUuid() string {
	id := uuid.NewV4()
	return id.String()
}

func GetClearUuid() string {
	id := GetUuid()
	return strings.ReplaceAll(id, "-", "")
}

var directory string = "23456789abcdefghigkmnpqrstuvwxyz"

func PruneUuid(uuidStr string) (string, error) {
	str := strings.ReplaceAll(uuidStr, "-", "")
	if len(str) != 32 {
		return "", errors.New("Illegal string " + uuidStr)
	}
	src := []byte(str)
	dst := make([]byte, len(src)/2)

	_, err := hex.Decode(dst, src)
	if err != nil {
		return "", err
	}

	xs := base32.NewEncoding(directory).EncodeToString(dst)
	return xs[:26], nil
}

func Parse2Uuid(pruneUuid string) (string, error) {
	val := pruneUuid + "======"
	bs, err := base32.NewEncoding(directory).DecodeString(val)
	if err != nil {
		return "", err
	}
	ret := hex.EncodeToString(bs)
	return ret, nil
}
