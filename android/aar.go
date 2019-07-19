package client

import (
	"encoding/json"
	"fmt"
	"git.iotserv.com/iotserv/client/config"
	"git.iotserv.com/iotserv/client/services"
	"git.iotserv.com/iotserv/utils/crypto"
	"git.iotserv.com/iotserv/utils/models"
	"github.com/gorilla/mux"
	"github.com/iotdevice/zeroconf"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

var ConfigMode *models.ClientFlat

func init() {
	ConfigMode = &models.ClientFlat{
		1082,
		uuid.Must(uuid.NewV4()).String(),
		"tcp",
		"guonei.nat-cloud.com",
		"34320",
		"34320",
		"34321",
		"34321",
		"HLLdsa544&*S",
	}
}

func Run() {
	port, _ := strconv.Atoi(config.Setting["apiPort"])
	//mDNS注册服务
	_, err := zeroconf.Register("nat-cloud-client", "_nat-cloud-client._tcp", "local.", port, []string{}, nil)
	//
	r := mux.NewRouter()
	r.HandleFunc("/", getExplorerToken)
	r.HandleFunc("/getExplorerToken", getExplorerToken)
	r.HandleFunc("/loginServer", loginServer)
	http.Handle("/", r)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", config.Setting["apiPort"]), nil) //设置监听的端口
	if err != nil {
		fmt.Printf("请检查端口%s是否被占用", config.Setting["apiPort"])
	}
}

func loginServer(w http.ResponseWriter, r *http.Request) {
	if config.Loged {
		response := Response{
			Code: 1,
			Msg:  "Already logged in",
		}
		responseJson, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJson)
		return
	}
	//body, _ := ioutil.ReadAll(r.Body)
	//fmt.Errorf(string(body))
	//err := json.Unmarshal(body, &ConfigMode)
	r.ParseForm()
	ConfigMode.LastId = r.FormValue("last_id")
	ConfigMode.ServerHost = r.FormValue("server_host")
	ConfigMode.TcpPort = r.FormValue("tcp_port")
	ConfigMode.UdpApiPort = r.FormValue("udp_p2p_port")
	ConfigMode.LoginKey = r.FormValue("login_key")
	if ConfigMode.LastId == "" {
		ConfigMode.LastId = uuid.Must(uuid.NewV4()).String()
	}
	tcpP, err := strconv.Atoi(ConfigMode.TcpPort)
	kcpP, err := strconv.Atoi(ConfigMode.KcpPort)
	tlsP, err := strconv.Atoi(ConfigMode.TlsPort)
	udpApiP, err := strconv.Atoi(ConfigMode.UdpApiPort)
	if err != nil {
		response := Response{
			Code: 1,
			Msg:  err.Error(),
		}
		responseJson, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJson)
		return
	}
	clientToken, err := crypto.GetToken(ConfigMode.LoginKey, ConfigMode.LastId, ConfigMode.ServerHost, tcpP,
		kcpP, tlsP, udpApiP, 1, 200000000000)
	if err != nil {
		response := Response{
			Code: 1,
			Msg:  err.Error(),
		}
		responseJson, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJson)
		return
	}
	err = services.RunNATManager(ConfigMode.LoginKey, clientToken)
	if err != nil {
		response := Response{
			Code: 1,
			Msg:  err.Error(),
		}
		responseJson, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJson)
		return
	}
	config.Loged = true
	config.Setting["explorerToken"], err = crypto.GetToken(ConfigMode.LoginKey, ConfigMode.LastId, ConfigMode.ServerHost, tcpP,
		kcpP, tlsP, udpApiP, 2, 200000000000)
	err = config.WriteConfigFile(models.ClientConfig{
		ExplorerTokenHttpPort: ConfigMode.ExplorerTokenHttpPort,
		Server: models.Srever{
			ConnectionType: ConfigMode.ConnectionType,
			ServerHost:     ConfigMode.ServerHost,
			TcpPort:        tcpP,
			KcpPort:        kcpP,
			UdpApiPort:     udpApiP,
			TlsPort:        tlsP,
			LoginKey:       ConfigMode.LoginKey,
		},
		LastId: ConfigMode.LastId,
	}, config.Setting["configFilePath"])
	if err != nil {
		fmt.Println(err.Error())
	}
	response := Response{
		Code: 0,
		Msg:  "success",
	}
	responseJson, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
	return
}

func getExplorerToken(w http.ResponseWriter, r *http.Request) {
	if config.Loged != true {
		response := Response{
			Code: 1,
			Msg:  "你还没有登录",
		}
		responseJson, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJson)
		return
	}
	response := Response{
		Code:  0,
		Msg:   "success",
		Token: config.Setting["explorerToken"],
	}
	responseJson, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
	return
}

type Response struct {
	Code  int
	Msg   string
	Token string
}
