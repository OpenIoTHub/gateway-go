package connect

import (
	"golang.org/x/crypto/ssh"
	"net"
	"strconv"
)

//func JoinSSHold(stream *mux.Stream,targetIP string,targetPort int,userName,passWord string) error  {
//	sshConn,err := ConnectToSSH(targetIP,targetPort,userName,passWord)
//	if err!=nil{
//		return err
//	}
//	go io.Join(stream,sshConn)
//	return nil
//}

func JoinSSH(stream net.Conn, remoteIP string, remotePort int, userName, passWord string) (err error) {
	client, err := ssh.Dial("tcp", net.JoinHostPort(remoteIP, strconv.Itoa(remotePort)), &ssh.ClientConfig{
		User:            userName,
		Auth:            []ssh.AuthMethod{ssh.Password(passWord)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stdout = stream
	session.Stderr = stream
	session.Stdin = stream
	//go func() {
	//	time.Sleep(time.Second)
	//	Join(session.Stdin,session.Stdout,stream)
	//}()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", 25, 80, modes)
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}
	session.Wait()
	return
}
