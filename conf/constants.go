// Copyright © 2024 Maxr1998 <max@maxr1998.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package conf

import "crypto/aes"

const (
	ServiceName  = "s3share"
	AppName      = ServiceName + "-cli"
	FileIdLength = 10            // 80 bits
	KeyLength    = 32            // AES-256
	IvLength     = aes.BlockSize // 16 bytes
	IvSize       = 8             // actual IV size (8 bytes) for AES-CTR without the counter
)

var ConfigFileLocations = []string{
	"$XDG_CONFIG_HOME/" + ServiceName,
	"$HOME/.config/" + ServiceName,
	"$HOME/." + ServiceName,
	".",
}