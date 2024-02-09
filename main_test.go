package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

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
	assert := assert.New(t)

	// Test the file format
	plugin.StatisticsFormat = "file"
	plugin.StatisticsFilePath = "tests/named.stats"
	ok, err := executeCheck(nil)
	assert.Equal(ok, 0)
	assert.NoError(err)

	// Test not able to read file
	os.Chmod("tests/unreadable.stats", 0100)
	plugin.StatisticsFormat = "file"
	plugin.StatisticsFilePath = "tests/unreadable.stats"
	ok, err = executeCheck(nil)
	assert.Equal(ok, 2)
	assert.Error(err)
	os.Chmod("tests/unreadable.stats", 0644)

	var bindXmlStats []byte
	var bindJsonStats []byte

	statsXmlFile, _ := os.Open("tests/named.xml")
	statsXmlFile.Read(bindXmlStats)
	statsXmlFile.Close()

	statsJsonFile, _ := os.Open("tests/named.json")
	statsJsonFile.Read(bindJsonStats)
	statsJsonFile.Close()

	// Setup server for testing
	xmlServer := &testServer{
		IP:          nil,
		Port:        0,
		StatsFormat: "xml",
		Content:     bindXmlStats,
	}

	xmlServe := startTestServer(xmlServer)

	defer xmlServe.Close()

	// Wait for the server to start
	for xmlServer.IP == nil {
		time.Sleep(1 * time.Second)
	}

	// Test the xml format
	plugin.StatisticsFormat = "xml"
	plugin.StatisticsFilePath = ""
	plugin.StatisticsIP = xmlServer.IP.String()
	plugin.StatisticsPort = xmlServer.Port
	ok, err = executeCheck(nil)
	assert.Equal(ok, 0)
	assert.NoError(err)

	// Shut down the server
	xmlServe.Close()

	// Setup server for testing json
	jsonServer := &testServer{
		IP:          nil,
		Port:        0,
		StatsFormat: "json",
		Content:     bindJsonStats,
	}

	jsonServe := startTestServer(jsonServer)

	defer jsonServe.Close()

	// Wait for the server to start
	for jsonServer.IP == nil {
		time.Sleep(1 * time.Second)
	}

	// Test the json format
	plugin.StatisticsFormat = "json"
	plugin.StatisticsFilePath = ""
	plugin.StatisticsIP = jsonServer.IP.String()
	plugin.StatisticsPort = jsonServer.Port
	ok, err = executeCheck(nil)
	assert.Equal(ok, 0)
	assert.NoError(err)

	// Shut down the server
	jsonServe.Close()

}

type testServer struct {
	IP          net.IP
	Port        int
	StatsFormat string
	Content     []byte
}

func startTestServer(runningServer *testServer) *httptest.Server {
	// Setup the test server
	// Load the data to return
	dns_stats := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure we have a GET request
		if r.Method != "GET" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Make sure we have the correct path
		if runningServer.StatsFormat == "xml" {
			if strings.HasPrefix(r.URL.Path, "/xml/v3") || r.URL.Path == "/" {
				// Return the xml
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				w.Write(runningServer.Content)
				return
			}
		} else if runningServer.StatsFormat == "json" {
			if strings.HasPrefix(r.URL.Path, "/json/v1") {
				// Return the json
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(runningServer.Content)
				return
			}
		}
	}))

	// Get the ip and port of the server
	var tempIP, tempPort string
	tempIP, tempPort, _ = net.SplitHostPort(dns_stats.Listener.Addr().String())
	runningServer.IP = net.ParseIP(tempIP)
	runningServer.Port, _ = strconv.Atoi(tempPort)

	return dns_stats
}
