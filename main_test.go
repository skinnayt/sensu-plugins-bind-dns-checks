package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCheckArgs(t *testing.T) {
	// Arguments that should be checked
	// The bind DNS server will listen on an ip address and port for the statistics channel
	// The bind DNS server might be configured with XML or JSON for returning the statistics
	// THe bind DNS server can also dump the statistics to a file

	// Setup the assert object for running the tests
	assert := assert.New(t)

	// Make sure we can only pass in file, xml or json for the statistics format
	// When the format is file, we need to specify the file path
	plugin.StatisticsFormat = "file"
	plugin.StatisticsFilePath = "tests/named.stats"
	ok, err := checkArgs(nil)
	assert.Equal(ok, 0)
	assert.NoError(err)

	// Make sure we get an error when the format is not file, xml or json
	plugin.StatisticsFormat = "invalid"
	ok, err = checkArgs(nil)
	assert.Equal(ok, 3)
	assert.Error(err)

	// Make sure we get an error when format is file but not filepath set
	plugin.StatisticsFormat = "file"
	plugin.StatisticsFilePath = ""
	ok, err = checkArgs(nil)
	assert.Equal(ok, 3)
	assert.Error(err)

	// Make sure we get an error when format is file but filepath is invalid
	plugin.StatisticsFormat = "file"
	plugin.StatisticsFilePath = "invalid.file.path"
	ok, err = checkArgs(nil)
	assert.Equal(ok, 3)
	assert.Error(err)

	// Make sure we get an error when file is not readable
	os.Chmod("tests/unreadable.stats", 0100)
	plugin.StatisticsFormat = "file"
	plugin.StatisticsFilePath = "tests/unreadable.stats"
	ok, err = checkArgs(nil)
	assert.Equal(ok, 3)
	assert.Error(err)
	os.Chmod("tests/unreadable.stats", 0644)

	// For both XML and JSON, we don't need to specify a file path
	// But we need to specify the ip address and port
	plugin.StatisticsFormat = "xml"
	plugin.StatisticsFilePath = ""
	plugin.StatisticsIP = "127.0.0.1"
	plugin.StatisticsPort = 8053
	ok, err = checkArgs(nil)
	assert.Equal(ok, 0)
	assert.NoError(err)

	// Set bad ip address
	plugin.StatisticsFormat = "xml"
	plugin.StatisticsFilePath = ""
	plugin.StatisticsIP = "bad.ip.address"
	plugin.StatisticsPort = 8053
	ok, err = checkArgs(nil)
	assert.Equal(ok, 3)
	assert.Error(err)

	// Set missing ip address
	plugin.StatisticsFormat = "xml"
	plugin.StatisticsFilePath = ""
	plugin.StatisticsIP = ""
	plugin.StatisticsPort = 8053
	ok, err = checkArgs(nil)
	assert.Equal(ok, 3)
	assert.Error(err)

	// Set bad port
	plugin.StatisticsFormat = "xml"
	plugin.StatisticsFilePath = ""
	plugin.StatisticsIP = "127.0.0.1"
	plugin.StatisticsPort = 65536
	ok, err = checkArgs(nil)
	assert.Equal(ok, 3)
	assert.Error(err)

	plugin.StatisticsFormat = "json"
	plugin.StatisticsFilePath = ""
	plugin.StatisticsIP = "127.0.0.1"
	plugin.StatisticsPort = 8053
	ok, err = checkArgs(nil)
	assert.Equal(ok, 0)
	assert.NoError(err)
}

func TestExecuteCheck(t *testing.T) {
}
