package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"net/http"
	"unicode/utf8"
	"webook/internal/domain"
	"webook/internal/service"
)

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp // 预编译
	svc         *service.UserService
}

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		// 在引入三方库之后，要加一个option regexp.None
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:         svc,
	}
}
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", h.SignUp)
	// 若又很多个相同前缀，可以
	ug := server.Group("/users")
	// POST /users/login
	ug.POST("/login", h.Login)
	// POST /users/profile
	ug.GET("/profile", h.Profile)
	// POST /users/edit
	ug.POST("/edit", h.Edit)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json: "email"`
		Password        string `json: "password"`
		ConfirmPassword string `json: "confirmPassword"`
	}

	var req SignupReq
	// Bind根据content-type解决数据
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 原本正则匹配方法 -> 进阶预编译
	//isEmail, err := regexp.Match(emailRegexPattern, []byte(req.Email))
	// 预编译 -> go自带正则较弱，这里引用别的库
	//isEmail := h.emailExp.Match([]byte(req.Email))
	isEmail, err := h.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入密码不一致")
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱不正确")
		return
	}

	isPassword, err := h.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码格式不正确，包含数字，特殊字符且长度大于8位")
		return
	}
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrorDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱已经注册!")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		// 返回400 ， 可以不返回错误
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			// 15min
			MaxAge: 900,
			//HttpOnly: true,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "session保存错误！")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码不正确")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	id := sess.Get("userId").(int64)
	u, err := h.svc.Select(ctx, id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	type User struct {
		NickName string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}
	ctx.JSON(http.StatusOK, User{
		NickName: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday,
	})

}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
		Nickname string `json:"nickname"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		// 返回400 ， 可以不返回错误
		return
	}
	// 判断自我介绍和昵称是否符合长度
	if utf8.RuneCountInString(req.AboutMe) > 256 {
		ctx.String(http.StatusOK, "自我介绍长度不大于256！")
	}

	if utf8.RuneCountInString(req.Nickname) > 10 {
		ctx.String(http.StatusOK, "自我介绍长度不大于10！")
	}

	err := h.svc.Edit(ctx, domain.User{
		AboutMe:  req.AboutMe,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
	})
	if err != nil {
		ctx.String(http.StatusOK, "更新失败")
	}
	ctx.String(http.StatusOK, "更新成功")
}
