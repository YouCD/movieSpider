package tools

//
// RemoveSpaceItem
//  @Description: 去除数组中的空格
//  @param a
//  @return ret
//
func RemoveSpaceItem(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
