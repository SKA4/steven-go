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

package steven

import (
	"encoding/base64"
	"encoding/json"
	"sort"
	"strings"

	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/ui"
)

var playerList = map[protocol.UUID]*playerInfo{}

type playerInfo struct {
	name        string
	uuid        protocol.UUID
	displayName chat.AnyComponent
	gameMode    gameMode
	ping        int

	skin     *render.TextureInfo
	skinHash string
}

type playerListUI struct {
	enabled bool

	background *ui.Image
	elements   []*ui.Text
	icons      [][2]*ui.Image
}

func (p *playerListUI) init() {
	p.background = ui.NewImage(render.GetTexture("solid"), 0, 16, 16, 16, 0, 0, 1, 1, 0, 0, 0)
	p.background.A = 120
	p.background.Visible = false
	ui.AddDrawable(p.background, ui.Top, ui.Center)
}

func (p *playerListUI) set(enabled bool) {
	p.enabled = enabled
	p.background.Visible = enabled
	for _, e := range p.elements {
		e.Visible = enabled
	}
	for _, i := range p.icons {
		i[0].Visible = enabled
		i[1].Visible = enabled
	}
}

func (p *playerListUI) render(delta float64) {
	if !p.enabled {
		return
	}
	for _, e := range p.elements {
		e.Visible = false
	}
	for _, i := range p.icons {
		i[0].Visible = false
		i[1].Visible = false
	}
	offset := 0
	count := 0
	width := 0.0
	for i, pl := range p.players() {
		if offset >= len(p.elements) {
			text := ui.NewText("", 8, 0, 255, 255, 255)
			p.elements = append(p.elements, text)
			ui.AddDrawable(text, ui.Top, ui.Center)
			icon := ui.NewImage(pl.skin, 0, 0, 16, 16, 8/64.0, 8/64.0, 8/64.0, 8/64.0, 255, 255, 255)
			ui.AddDrawable(icon, ui.Top, ui.Center)
			iconHat := ui.NewImage(pl.skin, 0, 0, 16, 16, 40/64.0, 8/64.0, 8/64.0, 8/64.0, 255, 255, 255)
			ui.AddDrawable(iconHat, ui.Top, ui.Center)
			p.icons = append(p.icons, [2]*ui.Image{icon, iconHat})
			text.Parent = p.background
			icon.Parent = p.background
			iconHat.Parent = p.background
		}
		text := p.elements[offset]
		icons := p.icons[offset]
		offset++
		text.Visible = true
		text.Y = 1 + 18*float64(i)
		text.Update(pl.name)
		count++
		if text.Width > width {
			width = text.Width
		}
		for _, ic := range icons {
			ic.Visible = true
			ic.Y = 1 + 18*float64(i)
			ic.Texture = pl.skin
		}
	}
	for _, i := range p.icons {
		if i[0].Visible {
			i[0].X = -width/2 - 4
			i[1].X = -width/2 - 4
		}
	}

	p.background.W = width + 32
	p.background.H = float64(count * 18)
}

func (p *playerListUI) players() (out []*playerInfo) {
	for _, pl := range playerList {
		out = append(out, pl)
	}
	sort.Sort(sortedPlayerList(out))
	return
}

type sortedPlayerList []*playerInfo

func (s sortedPlayerList) Len() int { return len(s) }
func (s sortedPlayerList) Less(a, b int) bool {
	if s[a].name < s[b].name {
		return true
	}
	return false
}
func (s sortedPlayerList) Swap(a, b int) { s[a], s[b] = s[b], s[a] }

func (handler) PlayerListInfo(p *protocol.PlayerInfo) {
	for _, pl := range p.Players {
		if _, ok := playerList[pl.UUID]; (!ok && p.Action != 0) || (ok && p.Action == 0) {
			continue
		}
		switch p.Action {
		case 0: // Add
			i := &playerInfo{
				name:        pl.Name,
				uuid:        pl.UUID,
				displayName: pl.DisplayName,
				gameMode:    gameMode(pl.GameMode),
				ping:        int(pl.Ping),
			}
			for _, prop := range pl.Properties {
				if prop.Name == "textures" {
					data, err := base64.URLEncoding.DecodeString(prop.Value)
					if err != nil {
						continue
					}
					var blob skinBlob
					err = json.Unmarshal(data, &blob)
					if err != nil {
						continue
					}
					url := blob.Textures.Skin.Url
					if strings.HasPrefix(url, "http://textures.minecraft.net/texture/") {
						i.skinHash = url[len("http://textures.minecraft.net/texture/"):]
						render.RefSkin(i.skinHash)
						i.skin = render.Skin(i.skinHash)
					}
				}
			}
			if i.skin == nil {
				i.skin = render.GetTexture("entity/steve")
			}
			playerList[pl.UUID] = i
		case 1: // Update gamemode
			playerList[pl.UUID].gameMode = gameMode(pl.GameMode)
		case 2: // Update ping
			playerList[pl.UUID].ping = int(pl.Ping)
		case 3: // Update display name
			playerList[pl.UUID].displayName = pl.DisplayName
		case 4: // Remove
			i := playerList[pl.UUID]
			if i.skinHash != "" {
				render.FreeSkin(i.skinHash)
			}
			delete(playerList, pl.UUID)
		}
	}
}

type skinBlob struct {
	Timestamp     int64
	ProfileID     string
	ProfileString string
	IsPublic      bool
	Textures      struct {
		Skin skinPath
		Cape skinPath
	}
}

type skinPath struct {
	Url string
}
