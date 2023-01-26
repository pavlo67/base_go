package viewshtml

//func Pagination(limits []uint64, sortBy []string, opt *httplib.ReadOptionsHTTP) string {
//	if opt == nil {
//		return ""
//	}
//
//	if opt.CGIParams == "" {
//		opt.CGIParams += "?"
//	} else {
//		opt.CGIParams += "&"
//	}
//
//	paginationHTML := ""
//
//	var pageLength uint64
//	if len(limits) > 1 {
//		pageLength = limits[1]
//	} else if len(limits) > 0 {
//		pageLength = limits[0]
//	}
//
//	pageNum := opt.PageNum
//	if pageNum == 0 {
//		pageNum = 1
//	}
//
//	if pageLength > 0 && opt.AllCnt > pageLength {
//		var showNextPrevPage = 2
//		paginationHTML += "\n"
//		pages := opt.AllCnt / pageLength
//		if pages*pageLength < opt.AllCnt {
//			pages++
//		}
//		threePoints := false
//		for i := uint64(1); i <= pages; i++ {
//			if pages > 5 {
//				if i != 1 && i != pages {
//					if int(math.Abs(float64(pageNum-i))) > showNextPrevPage {
//						if !threePoints {
//							threePoints = true
//							paginationHTML += `...
//`
//						}
//						continue
//					}
//					threePoints = false
//				}
//			}
//			if i != pageNum {
//				href := opt.WithParams + opt.CGIParams +
//					"sort=" + strings.Join(sortBy, "+") +
//					"&page=" + strconv.FormatUint(i, 10)
//				paginationHTML += `
//				[<a href="` + href + `">` + strconv.FormatUint(i, 10) + `</a>]
//`
//			} else {
//				paginationHTML += "\n" + strconv.FormatUint(i, 10) + "\n"
//			}
//		}
//		paginationHTML += "\n"
//	}
//
//	return paginationHTML
//}
