package login_service_v1

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"strconv"
	common "test.com/project-common"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/jwts"
	user_grpc "test.com/project-grpc/user/login"
	"test.com/project-user/config"
	"test.com/project-user/internal/dao"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/data/organization"
	"test.com/project-user/internal/database"
	"test.com/project-user/internal/database/tran"
	"test.com/project-user/internal/repo"
	"test.com/project-user/pkg/model"
	"time"
)

func New() *LoginService {
	return &LoginService{
		cache:            dao.Rc,
		organizationRepo: dao.NewOrganizationDao(),
		memberRepo:       dao.NewMemberDao(),
		transaction:      dao.NewTransaction(),
	}
}

type LoginService struct { //注册登录微服务【类】。
	// 微服务Server 1.需要继承Unimplemented_xxx类(_grpc.pb.go文件中的),然后这个类还需要
	//2.实现你proto文件中定义的那些函数。
	//那么这个微服务类才是可以注册到你grpcServer中的Server
	user_grpc.UnimplementedLoginServiceServer
	memberRepo       repo.MemberRepo       //member表
	organizationRepo repo.OrganizationRepo //organization表
	transaction      tran.Transaction
	cache            repo.Cache //自己新增参数,这是一个自定义的缓存
}

func (lg *LoginService) GetCaptcha(c context.Context, msg *user_grpc.CaptchaRequest) (*user_grpc.CaptchaResponse, error) {
	mobile := msg.Mobile
	if !common.VerifyMobile(mobile) {
		return nil /*model.IllegalMobile*/, errs.GrpcError(model.IllegalMobile)
	}
	code := RandomCaptCha()
	go func() { //发送短信
		zap.L().Info("api发送短信 info") //模拟发送成功，输出到日志
		//logs.LG.Debug("api发送短信 debug")
		//zap.L().Error("api发送短信 error")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := lg.cache.Put(ctx, "REGISTER:"+mobile, code, 15*time.Minute) //存入redis
		if err != nil {
			log.Printf("验证码存入redis出错,cause by: %v\n", err)
		}
	}()
	time.Sleep(1 * time.Second) //这样来做并发控制，不行
	return &user_grpc.CaptchaResponse{Code: code}, nil
}

func RandomCaptCha() string {
	rand.Seed(time.Now().UnixNano()) // 使用当前时间作为随机数种子

	min := 100000
	max := 999999
	randomNumber := rand.Intn(max-min+1) + min

	return strconv.Itoa(randomNumber)
}

