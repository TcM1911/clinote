/*
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (C) Joakim Kennedy, 2016
 */

package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallbackHandler(t *testing.T) {
	assert := assert.New(t)
	tempToken := "internal-dev.14CD91FCE1F.687474703A2F2F6C6F63616C686F7374.6E287AD298969B6F8C0B4B1D67BCAB1D"
	verifier := "40793F8BAE15D4E3B6DD5CA8AB4BF62F"
	sandbox := "false"

	c := make(chan *callbackValues)
	url := fmt.Sprintf("http://www.sample.com/?oauth_token=%s&&oauth_verifier=%s&&sandbox_lnb=%s", tempToken, verifier, sandbox)
	r := httptest.NewRequest(http.MethodGet, url, nil)
	go oathCallbackHandler(c).ServeHTTP(nil, r)
	vals := <-c

	assert.Equal(verifier, vals.Verifier)
	assert.Equal(tempToken, vals.TempToken)
	assert.False(vals.SandboxLnb)
}
