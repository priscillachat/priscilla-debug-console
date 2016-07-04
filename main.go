package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/priscillachat/prislog"
	"gopkg.in/readline.v1"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const (
	engagement = `{
		"type": "command",
		"source": "debug-console",
		"to": "server",
		"command":
		{
			"action": "engage",
			"type": "%s",
			"time": %d,
			"data":"%s"
		}
	}`
)

var logger *prislog.PrisLog

func main() {

	host := flag.String("server", "127.0.0.1", "Priscilla server host")
	port := flag.Int("port", 4517, "Priscilla server port")
	mode := flag.String("mode", "responder", "adapter or responder")
	secret := flag.String("secret", "abcdefghi",
		"secret for access priscilla server")

	flag.Parse()

	var err error
	logger, err = prislog.NewLogger(os.Stdout, "error")

	if err != nil {
		fmt.Println("Error initializing logger: ", err)
		os.Exit(-1)
	}

	prefix := fmt.Sprint(*host, ":", *port)
	rl, err := readline.New(prefix + " > ")

	if err != nil {
		logger.Error.Fatal("Error initializing readline:", err)
	}
	defer rl.Close()

	conn, err := net.Dial("tcp", prefix)

	if err != nil {
		logger.Error.Fatal("Error connecting to Priscilla server:", err)
	}

	go listen(conn)

	timestamp := time.Now().Unix()
	authMsg := fmt.Sprintf("%d%s%s", timestamp, "debug-console", *secret)

	mac := hmac.New(sha256.New, []byte(*secret))
	mac.Write([]byte(authMsg))

	fmt.Fprintf(conn, engagement, *mode, timestamp,
		base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	time.Sleep(time.Second)

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := rl.Readline()
		if err != nil {
			if err.Error() != "EOF" {
				logger.Error.Fatal("Unexpected error:", err)
			}
			break
		}
		trimmed := strings.Trim(line, " ")

		switch trimmed {
		case "exit":
			break
		case "put":
			reader.Reset(os.Stdin)
			str, _ := reader.ReadString('\x00')
			fmt.Printf("\nReceived:\n%s\n", str)
			fmt.Fprintf(conn, str)
			if err != nil {
				logger.Error.Println("Error sending:", err)
			}
		}
	}
}

func listen(reader io.Reader) {
	buf := make([]byte, 4096)

	for {
		count, err := reader.Read(buf)

		if err == nil {
			fmt.Printf("Received from server:\n%s\n", string(buf[:count]))
		} else {
			logger.Error.Println(err)
			time.Sleep(time.Second)
		}
	}
}
