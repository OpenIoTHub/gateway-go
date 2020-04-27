package client

import (
	"encoding/json"
	"fmt"
	_ "github.com/OpenIoTHub/gateway-go/component"
	"github.com/OpenIoTHub/gateway-go/config"
	"github.com/OpenIoTHub/gateway-go/services"
	"github.com/OpenIoTHub/utils/models"
	"github.com/gorilla/mux"
	"github.com/iotdevice/zeroconf"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"strconv"
)

var ConfigMode *models.GatewayConfig

func Run() {
	port, _ := strconv.Atoi(config.Setting["apiPort"])
	//mDNS注册服务
	_, err := zeroconf.Register("OpenIoTHubGateway", "_openiothub-gateway._tcp", "local.", port, []string{}, nil)
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
	var err error
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
	ConfigMode.Server.ServerHost = r.FormValue("server_host")
	ConfigMode.Server.TcpPort, err = strconv.Atoi(r.FormValue("tcp_port"))
	ConfigMode.Server.UdpApiPort, err = strconv.Atoi(r.FormValue("udp_p2p_port"))
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
	ConfigMode.Server.LoginKey = r.FormValue("login_key")
	if ConfigMode.LastId == "" {
		ConfigMode.LastId = uuid.Must(uuid.NewV4()).String()
	}

	clientToken, err := models.GetToken(*ConfigMode, 1, 200000000000)
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
	err = services.RunNATManager(ConfigMode.Server.LoginKey, clientToken)
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
	config.Setting["explorerToken"], err = models.GetToken(*ConfigMode, 2, 200000000000)
	err = config.WriteConfigFile(*ConfigMode, config.Setting["configFilePath"])
	if err != nil {
		log.Println(err.Error())
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
