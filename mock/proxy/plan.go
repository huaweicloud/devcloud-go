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
 */

package proxy

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// ErrNotFound is returned iff SelectRule can't find a Rule that applies
	ErrNotFound = errors.New("no matching rule found")
)

func marshalCommand(cmd string) []byte {
	return []byte(cmd)
}

type Plan struct {
	rulesMap sync.Map
}

func (p *Plan) AddRule(r Rule) error {
	if r.Percentage < 0 || r.Percentage > 100 {
		return fmt.Errorf("percentage in rule #%s is malformed. it must within 0-100", r.Name)
	}
	if len(r.Name) <= 0 {
		return fmt.Errorf("name of rule is required")
	}
	if len(r.Command) > 0 {
		r.marshaledCmd = marshalCommand(r.Command)
	}
	if _, ok := p.rulesMap.Load(r.Name); ok {
		return fmt.Errorf("a rule by the same name exists")
	}
	p.rulesMap.Store(r.Name, &r)
	return nil
}

func (p *Plan) SelectRule(clientAddr string, buf []byte) *Rule {
	var chosenRule *Rule
	p.rulesMap.Range(func(key, value interface{}) bool {
		rule, ok := value.(*Rule)
		if !ok {
			return true
		}
		if len(rule.ClientAddr) > 0 && strings.HasPrefix(clientAddr, rule.ClientAddr) {
			return true
		}

		if len(rule.Command) > 0 && !bytes.Contains(buf, rule.marshaledCmd) {
			return true
		}

		chosenRule = rule
		return false
	})
	if chosenRule == nil {
		return nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if chosenRule.Percentage > 0 && r.Intn(100) > chosenRule.Percentage {
		return nil
	}
	atomic.AddUint64(&chosenRule.hits, 1)
	return chosenRule
}

func (p *Plan) DeleteRule(name string) error {
	_, ok := p.rulesMap.Load(name)
	if !ok {
		return ErrNotFound
	}
	p.rulesMap.Delete(name)
	return nil
}

func (p *Plan) DeleteAllRule() {
	p.rulesMap.Range(func(key, value interface{}) bool {
		p.rulesMap.Delete(key)
		return true
	})
}
