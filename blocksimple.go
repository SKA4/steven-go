// Copyright 2015 Matthew Collins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

type blockSimple struct {
	baseBlock
}

type simpleConfig struct {
	NotCullAgainst bool
	NoCollision    bool
	NotRendable    bool
}

func initSimple(name string) *BlockSet {
	return initSimpleConfig(name, simpleConfig{})
}

func initSimpleConfig(name string, config simpleConfig) *BlockSet {
	s := &blockSimple{}
	s.init(name)
	set := alloc(s)

	s.cullAgainst = !config.NotCullAgainst
	s.collidable = !config.NoCollision
	s.renderable = !config.NotRendable

	return set
}

func (b *blockSimple) toData() int {
	if b == b.Parent.Base {
		return 0
	}
	return -1
}
