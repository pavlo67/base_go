package views_html

//func HTMLPagination(opt *db.ReadAllHTTPOptions) string {
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
//	maxCnt := uint64(0)
//	pageLength := uint64(0)
//	if len(opt.Limits) > 0 {
//		maxCnt += opt.Limits[0]
//		pageLength = opt.Limits[0]
//		if len(opt.Limits) > 1 {
//			maxCnt += opt.Limits[1]
//			pageLength = opt.Limits[1]
//		}
//	}
//
//	if maxCnt > 0 && opt.AllCnt > maxCnt {
//		var showNextPrevPage = 2
//		paginationHTML += `
//`
//		pages := opt.AllCnt / pageLength
//		if pages*pageLength < opt.AllCnt {
//			pages++
//		}
//		threePoints := false
//		for i := uint64(1); i <= pages; i++ {
//			if pages > 5 {
//				if i != 1 && i != pages {
//					if int(math.Abs(float64(opt.PageNum-i+1))) > showNextPrevPage {
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
//			if i-1 != opt.PageNum {
//				href := opt.Path + opt.CGIParams + "sort=" + strings.Join(opt.SortBy, "+") + "&page=" + strconv.FormatUint(i-1, 10)
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
