package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"

	"github.com/spf13/cobra"
	//"time"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "8010"
	CONN_TYPE = "tcp"
	CONN_BASE = "http://localhost:4044/api/gateway"
)

type meshData struct {
	addr string `json:"addr"`
	data string `json:"data"`
}

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

type accessRequest struct {
	MAC  string `json:"MAC"`
	Card string `json:"Card"`
	Mid  uint64 `json:"mid"`
}

type hbRequest struct {
	MAC  string `json:"MAC"`
	Tick uint64 `json:"tick"`
}

type meshResponse struct {
	Addr string                 `json:"addr"`
	Data map[string]interface{} `json:"data"`
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.

	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	dlen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// `&myStoredVariable` is the address of the variable we want to store our
	// parsed data in

	var data map[string]interface{}

	if err := json.Unmarshal(buf[:dlen], &data); err != nil {
		//panic(err)
		fmt.Println("error:", err)
		return
	}

	if _, ok := data["addr"]; !ok {
		return
	}

	if _, ok := data["data"]; !ok {
		return
	}

	reqT := data["data"].(map[string]interface{})["type"]
	if reqT != nil {
		fmt.Printf("T: %+v\n", reflect.TypeOf(reqT))
		rType := int64(reqT.(float64))
		fmt.Printf("T: %v %v\n", rType, rType == 1)
		if rType == 0 {
			acc := &accessRequest{
				data["addr"].(string),
				data["data"].(map[string]interface{})["card"].(string),         /**/
				uint64(data["data"].(map[string]interface{})["mid"].(float64)), /**/
			}

			go handleAccPost(acc, conn)
		} else if rType == 1 {
			acc := &accessRequest{
				data["addr"].(string),
				data["data"].(map[string]interface{})["card"].(string), /**/
				uint64(data["data"].(map[string]interface{})["mid"].(float64)),
			}

			go handlePropPost(acc, conn)
		} else {

			//fmt.Printf("%+v\n", data)
			//fmt.Printf("%s\n", buf[:dlen])

		}
	}

	tc := data["data"].(map[string]interface{})["tick"]
	if tc != nil {
		rg := &hbRequest{
			data["addr"].(string),
			uint64(data["data"].(map[string]interface{})["tick"].(float64)),
		}
		go handleHB(rg, conn)
	}

	tm := data["data"].(map[string]interface{})["time"]
	if tm != nil {
		conn.Close()
		fmt.Printf("\nMessage round trip: %v\n", tm)

		nm := fmt.Sprintf("trip-%v.log", data["addr"].(string))
		f, err := os.OpenFile(nm,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		txt := fmt.Sprintf("%v\n", tm)
		if _, err := f.WriteString(txt); err != nil {
			fmt.Println(err)
		}
	}
	//err = json.Unmarshal(buf[:dlen], &data)
	//enc := base64.StdEncoding.EncodeToString(buf[:dlen])

	//fmt.Println(string(enc))
	//fmt.Printf("%+v\n", data.addr)
	//time.Sleep(2 * time.Second)
	// Send a response back to person contacting us.
	//conn.Write(buf[:dlen])
	// Close the connection when you're done with it.
	//conn.Close()
}

func handleAccPost(data *accessRequest, conn net.Conn) {
	//req, err := http.NewRequest("POST", CONN_BASE, bytes.NewBuffer(data))
	//req.Header.Set("X-Custom-Header", "myvalue")
	//req.Header.Set("Content-Type", "application/json")
	defer conn.Close()

	b, err := json.MarshalIndent(&data, "", "\t")
	fmt.Printf("sending: %+v\n", b)
	if err != nil {
		fmt.Println("error:", err)
	}

	var host = gatewayIP

	if len(gatewayPort) != 0 {
		host += ":" + gatewayPort
	}

	host += "/api/gateway/access"

	resp, err := http.Post(CONN_BASE+"/access", "base64", bytes.NewBuffer(b))

	if err != nil {
		fmt.Println(err)
		//panic(err)
	} else {

		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		//body, _ := ioutil.ReadAll(resp.Body)

		res := &meshResponse{
			Addr: data.MAC,
			//Data: string(body),
		}

		b, err := json.MarshalIndent(&res, "", "\t")
		fmt.Printf("sending: %+v\n", b)
		if err != nil {
			fmt.Println("error:", err)
		}

		n, err := conn.Write(b)

		fmt.Println("write: %+v %+v\n", n, err)
	}

}

func handlePropPost(data *accessRequest, conn net.Conn) {
	//req, err := http.NewRequest("POST", CONN_BASE, bytes.NewBuffer(data))
	//req.Header.Set("X-Custom-Header", "myvalue")
	//req.Header.Set("Content-Type", "application/json")
	defer conn.Close()

	b, err := json.MarshalIndent(&data, "", "\t")
	//fmt.Printf("sending: %+v\n", b)
	if err != nil {
		fmt.Println("error:", err)
	}

	var host = gatewayIP

	if len(gatewayPort) != 0 {
		host += ":" + gatewayPort
	}

	host += "/api/gateway/props"

	resp, err := http.Post(CONN_BASE+"/props", "base64", bytes.NewBuffer(b))

	if err != nil {
		fmt.Println(err)
		//panic(err)
	} else {

		defer resp.Body.Close()
		//fmt.Println("response Status:", resp.Status)
		//fmt.Println("response Headers:", resp.Header)

		body := make([]byte, 2048)
		// Read the incoming connection into the buffer.
		blen, err := resp.Body.Read(body)

		//body, _ := ioutil.ReadAll(resp.Body)

		var resData map[string]interface{}

		if err := json.Unmarshal(body[:blen], &resData); err != nil {
			//panic(err)
			fmt.Println("error:", err)
			return
		}

		res := &meshResponse{
			Addr: data.MAC,
			Data: resData,
		}

		b, err := json.MarshalIndent(&res, "", "\t")
		//fmt.Printf("sending: %s\n", b)
		if err != nil {
			fmt.Println("error:", err)
		}

		/*n, err :=*/
		conn.Write(b)

		//fmt.Println("write: %+v %+v\n", n, err)
	}

}

func handleHB(data *hbRequest, conn net.Conn) {
	//req, err := http.NewRequest("POST", CONN_BASE, bytes.NewBuffer(data))
	//req.Header.Set("X-Custom-Header", "myvalue")
	//req.Header.Set("Content-Type", "application/json")
	defer conn.Close()

	b, err := json.MarshalIndent(&data, "", "\t")
	//fmt.Printf("sending: %+v\n", b)
	if err != nil {
		fmt.Println("error:", err)
	}

	var host = gatewayIP

	if len(gatewayPort) != 0 {
		host += ":" + gatewayPort
	}

	host += "/api/gateway/props"

	resp, err := http.Post(CONN_BASE+"/hb", "base64", bytes.NewBuffer(b))

	if err != nil {
		fmt.Println(err)
		//panic(err)
	} else {

		defer resp.Body.Close()
		//fmt.Println("response Status:", resp.Status)
		//fmt.Println("response Headers:", resp.Header)

		body := make([]byte, 2048)
		// Read the incoming connection into the buffer.
		blen, err := resp.Body.Read(body)

		//body, _ := ioutil.ReadAll(resp.Body)

		var resData map[string]interface{}

		if err := json.Unmarshal(body[:blen], &resData); err != nil {
			//panic(err)
			fmt.Println("error:", err)
			return
		}

		res := &meshResponse{
			Addr: data.MAC,
			Data: resData,
		}

		b, err := json.MarshalIndent(&res, "", "\t")
		//fmt.Printf("sending: %s\n", b)
		if err != nil {
			fmt.Println("error:", err)
		}

		/*n, err :=*/
		conn.Write(b)

		//fmt.Println("write: %+v %+v\n", n, err)
	}

}

//gatewayIP != ""
