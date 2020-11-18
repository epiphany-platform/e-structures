package to

func StrPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func BooPtr(b bool) *bool {
	return &b
}
