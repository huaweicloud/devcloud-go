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
 * 2022.01.20-Adapt to Redis and MySQL fault injection, delete unnecessary functions
 * add required logic.
 * 			Huawei Technologies Co., Ltd.
 */

package proxy

import (
	"fmt"
	"strings"
	"sync/atomic"
)

type Rule struct {
	Name        string `json:"name,omiempty"`
	Delay       int    `json:"delay,omitempty"`
	Jitter      int    `json:"jitter,omitempty"`
	Drop        bool   `json:"drop,omitempty"`
	ReturnEmpty bool   `json:"return_empty,omitempty"`
	ReturnErr   error  `json:"return_err,omitempty"`
	Percentage  int    `json:"percentage,omitempty"`
	// SelectRule does prefix matching on this value
	ClientAddr string `json:"client_addr,omitempty"`
	Command    string `json:"command,omitempty"`
	// filled by marshalCommand
	marshaledCmd []byte
	hits         uint64
}

// setPCC set percentage clientAddr command
func (r *Rule) setPCC(percentage int, clientAddr, command string) {
	if percentage >= 0 && percentage < 100 {
		r.Percentage = percentage
	} else {
		r.Percentage = 0
	}
	if clientAddr != "" {
		r.ClientAddr = clientAddr
	}
	if command != "" {
		r.Command = command
	}
}

func (r *Rule) String() string {
	buf := make([]string, 0)
	buf = append(buf, r.Name)

	// count hits
	hits := atomic.LoadUint64(&r.hits)
	buf = append(buf, fmt.Sprintf("hits=%d", hits))

	if r.Delay > 0 {
		buf = append(buf, fmt.Sprintf("delay=%d", r.Delay))
	}
	if r.Jitter > 0 {
		buf = append(buf, fmt.Sprintf("jitter=%d", r.Jitter))
	}
	if r.Drop {
		buf = append(buf, fmt.Sprintf("drop=%t", r.Drop))
	}
	if r.ReturnEmpty {
		buf = append(buf, fmt.Sprintf("return_empty=%t", r.ReturnEmpty))
	}
	if r.ReturnErr != nil {
		buf = append(buf, fmt.Sprintf("return_err=%s", r.ReturnErr))
	}
	if len(r.ClientAddr) > 0 {
		buf = append(buf, fmt.Sprintf("client_addr=%s", r.ClientAddr))
	}
	if len(r.Command) > 0 {
		buf = append(buf, fmt.Sprintf("command=%s", r.Command))
	}
	if r.Percentage > 0 {
		buf = append(buf, fmt.Sprintf("percentage=%d", r.Percentage))
	}

	return strings.Join(buf, " ")
}
