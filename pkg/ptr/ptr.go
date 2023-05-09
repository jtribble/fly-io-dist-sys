package ptr

func ToInt(value int) *int {
	return &value
}

func ToString(value string) *string {
	return &value
}

func ToStringSlice(value []string) *[]string {
	return &value
}
