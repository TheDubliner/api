// Vikunja is a todo-list application to facilitate your life.
// Copyright 2018-2020 Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package files

import (
	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/testfixtures.v2"
	"os"
	"path/filepath"
	"testing"
)

// This file handles storing and retrieving a file for different backends
var fs afero.Fs
var afs *afero.Afero

// InitFileHandler creates a new file handler for the file backend we want to use
func InitFileHandler() {
	fs = afero.NewOsFs()
	afs = &afero.Afero{Fs: fs}
}

// InitTestFileHandler initializes a new memory file system for testing
func InitTestFileHandler() {
	fs = afero.NewMemMapFs()
	afs = &afero.Afero{Fs: fs}
}

func initFixtures(t *testing.T) {
	// Init db fixtures
	err := db.LoadFixtures()
	assert.NoError(t, err)

	InitTestFileFixtures(t)
}

//InitTestFileFixtures initializes file fixtures
func InitTestFileFixtures(t *testing.T) {
	// Init fixture files
	filename := config.FilesBasePath.GetString() + "/1"
	err := afero.WriteFile(afs, filename, []byte("testfile1"), 0644)
	assert.NoError(t, err)
}

// InitTests handles the actual bootstrapping of the test env
func InitTests() {
	var err error
	x, err = db.CreateTestEngine()
	if err != nil {
		log.Fatal(err)
	}

	err = x.Sync2(GetTables()...)
	if err != nil {
		log.Fatal(err)
	}

	config.InitDefaultConfig()
	// We need to set the root path even if we're not using the config, otherwise fixtures are not loaded correctly
	config.ServiceRootpath.Set(os.Getenv("VIKUNJA_SERVICE_ROOTPATH"))

	// Sync fixtures
	var fixturesHelper testfixtures.Helper = &testfixtures.SQLite{}
	if config.DatabaseType.GetString() == "mysql" {
		fixturesHelper = &testfixtures.MySQL{}
	}
	fixturesDir := filepath.Join(config.ServiceRootpath.GetString(), "pkg", "files", "fixtures")
	err = db.InitFixtures(fixturesHelper, fixturesDir)
	if err != nil {
		log.Fatal(err)
	}

	InitTestFileHandler()
}

// FileStat stats a file. This is an exported function to be able to test this from outide of the package
func FileStat(filename string) (os.FileInfo, error) {
	return afs.Stat(filename)
}
