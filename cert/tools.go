package cert

func StrToArray(str string) []string {
	if str == "" {
		return nil
	}
	return []string{str}
}
