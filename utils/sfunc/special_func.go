package sfunc

// InSlice seperti fungsi in, apakah target ada didalam slice.
// return true jika ada
func InSlice(target string, slice []string) bool {
	if len(slice) == 0 {
		return false
	}

	for _, value := range slice {
		if target == value {
			return true
		}
	}

	return false
}

// AllValueInSliceIsValid memasukkan input request berupa slice dan
// membandingkan isi slicenya apakah tersedia untuk digunakan
func AllValueInSliceIsValid(inputSlice []string, validSlice []string) bool {
	for _, input := range inputSlice {
		if !InSlice(input, validSlice) {
			return false
		}
	}
	return true
}