func (lg *LoginService) Register(ctx context.Context, msg *user_grpc.RegisterRequest) (*user_grpc.RegisterResponse, error) {
	//校验参数
	//fmt.Println("获取的电话号码", msg.Mobile)
	c := context.Background()
	//校验验证码，是否已经存在
	redisCode, err := lg.cache.Get(c, "REGISTER:"+msg.Mobile)
	if err == redis.Nil {
		//zap.L().Error()
		return nil, errs.GrpcError(model.CaptchaNotExist)
	}
	if err != nil {
		zap.L().Error("redis查库出错:", zap.Error(err))
		return nil, errs.GrpcError(model.InCorrectCaptcha)
	}

	if redisCode != msg.Captcha {
		return nil, errs.GrpcError(model.InCorrectCaptcha)
	}

	//邮箱，账号，手机号是否被注册
	exist, err := lg.memberRepo.GetMemberByEmail(c, msg.Email)
	if err != nil {
		zap.L().Error("DB错误", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExisted)
	}

	exist, err = lg.memberRepo.GetMemberByAccount(c, msg.Name)
	if err != nil {
		zap.L().Error("DB错误", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExisted)
	}

	exist, err = lg.memberRepo.GetMemberByMobile(c, msg.Mobile)
	if err != nil {
		zap.L().Error("DB错误", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.MobileExisted)
	}

	//执行业务，将数据存入数据库：member和组织表
	//pwd := encrypts.Md5(msg.Password)
	fmt.Println("grpc调用方传的password:", msg.Password)
	mem := &member.Member{
		Account:       msg.Name,
		Password:      msg.Password,
		Name:          msg.Name,
		Mobile:        msg.Mobile,
		Email:         msg.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        1, //model.Normal
	}
	err = lg.transaction.Action(func(conn database.DbConn) error { //事务：保障操作的原子性
		if err := lg.memberRepo.SaveMember(conn, c, mem); err != nil {
			fmt.Println("出现db错误")
			zap.L().Error("register save member db error", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}

		org := &organization.Organization{
			Name:       mem.Name + "个人组织",
			MemberId:   mem.Id,
			CreateTime: time.Now().UnixMilli(),
			Personal:   1, //model.Personal
			Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		}
		//业务：创建一个人的账号再自动给他创建个人组织
		err = lg.organizationRepo.SaveOrganization(conn, c, org)
		if err != nil {
			zap.L().Error("register SaveOrganization db err", zap.Error(err))
			return errs.GrpcError(model.DBError) //数据库错误
		}

		return nil
	})

	return &user_grpc.RegisterResponse{}, err
}

func (lg *LoginService) Login(ctx context.Context, msg *user_grpc.LoginMessage) (*user_grpc.LoginResponse, error) { //实现微服务login
	c := context.Background()
	//msg.Password = encrypts.Md5(msg.Password)
	fmt.Println("Login/rpc传过来的密码为", msg.Password)
	memFormDB, err := lg.memberRepo.FindMember(c, msg.Account, msg.Password) //因为注册前端已经做了加密
	if err != nil {
		fmt.Println(err)
		zap.L().Error("查询登录用户时数据库出现异常", zap.Error(err)) //查询用户失败!!!!!
		return nil, errs.GrpcError(model.DBError)
	}
	if memFormDB == nil { //注意为什么要接收判断这2个错误
		return nil, errs.GrpcError(model.AccuntAndPwdError)
	}
	memMsg := &user_grpc.MemberMessage{}
	err = copier.Copy(memMsg, memFormDB)
	memMsg.Code, _ = encrypts.EncryptInt64(memFormDB.Id, model.AESkey)
	memMsg.LastLoginTime = common.FormatByMill(memFormDB.LastLoginTime)
	memMsg.CreateTime = common.FormatByMill(memFormDB.CreateTime)

	orgs, err := lg.organizationRepo.FindOrganizationByMemId(c, memFormDB.Id)
	if err != nil {
		zap.L().Error("查询登录用户的组织失败", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var orgMessage []*user_grpc.OrganizationMessage
	err = copier.Copy(&orgMessage, orgs)
	orgMap := organization.ToMap(orgs)
	for _, v := range orgMessage { //手动填入code,ownerCode和时间
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESkey)
		v.OwnerCode = memMsg.Code
		v.CreateTime = common.FormatByMill(orgMap[v.Id].CreateTime) //获取原始时间

	}
	if len(orgs) > 0 {
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgs[0].Id, encrypts.AESKEY) //加密第一个组织的id?
	}

	exp := time.Duration(config.C.JwtConfig.AccessExp) * 3600 * 24 * time.Second //与time.duration相乘需要匹配类型
	rExp := time.Duration(config.C.JwtConfig.RefreshExp) * 3600 * 24 * time.Second
	idStr := strconv.FormatInt(memFormDB.Id, 10)

	token := jwts.CreateToken(idStr, exp, config.C.JwtConfig.AccessSecret, rExp, config.C.JwtConfig.RefreshSecret) //负载中只存了id

	tokenList := &user_grpc.TokenMessage{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		AccessTokenExp: token.AccessExp,
		TokenType:      "bearer",
	}
	//token:=jwts.CreateToken(string(memFormDB.Id),)
	resp := &user_grpc.LoginResponse{
		Member:           memMsg,
		OrganizationList: orgMessage,
		TokenList:        tokenList,
	}
	return resp, nil
}

func (lg *LoginService) TokenVerify(ctx context.Context, msg *user_grpc.LoginMessage) (*user_grpc.LoginResponse, error) {
	tokenStr := msg.Token
	//fmt.Println("grpc通信获取的token为:", tokenStr)

	parseToken, err := jwts.ParseToken(tokenStr, config.C.JwtConfig.AccessSecret)
	if err != nil {
		zap.L().Error("Token验证失败:", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}
	//数据库查询
	//fmt.Println(parseToken)
	id, _ := strconv.ParseInt(parseToken, 10, 64)
	memberbyId, err := lg.memberRepo.FindMemberById(ctx, id)
	if err != nil {
		zap.L().Error("查询id用户失败，数据库错误", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	memMsg := &user_grpc.MemberMessage{}
	copier.Copy(memMsg, memberbyId)
	memMsg.Code, _ = encrypts.EncryptInt64(memberbyId.Id, model.AESkey)

	orgsDB, err := lg.organizationRepo.FindOrganizationByMemId(ctx, id) //【tokenVerify中间件中，用户，组织信息每次都在数据库中查询？先思考一个问题，用户信息不能直接放在jwt的载荷中吗？如果不能，只能通过载荷中的id到数据库中去查询，这样可能开销有一些大，因为每个接口都会走一次verify的中间件。加缓存？？】
	if err != nil {
		zap.L().Error("查询登录用户的组织失败", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}

	if len(orgsDB) > 0 { //只需要加密组织字段
		memMsg.OrganizationCode, _ = encrypts.EncryptInt64(orgsDB[0].Id, encrypts.AESKEY) //加密第一个组织的id?
	}

	return &user_grpc.LoginResponse{Member: memMsg}, nil
}

func (lg *LoginService) MyOrganization(ctx context.Context, in *user_grpc.MyOrgReqGrpc) (*user_grpc.MyOrgRspGrpc, error) {
	orgList, err := lg.organizationRepo.FindOrganizationByMemId(ctx, in.Id)
	if err != nil {
		zap.L().Error("查询组织列表失败", zap.Error(err))

		return &user_grpc.MyOrgRspGrpc{OrganizationList: nil}, err
	}

	ogrs := []*user_grpc.OrganizationMessage{}
	copier.Copy(&ogrs, orgList)
	return &user_grpc.MyOrgRspGrpc{OrganizationList: ogrs}, nil
}
