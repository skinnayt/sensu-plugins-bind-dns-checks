package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
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
	err = os.Chmod("tests/unreadable.stats", 0100)
	if err != nil {
		assert.FailNow("Unable to change file permissions")
	}
	plugin.StatisticsFormat = "file"
	plugin.StatisticsFilePath = "tests/unreadable.stats"
	ok, err = checkArgs(nil)
	assert.Equal(ok, 3)
	assert.Error(err)
	err = os.Chmod("tests/unreadable.stats", 0644)
	if err != nil {
		assert.FailNow("Unable to change file permissions")
	}

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
	err = os.Chmod("tests/unreadable.stats", 0100)
	if err != nil {
		assert.FailNow("Unable to change file permissions")
	}
	plugin.StatisticsFormat = "file"
	plugin.StatisticsFilePath = "tests/unreadable.stats"
	ok, err = executeCheck(nil)
	assert.Equal(ok, 2)
	assert.Error(err)
	err = os.Chmod("tests/unreadable.stats", 0644)
	if err != nil {
		assert.FailNow("Unable to change file permissions")
	}

	var namedXmlStats []byte
	var namedJsonStats []byte

	namedXmlStats, _ = os.ReadFile("tests/named.xml")
	namedJsonStats, _ = os.ReadFile("tests/named.json")

	tt := []struct {
		StatsFormat  string
		OutputFormat string
		Content      []byte
	}{
		{"xml", "", namedXmlStats},
		{"xml", "graphite", namedXmlStats},
		{"json", "", namedJsonStats},
		{"json", "graphite", namedJsonStats},
	}

	for _, tc := range tt {
		// Setup server for testing
		testConfig := &testServer{
			IP:          nil,
			Port:        0,
			StatsFormat: tc.StatsFormat,
			Content:     tc.Content,
			WaitChan:    make(chan bool),
		}

		testServe := startTestServer(testConfig)
		defer testServe.Close()

		go func() {
			testServe.Start()
			<-testConfig.WaitChan
			testServe.Close()
		}()

		// Wait for the server to start
		for testConfig.IP == nil {
			time.Sleep(1 * time.Second)
		}

		// Test the xml format
		plugin.StatisticsFormat = tc.StatsFormat
		plugin.StatisticsFilePath = ""
		plugin.OutputFormat = tc.OutputFormat
		plugin.StatisticsIP = testConfig.IP.String()
		plugin.StatisticsPort = testConfig.Port
		ok, err = executeCheck(nil)
		assert.Equal(ok, 0)
		assert.NoError(err)

		testConfig.WaitChan <- true

		// Shut down the server
		testServe.Close()
	}
}

type testServer struct {
	IP          net.IP
	Port        int
	StatsFormat string
	Content     []byte
	IsActive    bool
	WaitChan    chan bool
}

func startTestServer(runningServer *testServer) *httptest.Server {
	// Setup the test server
	// Load the data to return
	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure we have a GET request
		if r.Method != "GET" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		accept := strings.Split(r.Header.Get("Accept"), ",")

		// Make sure we have the correct path
		if runningServer.StatsFormat == "xml" && slices.Contains(accept, "application/xml") {
			if strings.HasPrefix(r.URL.Path, "/xml/v3") || r.URL.Path == "/" {
				// Return the xml
				w.Header().Set("Content-Type", "application/xml")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(runningServer.Content)
				if err != nil {
					http.Error(w, "Error writing to response", http.StatusInternalServerError)
					return
				}

				w.(http.Flusher).Flush()
				r.Context().Done()

				return
			}
		} else if runningServer.StatsFormat == "json" && slices.Contains(accept, "application/json") {
			if strings.HasPrefix(r.URL.Path, "/json/v1") {
				// Return the json
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(runningServer.Content)
				if err != nil {
					http.Error(w, "Error writing to response", http.StatusInternalServerError)
					return
				}

				w.(http.Flusher).Flush()
				r.Context().Done()

				return
			}
		}
	})
	dns_stats := httptest.NewUnstartedServer(httpHandler)

	runningServer.IsActive = true

	// Get the ip and port of the server
	var tempIP, tempPort string
	tempIP, tempPort, _ = net.SplitHostPort(dns_stats.Listener.Addr().String())
	runningServer.IP = net.ParseIP(tempIP)
	runningServer.Port, _ = strconv.Atoi(tempPort)

	return dns_stats
}
