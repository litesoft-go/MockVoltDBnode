package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/litesoft-go/mockvoltdb/pkg/canned"
	"github.com/litesoft-go/mockvoltdb/pkg/utils"
	"github.com/litesoft-go/mockvoltdb/version"
)

const (
	PORT_0_http     = 8080
	PORT_1_internal = 3021
	PORT_7_status   = 11780
	// PORT_2_rep   = 5555
	// PORT_3_zk    = 7181
	// PORT_4_jmx   = 9090
	// PORT_5_admin = 21211
	// PORT_6_client= 21212
)

var responses = make(map[string]*utils.Responder)

func splitHost(pHost string) (rHost, rPort string) {
	rHost = pHost
	if at := strings.IndexByte(pHost, ':'); at != -1 {
		rHost = pHost[:at]
		rPort = pHost[at:]
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	host, port := splitHost(r.Host)
	url := r.URL
	query := url.RawQuery
	urlPathQuery := url.Path
	if query != "" {
		urlPathQuery += "?" + query
	}
	resp := responses[port].GetResponse(urlPathQuery)
	cType := "application/json"
	body := resp.GetBody()
	status := resp.GetStatus()
	msg := fmt.Sprintf("Request: %s%s%s%s", ourPrivateClassA, host, port, urlPathQuery)
	if status == 404 {
		msg = "**** 404 - " + msg
		body = "port[" + port + "] " + body
		cType = "text/plain"
	}
	fmt.Println(msg)
	w.Header().Set("Content-Type", cType)
	w.WriteHeader(status)
	count, err := fmt.Fprint(w, body)
	if err != nil {
		fmt.Printf("******************** type '%s' %d, bytes: %d, error: %v\n", cType, status, count, err)
	}
}

var ourPrivateClassA = ""

func main() {
	fmt.Printf("MockVoltDB Version: %s\n", version.Version)
	fmt.Printf("Args: %v\n", os.Args[1:])

	ipv4s, err := utils.IPv4s()

	fmt.Println("(IPv4s):")
	for _, ipv4 := range ipv4s {
		fmt.Printf("   %s\n", ipv4.String())
		if ipv4.IsPrivateClassA() {
			ourPrivateClassA = ipv4.String() + "|"
		}
	}
	if err == nil {
		var dir string
		dir, err = os.Getwd()
		if err == nil {
			fmt.Println("Working Dir: ", dir)
		}
	}
	if err == nil {
		if (len(os.Args) > 1) && (os.Args[1] == "init") {
			err = createFile("voltdbroot", ".initialized", "")
			if err == nil {
				err = createFile("voltdbroot/config", "path.properties", "forGoodTime=555-123-4567\n")
			}
			if err == nil {
				fmt.Println("Initialize! ********************************************")
				os.Exit(0)
			}
			fmt.Println("Initialize! SHIT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n", err)
		} else {
			responses[":11780"] = canned.PortResonses_status() // PORT_7_status___
			responses[":8080"] = canned.PortResonses_http()    // PORT_0_http_____

			http.HandleFunc("/", handler)
			err = doAllWork(
				PORT_0_http,
				PORT_1_internal,
				PORT_7_status,
				// PORT_2_rep,
				// PORT_3_zk,
				// PORT_4_jmx,
				// PORT_5_admin,
				// PORT_6_client,
			)
		}
	}
	log.Fatal(err)
}

func createFile(dirs, filename, body string) (err error) {
	if dirs != "" {
		dirs, err = makeDirs(dirs)
	}
	if (err == nil) && (filename != "") {
		err = makeFile(dirs+filename, body)
	}
	return
}

func makeDirs(dirs string) (pathDir string, err error) {
	_, err = os.Stat(dirs)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirs, 0755) // read & execute: everyone; write: owner
	}
	if err == nil {
		pathDir = dirs + "/"
	}
	return
}

func makeFile(filepath, body string) (err error) {
	_, err = os.Stat(filepath)
	if os.IsNotExist(err) {
		var file *os.File
		file, err = os.Create(filepath)
		if err == nil {
			_, err = file.WriteString(body)
			err2 := file.Close()
			if err == nil {
				err = err2
			}
		}
	}
	return
}

func doAllWork(ports ...int) error {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure it's called to release resources even if no errors

	for _, port := range ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()

			err := work(port)
			if err != nil {
				fmt.Printf("Worker #%d, error: %v\n", port, err)
				cancel()
				return
			}
		}(port)
	}
	wg.Wait()

	return ctx.Err()
}

func work(port int) error {
	fmt.Printf("Listening on %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(addr, nil)
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}
