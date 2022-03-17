/*
 * Copyright (c) 2017, redfi
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * * Neither the name of the copyright holder nor the names of its
 *   contributors may be used to endorse or promote products derived from
 *   this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * 202.01.20-Adapt to Redis and MySQL fault injection, delete unnecessary functions
 * add required logic.
 * 			Huawei Technologies Co., Ltd.
 *
 */

// Package proxy TCP-based fault injection
package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/dolthub/vitess/go/mysql"
	"github.com/go-redis/redis/v8"
	"github.com/tidwall/redcon"
	"gopkg.in/fatih/pool.v2"
)

func factory(server string) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		return net.Dial("tcp", server)
	}
}

type MockType string

// agent type
const (
	Redis MockType = "redis"
	Mysql MockType = "mysql"
	Etcd  MockType = "etcd"
)

// Proxy real service agent
type Proxy struct {
	Server   string
	Addr     string
	plan     *Plan
	connMap  sync.Map
	connPool pool.Pool
	listener net.Listener
	mock     MockType
}

func NewProxy(server, addr string, mock MockType) *Proxy {
	plan := &Plan{}
	return &Proxy{
		Server: server,
		Addr:   addr,
		plan:   plan,
		mock:   mock,
	}
}

func (p *Proxy) StartProxy() error {
	var err error
	p.connPool, err = pool.NewChannelPool(5, 30, factory(p.Server))
	if err != nil {
		log.Println(err)
		return err
	}
	p.listener, err = net.Listen("tcp", p.Addr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("proxy [%s] started! ", p.Addr)
	go func() {
		for {
			conn, err := p.listener.Accept()
			if err != nil {
				log.Println(err)
				break
			}
			p.connMap.Store(conn.RemoteAddr().String(), &conn)
			go p.handle(conn)
		}
	}()
	return nil
}

// handle proxy service real service interactivity
func (p *Proxy) handle(conn net.Conn) {
	var wg sync.WaitGroup
	targetConn, err := p.connPool.Get()
	if err != nil {
		log.Fatal("failed to get a connection from connPool")
	}

	wg.Add(2)
	go func() {
		p.faulter(targetConn, conn)
		wg.Done()
	}()
	go func() {
		p.pipe(targetConn, conn)
		wg.Done()
	}()

	wg.Wait()
	p.connMap.Delete(conn.RemoteAddr().String())
	err = conn.Close()
	if err != nil {
		log.Println(err)
	}
}

// Write proxy service results
func (p *Proxy) faulter(dst, src net.Conn) {
	buf := make([]byte, 32<<10)
	for {
		n, err := src.Read(buf)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		rule := p.plan.SelectRule(src.RemoteAddr().String(), buf)
		if errflg := p.Write(src, rule); errflg {
			_, err = dst.Write(buf[:n])
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

// Write proxy service results write
func (p *Proxy) Write(src net.Conn, rule *Rule) bool {
	if rule == nil {
		return true
	}
	if rule.Delay > 0 {
		time.Sleep(time.Duration(rule.Delay) * time.Millisecond)
	}
	if rule.Jitter > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		time.Sleep(time.Duration(rule.Jitter*r.Intn(100)/100) * time.Millisecond)
	}
	if rule.Drop {
		if p.mock == Mysql {
			p.errWrite(src, mysql.NewSQLError(mysql.ERUnknownError, mysql.SSUnknownSQLState, "Server shutdown in progress"))
		} else {
			p.errWrite(src, RedisError("Server shutdown in progress"))
		}
		return false
	}
	if rule.ReturnEmpty {
		if p.mock == Mysql {
			p.errWrite(src, mysql.NewSQLError(mysql.ERUnknownError, mysql.SSUnknownSQLState, "nil"))
		} else {
			p.errWrite(src, redis.Nil)
		}
		return false
	}
	if rule.ReturnErr != nil {
		p.errWrite(src, rule.ReturnErr)
		return false
	}
	return true
}

// errWrite proxy service error results write
func (p *Proxy) errWrite(src net.Conn, srcErr error) {
	buf := make([]byte, 0)
	if p.mock == Mysql {
		sqlErr, err := srcErr.(*mysql.SQLError)
		log.Println(err)
		buf = writePacket(uint16(sqlErr.Num), sqlErr.State, "%v", sqlErr.Message)
	} else {
		redisErr, err := srcErr.(redis.Error)
		log.Println(err)
		buf = redcon.AppendError(buf, redisErr.Error())
	}
	_, err := src.Write(buf)
	if err != nil {
		log.Println(err)
	}
}

// pipe real service results
func (p *Proxy) pipe(dst, src net.Conn) {
	buf := make([]byte, 32<<10)

	for {
		n, err := dst.Read(buf)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		_, err = src.Write(buf[:n])
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (p *Proxy) StopProxy() {
	p.connMap.Range(func(key, value interface{}) bool {
		p.connMap.Delete(key)
		err := (*value.(*net.Conn)).Close()
		if err != nil {
			log.Println(err)
		}
		return true
	})
	err := p.listener.Close()
	if err != nil {
		log.Println(err)
	}
	p.connPool.Close()
	log.Printf("proxy [%s] stop! ", p.Addr)
}

func (p *Proxy) AddDelay(name string, delay, percentage int, clientAddr, command string) error {
	rule := Rule{
		Name:  name,
		Delay: delay,
	}
	rule.setPCC(percentage, clientAddr, command)
	return p.plan.AddRule(rule)
}

func (p *Proxy) AddJitter(name string, jitter, percentage int, clientAddr, command string) error {
	rule := Rule{
		Name:   name,
		Jitter: jitter,
	}
	rule.setPCC(percentage, clientAddr, command)
	return p.plan.AddRule(rule)
}

func (p *Proxy) AddDrop(name string, percentage int, clientAddr, command string) error {
	rule := Rule{
		Name: name,
		Drop: true,
	}
	rule.setPCC(percentage, clientAddr, command)
	return p.plan.AddRule(rule)
}

func (p *Proxy) AddReturnEmpty(name string, percentage int, clientAddr, command string) error {
	rule := Rule{
		Name:        name,
		ReturnEmpty: true,
	}
	if percentage > 0 && percentage < 100 {
		rule.Percentage = percentage
	}
	if clientAddr != "" {
		rule.ClientAddr = clientAddr
	}
	if command != "" {
		rule.Command = command
	}
	return p.plan.AddRule(rule)
}

func (p *Proxy) AddReturnErr(name string, returnErr error, percentage int, clientAddr, command string) error {
	rule := Rule{
		Name:      name,
		ReturnErr: returnErr,
	}
	rule.setPCC(percentage, clientAddr, command)
	return p.plan.AddRule(rule)
}

func (p *Proxy) DeleteAllRule() {
	p.plan.DeleteAllRule()
}

func writeByte(data []byte, pos int, value byte) int {
	data[pos] = value
	return pos + 1
}

func writeUint16(data []byte, pos int, value uint16) int {
	data[pos] = byte(value)
	data[pos+1] = byte(value >> 8)
	return pos + 2
}

func writeEOFString(data []byte, pos int, value string) int {
	pos += copy(data[pos:], value)
	return pos
}

func errorPacket(errorCode uint16, sqlState string, format string, args ...interface{}) ([]byte, error) {
	errorMessage := fmt.Sprintf(format, args...)
	length := 1 + 2 + 1 + 5 + len(errorMessage)

	data := make([]byte, length)
	pos := 0
	pos = writeByte(data, pos, mysql.ErrPacket)
	pos = writeUint16(data, pos, errorCode)
	pos = writeByte(data, pos, '#')
	if sqlState == "" {
		sqlState = mysql.SSUnknownSQLState
	}
	if len(sqlState) != 5 {
		return nil, errors.New("sqlState has to be 5 characters long")
	}
	pos = writeEOFString(data, pos, sqlState)
	_ = writeEOFString(data, pos, errorMessage)

	return data, nil
}

func writePacket(errorCode uint16, sqlState string, format string, args ...interface{}) []byte {
	buff := make([]byte, 0)
	buf, err := errorPacket(errorCode, sqlState, format, args)
	if err != nil {
		log.Println(err)
	}
	packetLength := len(buf)
	buff = append(buff, byte(packetLength))
	buff = append(buff, byte(packetLength>>8))
	buff = append(buff, byte(packetLength>>16))
	buff = append(buff, 1)
	buff = append(buff, buf...)
	return buff
}

// RedisError rediserror
type RedisError string

func (e RedisError) Error() string { return string(e) }

func (RedisError) RedisError() {}
