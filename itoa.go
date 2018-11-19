package utils

import "strconv"

func Itoa(i int64) string {
	return strconv.FormatInt(i,10)
}


func StringToInt64(s string) (int64,error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1,err
	}
	return i,nil
}

func Int64ToString(s int64) (string) {

	return strconv.FormatInt(s,10)
}

