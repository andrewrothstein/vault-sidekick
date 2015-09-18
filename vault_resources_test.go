/*
Copyright 2015 Home Office All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestSetResources(t *testing.T) {
	var items vaultResources

	assert.Nil(t, items.Set("secret:test:fn=filename.test,fmt=yaml"))
	assert.Nil(t, items.Set("secret:test:fn=filename.test,"))
	assert.Nil(t, items.Set("secret:/db/prod/username"))
	assert.Nil(t, items.Set("secret:/db/prod:fn=filename.test,fmt=yaml"))
	assert.Nil(t, items.Set("secret:test:fn=filename.test,"))
	assert.Nil(t, items.Set("pki:example-dot-com:cn=blah.example.com"))
	assert.Nil(t, items.Set("pki:example-dot-com:cn=blah.example.com,fn=/etc/certs/ssl/blah.example.com"))
	assert.Nil(t, items.Set("pki:example-dot-com:cn=blah.example.com,rn=10s"))
	assert.NotNil(t, items.Set("secret:"))
	assert.NotNil(t, items.Set("secret:test:fn=filename.test,fmt="))
	assert.NotNil(t, items.Set("secret::fn=filename.test,fmt=yaml"))
	assert.NotNil(t, items.Set("secret:te1st:fn=filename.test,fmt="))
	assert.NotNil(t, items.Set("fn=filename.test,fmt=yaml"))
}

func TestResourceSize(t *testing.T) {
	var items vaultResources
	items.Set("secret:test:fn=filename.test,fmt=yaml")
	items.Set("secret:test:fn=fileame.test")
	assert.Equal(t, 2, items.size())
}

func TestResources(t *testing.T) {
	var items vaultResources
	items.Set("secret:test:fn=filename.test,fmt=yaml")
	items.Set("secret:test:fn=fileame.test")

	if passed := assert.Equal(t, len(items.items), 2); !passed {
		t.FailNow()
	}

	rn := items.items[0]
	assert.Equal(t, "secret", rn.resource)
	assert.Equal(t, "test", rn.name)
	assert.Equal(t, 2, len(rn.options))
	assert.Equal(t, "filename.test", rn.options[OptionFilename])
	assert.Equal(t, "yaml", rn.options[OptionFormat])
	rn = items.items[1]
	assert.Equal(t, "secret", rn.resource)
	assert.Equal(t, "test", rn.name)
	assert.Equal(t, 1, len(rn.options))
	assert.Equal(t, "fileame.test", rn.options[OptionFilename])
}