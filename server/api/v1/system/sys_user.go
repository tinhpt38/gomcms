package system

import (
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Đăng nhập
// @Tags     Base
// @Summary  Đăng nhập người dùng
// @Produce   application/json
// @Param    data  body      systemReq.Login                                             true  "Tên người dùng, Mật khẩu, Mã xác nhận"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "Trả về thông tin người dùng, token, thời gian hết hạn"
// @Router   /base/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var l systemReq.Login
	err := c.ShouldBindJSON(&l)
	key := c.ClientIP()

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.LoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Kiểm tra xem mã xác nhận có được bật hay không
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha               // Có bật chống bùng nổ hay không
	openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // Thời gian chờ hết hạn bộ nhớ cache
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool = openCaptcha == 0 || openCaptcha < interfaceToInt(v)

	if !oc || (l.CaptchaId != "" && l.Captcha != "" && store.Verify(l.CaptchaId, l.Captcha, true)) {
		u := &system.SysUser{Username: l.Username, Password: l.Password}
		user, err := userService.Login(u)
		if err != nil {
			global.GVA_LOG.Error("Đăng nhập thất bại! Tên người dùng không tồn tại hoặc mật khẩu không đúng!", zap.Error(err))
			// Tăng số lần mã xác nhận
			global.BlackCache.Increment(key, 1)
			response.FailWithMessage("Tên người dùng không tồn tại hoặc mật khẩu không đúng", c)
			return
		}
		if user.Enable != 1 {
			global.GVA_LOG.Error("Đăng nhập thất bại! Người dùng bị cấm đăng nhập!")
			// Tăng số lần mã xác nhận
			global.BlackCache.Increment(key, 1)
			response.FailWithMessage("Người dùng bị cấm đăng nhập", c)
			return
		}
		b.TokenNext(c, *user)
		return
	}
	// Tăng số lần mã xác nhận
	global.BlackCache.Increment(key, 1)
	response.FailWithMessage("Mã xác nhận không đúng", c)
}

// TokenNext sau khi đăng nhập, phát hành jwt
func (b *BaseApi) TokenNext(c *gin.Context, user system.SysUser) {
	token, claims, err := utils.LoginToken(&user)
	if err != nil {
		global.GVA_LOG.Error("Lấy mã thông báo thất bại!", zap.Error(err))
		response.FailWithMessage("Lấy mã thông báo thất bại", c)
		return
	}
	if !global.GVA_CONFIG.System.UseMultipoint {
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "Đăng nhập thành công", c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			global.GVA_LOG.Error("Thiết lập trạng thái đăng nhập thất bại!", zap.Error(err))
			response.FailWithMessage("Thiết lập trạng thái đăng nhập thất bại", c)
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "Đăng nhập thành công", c)
	} else if err != nil {
		global.GVA_LOG.Error("Thiết lập trạng thái đăng nhập thất bại!", zap.Error(err))
		response.FailWithMessage("Thiết lập trạng thái đăng nhập thất bại", c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt không hợp lệ", c)
			return
		}
		if err := jwtService.SetRedisJWT(token, user.GetUsername()); err != nil {
			response.FailWithMessage("Thiết lập trạng thái đăng nhập thất bại", c)
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "Đăng nhập thành công", c)
	}
}

