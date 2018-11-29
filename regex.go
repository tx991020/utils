package utils

import "regexp"

func getRegexp(pattern string) (*regexp.Regexp, error) {

	if r, err := regexp.Compile(pattern); err == nil {
		return r, nil
	} else {
		return nil, err
	}
}


func Validate(pattern string) error {
	_, err := getRegexp(pattern)
	return err
}


//正则表达式是否匹配
func IsMatch(pattern string, src []byte) bool {
	if r, err := getRegexp(pattern); err == nil {
		return r.Match(src)
	}
	return false
}

func IsMatchString(pattern string, src string) bool {
	return IsMatch(pattern, []byte(src))
}




// 正则匹配，并返回匹配的列表
func MatchString(pattern string, src string) ([]string, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.FindStringSubmatch(src), nil
	} else {
		return nil, err
	}
}

func MatchAllString(pattern string, src string) ([][]string, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.FindAllStringSubmatch(src, -1), nil
	} else {
		return nil, err
	}
}



// 正则替换(全部替换)
func Replace(pattern string, replace, src []byte) ([]byte, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.ReplaceAll(src, replace), nil
	} else {
		return nil, err
	}
}

// 正则替换(全部替换)，字符串
func ReplaceString(pattern, replace, src string) (string, error) {
	r, e := Replace(pattern, []byte(replace), []byte(src))
	return string(r), e
}