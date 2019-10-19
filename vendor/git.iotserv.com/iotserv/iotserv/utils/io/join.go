package io

import (
	_ "git.iotserv.com/iotserv/iotserv/utils/pool"
	"io"
	_ "sync"
)

//func Join(c1 io.ReadWriteCloser, c2 io.ReadWriteCloser) (inCount int64, outCount int64) {
//	var wait sync.WaitGroup
//	pipe := func(to io.ReadWriteCloser, from io.ReadWriteCloser, count *int64) {
//		defer to.Close()
//		defer from.Close()
//		defer wait.Done()
//
//		buf := pool.GetBuf(16 * 1024)
//		defer pool.PutBuf(buf)
//		*count, _ = io.CopyBuffer(to, from, buf)
//	}
//
//	wait.Add(2)
//	go pipe(c1, c2, &inCount)
//	go pipe(c2, c1, &outCount)
//	wait.Wait()
//	return
//}

//func Join(c1 io.ReadWriteCloser, c2 io.ReadWriteCloser) (inCount int64, outCount int64) {
//	var wait sync.WaitGroup
//	pipe := func(to io.ReadWriteCloser, from io.ReadWriteCloser, count *int64) {
//		defer to.Close()
//		defer from.Close()
//		defer wait.Done()
//		*count, _ = io.Copy(to, from)
//	}
//
//	wait.Add(2)
//	go pipe(c1, c2, &inCount)
//	go pipe(c2, c1, &outCount)
//	wait.Wait()
//	return
//}

func Join(p1 io.ReadWriteCloser, p2 io.ReadWriteCloser) (inCount int64, outCount int64) {
	defer p1.Close()
	defer p2.Close()
	// start tunnel
	p1die := make(chan struct{})
	go func() { io.Copy(p1, p2); close(p1die) }()
	p2die := make(chan struct{})
	go func() { io.Copy(p2, p1); close(p2die) }()
	// wait for tunnel termination
	select {
	case <-p1die:
	case <-p2die:
	}
	return 0, 0
}