// Register
// @Tags     SysUser
// @Summary  Đăng ký tài khoản người dùng
// @Produce   application/json
// @Param    data  body      systemReq.Register                                            true  "Tên người dùng, Biệt danh, Mật khẩu, ID vai trò"
// @Success  200   {object}  response.Response{data=systemRes.SysUserResponse,msg=string}  "Đăng ký tài khoản người dùng, trả về thông tin người dùng"
// @Router   /user/admin_register [post]
func (b *BaseApi) Register(c *gin.Context) {
	var r systemReq.Register
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var authorities []system.SysAuthority
	for _, v := range r.AuthorityIds {
		authorities = append(authorities, system.SysAuthority{
			AuthorityId: v,
		})
	}
	user := &system.SysUser{Username: r.Username, NickName: r.NickName, Password: r.Password, HeaderImg: r.HeaderImg, AuthorityId: r.AuthorityId, Authorities: authorities, Enable: r.Enable, Phone: r.Phone, Email: r.Email}
	userReturn, err := userService.Register(*user)
	if err != nil {
		global.GVA_LOG.Error("Đăng ký thất bại!", zap.Error(err))
		response.FailWithDetailed(systemRes.SysUserResponse{User: userReturn}, "Đăng ký thất bại", c)
		return
	}
	response.OkWithDetailed(systemRes.SysUserResponse{User: userReturn}, "Đăng ký thành công", c)
}

// ChangePassword
// @Tags      SysUser
// @Summary   Người dùng thay đổi mật khẩu
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      systemReq.ChangePasswordReq    true  "Tên người dùng, Mật khẩu cũ, Mật khẩu mới"
// @Success   200   {object}  response.Response{msg=string}  "Người dùng thay đổi mật khẩu"
// @Router    /user/changePassword [post]
func (b *BaseApi) ChangePassword(c *gin.Context) {
	var req systemReq.ChangePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.ChangePasswordVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	u := &system.SysUser{GVA_MODEL: global.GVA_MODEL{ID: uid}, Password: req.Password}
	_, err = userService.ChangePassword(u, req.NewPassword)
	if err != nil {
		global.GVA_LOG.Error("Thay đổi thất bại!", zap.Error(err))
		response.FailWithMessage("Thay đổi thất bại, mật khẩu cũ không khớp với tài khoản hiện tại", c)
		return
	}
	response.OkWithMessage("Thay đổi thành công", c)
}

// GetUserList
// @Tags      SysUser
// @Summary   Lấy danh sách người dùng theo trang
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "Trang, Số lượng mục trên mỗi trang"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "Lấy danh sách người dùng theo trang, trả về danh sách, tổng số, trang, số lượng mục trên mỗi trang"
// @Router    /user/getUserList [post]
func (b *BaseApi) GetUserList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := userService.GetUserInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Lấy danh sách thất bại!", zap.Error(err))
		response.FailWithMessage("Lấy danh sách thất bại", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Lấy danh sách thành công", c)
}

// SetUserAuthority
// @Tags      SysUser
// @Summary   Thay đổi quyền của người dùng
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuth          true  "UUID người dùng, ID vai trò"
// @Success   200   {object}  response.Response{msg=string}  "Thay đổi quyền của người dùng"
// @Router    /user/setUserAuthority [post]
func (b *BaseApi) SetUserAuthority(c *gin.Context) {
	var sua systemReq.SetUserAuth
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if UserVerifyErr := utils.Verify(sua, utils.SetUserAuthorityVerify); UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	userID := utils.GetUserID(c)
	err = userService.SetUserAuthority(userID, sua.AuthorityId)
	if err != nil {
		global.GVA_LOG.Error("Thay đổi thất bại!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims := utils.GetUserInfo(c)
	j := &utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)} // Chữ ký duy nhất
	claims.AuthorityId = sua.AuthorityId
	if token, err := j.CreateToken(*claims); err != nil {
		global.GVA_LOG.Error("Thay đổi thất bại!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		c.Header("new-token", token)
		c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt.Unix(), 10))
		utils.SetToken(c, token, int((claims.ExpiresAt.Unix()-time.Now().Unix())/60))
		response.OkWithMessage("Thay đổi thành công", c)
	}
}

// SetUserAuthorities
// @Tags      SysUser
// @Summary   Thiết lập quyền của người dùng
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuthorities   true  "UUID người dùng, ID vai trò"
// @Success   200   {object}  response.Response{msg=string}  "Thiết lập quyền của người dùng"
// @Router    /user/setUserAuthorities [post]
func (b *BaseApi) SetUserAuthorities(c *gin.Context) {
	var sua systemReq.SetUserAuthorities
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.SetUserAuthorities(sua.ID, sua.AuthorityIds)
	if err != nil {
		global.GVA_LOG.Error("Thiết lập thất bại!", zap.Error(err))
		response.FailWithMessage("Thiết lập thất bại", c)
		return
	}
	response.OkWithMessage("Thiết lập thành công", c)
}

