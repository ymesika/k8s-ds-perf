package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rcrowley/go-metrics"
)

var (
	server string
	port   int
	count  int
)

// Client ..
type Client struct {
	addr string
}

// Call ..
func (c *Client) Call() (string, time.Duration, error) {
	start := time.Now()
	resp, err := http.Get(c.addr)
	end := time.Now()
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	return string(body), end.Sub(start), nil
}

func main() {
	flag.Parse()

	if server == "" {
		fmt.Println("Missing 'server' flag")
		os.Exit(-1)
	}

	client := &Client{
		addr: fmt.Sprintf("http://%s:%d/", server, port),
	}

	fmt.Printf("Client will connect to [%s] %d times and collect stats\n", client.addr, count)

	s := metrics.NewUniformSample(count)

	for i := 0; i < count; i++ {
		resp, duration, err := client.Call()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		h := metrics.GetOrRegisterHistogram(resp, metrics.DefaultRegistry, s)
		h.Update(duration.Nanoseconds())
	}

	metrics.WriteOnce(metrics.DefaultRegistry, os.Stdout)
}

func init() {
	flag.IntVar(&port, "port", 8080, "remote server port")
	flag.StringVar(&server, "server", "", "remote server name")
	flag.IntVar(&count, "count", 1000, "number of calls to the server")
}
