package worker

import (
	"time"

	"github.com/forgoer/openssl"
	"github.com/open-tdp/go-helper/logman"
	"github.com/open-tdp/go-helper/psutil"

	"tdp-cloud/cmd/args"
)

func (pod *SendPod) Register() (uint, error) {

	logman.Info("register:send")

	stat := psutil.Summary(true)
	cloudId := psutil.CloudInstanceId()
	workerId := openssl.Md5ToString(stat.HostId)

	data := &map[string]any{
		"CloudId":       cloudId,
		"WorkerId":      workerId,
		"WorkerMeta":    stat,
		"WorkerVersion": args.Version,
	}

	err := pod.WriteJson(&SocketData{
		Method:  "Register",
		TaskId:  0,
		Payload: data,
	})

	return 0, err

}

func (pod *RespPod) Register(rs *SocketData) {

	logman.Info("register:resp", "payload", rs.Payload)

	go KeepAlive(&SendPod{pod.WsConn})

}

//// 持续报送状态

func KeepAlive(pod *SendPod) error {

	for {
		time.Sleep(25 * time.Second)

		if _, err := pod.Ping(); err != nil {
			logman.Error("ping:fail", "error", err)
			return err
		}
	}

}
