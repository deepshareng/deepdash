package email

const (
	HOST        = "smtp.ym.163.com"
	SERVER_ADDR = "smtp.ym.163.com:25"
	//发送邮件的邮箱
	USER = "register@deepshare.io"
	//USER        = "deepshare@misingularity.com"

	PASSWORD = "dsqdjzds"        //发送邮件邮箱的授权码
	RECEIVER = "bd@deepshare.io" //接收邮件的邮箱
	//RECEIVER = `deepshare@misingularity.io,johney.song@misingularity.io,zhaojian@misingularity.io,huiyan@misingularity.io,teng.zhang@misingularity.io,nanxi.li@misingularity.io,asang@misingularity.io` //接收邮件的邮箱
)

type Email struct {
	to      string
	subject string
	msg     string
	format  string
}

type HostMail struct {
	Identity, User, Password, Host string
}
