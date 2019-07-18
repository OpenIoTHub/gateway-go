package client

import (
	"encoding/json"
	"fmt"
	"git.iotserv.com/iotserv/client/services"
	"git.iotserv.com/iotserv/utils/crypto"
	"git.iotserv.com/iotserv/utils/models"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

var loged = false
var configMode *models.ClientFlat

func init() {
	configMode = &models.ClientFlat{
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
	r := mux.NewRouter()
	r.HandleFunc("/", getExplorerToken)
	r.HandleFunc("/getExplorerToken", getExplorerToken)
	r.HandleFunc("/loginServer", loginServer)
	http.Handle("/", r)
	err := http.ListenAndServe("127.0.0.1:1082", nil) //设置监听的端口
	if err != nil {
		fmt.Printf("请检查端口1082是否被占用")
	}
}

func loginServer(w http.ResponseWriter, r *http.Request) {
	if loged {
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
	//err := json.Unmarshal(body, &configMode)
	r.ParseForm()
	configMode.LastId = r.FormValue("last_id")
	configMode.ServerHost = r.FormValue("server_host")
	configMode.TcpPort = r.FormValue("tcp_port")
	configMode.UdpApiPort = r.FormValue("udp_p2p_port")
	configMode.LoginKey = r.FormValue("login_key")
	if configMode.LastId == "" {
		configMode.LastId = uuid.Must(uuid.NewV4()).String()
	}
	tcpP, err := strconv.Atoi(configMode.TcpPort)
	kcpP, err := strconv.Atoi(configMode.KcpPort)
	tlsP, err := strconv.Atoi(configMode.TlsPort)
	udpApiP, err := strconv.Atoi(configMode.UdpApiPort)
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
	clientToken, err := crypto.GetToken(configMode.LoginKey, configMode.LastId, configMode.ServerHost, tcpP,
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
	err = services.RunNATManager(configMode.LoginKey, clientToken)
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
	loged = true
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
	tcpP, err := strconv.Atoi(configMode.TcpPort)
	kcpP, err := strconv.Atoi(configMode.KcpPort)
	tlsP, err := strconv.Atoi(configMode.TlsPort)
	udpApiP, err := strconv.Atoi(configMode.UdpApiPort)
	explorerToken, err := crypto.GetToken(configMode.LoginKey, configMode.LastId, configMode.ServerHost, tcpP,
		kcpP, tlsP, udpApiP, 2, 200000000000)
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
	response := Response{
		Code:  0,
		Msg:   "success",
		Token: explorerToken,
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
