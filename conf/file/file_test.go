package file

import (
	"os"
	"path/filepath"
	"testing"

	"project-layout/pkg/conf"
)

const (
	testJSON = `
{
    "server":{
        "http":{
            "addr":"0.0.0.0",
			"port":80,
            "timeout":0.5,
			"enable_ssl":true
        },
        "grpc":{
            "addr":"0.0.0.0",
			"port":10080,
            "timeout":0.2
        }
    },
    "data":{
        "database":{
            "driver":"mysql",
            "source":"root:root@tcp(127.0.0.1:3306)/test?parseTime=true"
        }
    },
	"endpoints":[
		"www.aaa.com",
		"www.bbb.org"
	],
    "foo":[
        {
            "name":"nihao",
            "age":18
        },
        {
            "name":"nihao",
            "age":18
        }
    ]
}`

	testJSONUpdate = `
{
    "server":{
        "http":{
            "addr":"0.0.0.0",
			"port":80,
            "timeout":0.5,
			"enable_ssl":true
        },
        "grpc":{
            "addr":"0.0.0.0",
			"port":10090,
            "timeout":0.2
        }
    },
    "data":{
        "database":{
            "driver":"mysql",
            "source":"root:root@tcp(127.0.0.1:3306)/test?parseTime=true"
        }
    },
	"endpoints":[
		"www.aaa.com",
		"www.bbb.org"
	],
    "foo":[
        {
            "name":"nihao",
            "age":18
        },
        {
            "name":"nihao",
            "age":18
        }
    ],
	"bar":{
		"event":"update"
	}
}`
)

type testConfigStruct struct {
	Server struct {
		HTTP struct {
			Addr      string  `json:"addr" yaml:"addr"`
			Port      int     `json:"port" yaml:"port"`
			Timeout   float64 `json:"timeout" yaml:"timeout"`
			EnableSSL bool    `json:"enable_ssl" yaml:"enableSSL"`
		} `json:"http" yaml:"http"`
		GRPC struct {
			Addr    string  `json:"addr" yaml:"addr"`
			Port    int     `json:"port" yaml:"port"`
			Timeout float64 `json:"timeout" yaml:"timeout"`
		} `json:"grpc" yaml:"grpc"`
	} `json:"server"`
	Data struct {
		Database struct {
			Driver string `json:"driver" yaml:"driver"`
			Source string `json:"source" yaml:"source"`
		} `json:"database" yaml:"database" yaml:"database"`
	} `json:"data" yaml:"data"`
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
}

func TestFile(t *testing.T) {
	var (
		path = filepath.Join(os.TempDir(), "test_config")
		file = filepath.Join(path, "test.json")
		data = []byte(testJSON)
	)
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Error(err)
	}
	if err := os.WriteFile(file, data, 0o666); err != nil {
		t.Error(err)
	}

	if err := os.MkdirAll(path, 0700); err != nil {
		t.Error(err)
	}

	testSource(t, file)
	testSource(t, path)

}

func testSource(t *testing.T, path string) {
	t.Logf("path: %s", path)
	s := NewSource(path)
	kvs, err := s.Load()
	if err != nil {
		t.Error(err)
	}
	for _, f := range kvs {
		t.Logf("文件名 Key: %s, Format: %s, Data: %s", f.Key, f.Format, f.Value)
	}
}

func TestConfig(t *testing.T) {
	var (
		path  = filepath.Join(os.TempDir(), "test_config")
		file  = filepath.Join(path, "test.json")
		file2 = filepath.Join(path, "config.json")
		data  = []byte(testJSON)
		data2 = []byte(testJSONUpdate)
	)
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Error(err)
	}
	if err := os.WriteFile(file, data, 0o666); err != nil {
		t.Error(err)
	}

	if err := os.WriteFile(file2, data2, 0o666); err != nil {
		t.Error(err)
	}

	if err := os.MkdirAll(path, 0700); err != nil {
		t.Error(err)
	}

	c := conf.New(conf.WithSource(
		NewSource(path),
	))

	testConfig(t, c)
}

func testConfig(t *testing.T, c conf.Conf) {
	var (
		httpAddr       = "0.0.0.0"
		httpTimeout    = 0.5
		grpcPort       = 10080
		endpoint1      = "www.aaa.com"
		databaseDriver = "mysql"
	)

	c.Load()
	c.Watch()

	v := c.File("test").Get("server")
	t.Logf("app: %v", v)
	config := c.File("config").Get("server")
	t.Logf("app: %v", config)
	driver := c.File("config").GetString("data.database.driver")
	t.Logf("data.database.driver: %s", driver)

	if databaseDriver != driver {
		t.Fatal("databaseDriver is not equal to val")
	}

	var testConf testConfigStruct
	appConf := c.File("test")
	err := appConf.Unmarshal(&testConf)
	if err != nil {
		t.Errorf("error: %d", err)
	}
	t.Logf("AppConfig: %v", testConf)

	if httpAddr != testConf.Server.HTTP.Addr {
		t.Errorf("testConf.Server.HTTP.Addr want: %s, got: %s", httpAddr, testConf.Server.HTTP.Addr)
	}
	if httpTimeout != testConf.Server.HTTP.Timeout {
		t.Errorf("testConf.Server.HTTP.Timeout want: %.1f, got: %.1f", httpTimeout, testConf.Server.HTTP.Timeout)
	}
	if !testConf.Server.HTTP.EnableSSL {
		t.Error("testConf.Server.HTTP.EnableSSL is not equal to true")
	}
	if grpcPort != testConf.Server.GRPC.Port {
		t.Errorf("testConf.Server.GRPC.Port want: %d, got: %d", grpcPort, testConf.Server.GRPC.Port)
	}
	if endpoint1 != testConf.Endpoints[0] {
		t.Errorf("testConf.Endpoints[0] want: %s, got: %s", endpoint1, testConf.Endpoints[0])
	}
	if len(testConf.Endpoints) != 2 {
		t.Error("len(testConf.Endpoints) is not equal to 2")
	}
}
