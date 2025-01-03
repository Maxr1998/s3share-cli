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

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitConfig(failOnError bool) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	for _, location := range ConfigFileLocations {
		viper.AddConfigPath(location)
	}

	if err := viper.ReadInConfig(); err != nil && failOnError {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Fprintln(os.Stderr, "Config file not found. Please create a config file `config.toml` at one of the following locations:")
			for _, location := range ConfigFileLocations {
				fmt.Fprintln(os.Stderr, " -", location)
			}
			os.Exit(1)
		}
	}
}