// DeleteUser
// @Tags      SysUser
// @Summary   Xóa người dùng
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "ID người dùng"
// @Success   200   {object}  response.Response{msg=string}  "Xóa người dùng"
// @Router    /user/deleteUser [delete]
func (b *BaseApi) DeleteUser(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	jwtId := utils.GetUserID(c)
	if jwtId == uint(reqId.ID) {
		response.FailWithMessage("Xóa thất bại, không thể tự tử", c)
		return
	}
	err = userService.DeleteUser(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("Xóa thất bại!", zap.Error(err))
		response.FailWithMessage("Xóa thất bại", c)
		return
	}
	response.OkWithMessage("Xóa thành công", c)
}

// SetUserInfo
// @Tags      SysUser
// @Summary   Thiết lập thông tin người dùng
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, Tên người dùng, Biệt danh, Link ảnh đại diện"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Thiết lập thông tin người dùng"
// @Router    /user/setUserInfo [put]
func (b *BaseApi) SetUserInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(user, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if len(user.AuthorityIds) != 0 {
		err = userService.SetUserAuthorities(user.ID, user.AuthorityIds)
		if err != nil {
			global.GVA_LOG.Error("Thiết lập thất bại!", zap.Error(err))
			response.FailWithMessage("Thiết lập thất bại", c)
			return
		}
	}
	err = userService.SetUserInfo(system.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		SideMode:  user.SideMode,
		Enable:    user.Enable,
	})
	if err != nil {
		global.GVA_LOG.Error("Thiết lập thất bại!", zap.Error(err))
		response.FailWithMessage("Thiết lập thất bại", c)
		return
	}
	response.OkWithMessage("Thiết lập thành công", c)
}

// SetSelfInfo
// @Tags      SysUser
// @Summary   Thiết lập thông tin người dùng
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, Tên người dùng, Biệt danh, Link ảnh đại diện"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Thiết lập thông tin người dùng"
// @Router    /user/SetSelfInfo [put]
func (b *BaseApi) SetSelfInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user.ID = utils.GetUserID(c)
	err = userService.SetSelfInfo(system.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		SideMode:  user.SideMode,
		Enable:    user.Enable,
	})
	if err != nil {
		global.GVA_LOG.Error("Thiết lập thất bại!", zap.Error(err))
		response.FailWithMessage("Thiết lập thất bại", c)
		return
	}
	response.OkWithMessage("Thiết lập thành công", c)
}

// GetUserInfo
// @Tags      SysUser
// @Summary   Lấy thông tin người dùng
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "Lấy thông tin người dùng"
// @Router    /user/getUserInfo [get]
func (b *BaseApi) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	ReqUser, err := userService.GetUserInfo(uuid)
	if err != nil {
		global.GVA_LOG.Error("Lấy thông tin thất bại!", zap.Error(err))
		response.FailWithMessage("Lấy thông tin thất bại", c)
		return
	}
	response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "Lấy thông tin thành công", c)
}

// ResetPassword
// @Tags      SysUser
// @Summary   Đặt lại mật khẩu người dùng
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      system.SysUser                 true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "Đặt lại mật khẩu người dùng"
// @Router    /user/resetPassword [post]
func (b *BaseApi) ResetPassword(c *gin.Context) {
	var user system.SysUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.ResetPassword(user.ID)
	if err != nil {
		global.GVA_LOG.Error("Đặt lại thất bại!", zap.Error(err))
		response.FailWithMessage("Đặt lại thất bại"+err.Error(), c)
		return
	}
	response.OkWithMessage("Đặt lại thành công", c)
}
