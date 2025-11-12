package properties

// type PropertyHandler struct {
// 	// UserSvc *svcs.UserSvc
// }

// func NewUserHandler(userSvc *svcs.UserSvc) *CommonPropertyHandler {
// 	return &CommonPropertyHandler{UserSvc: userSvc}
// }

// func (c *CommonPropertyHandler) GetAllProperties(w http.ResponseWriter, r *http.Request) {
// 	params := r.URL.Query()
// 	page, limit := 1, 10
// 	if p, err := strconv.Atoi(params.Get("page")); err != nil {
// 		page = p
// 	}
// 	if l, err := strconv.Atoi(params.Get("limit")); err != nil {
// 		limit = l
// 	}
// 	offset := (page - 1) * limit

// 	// userType:=
// }
