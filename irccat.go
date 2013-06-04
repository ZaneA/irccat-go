package main

// I'm sure this list can be trimmed a bit...
import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

// Flags
var verbose = flag.Bool("verbose", false, "Be verbose (and stay in the foreground)")
var nick = flag.String("nick", "irccat", "Nickname")
var dest = flag.String("dest", "", "Where to send the text (nickname or channel)")
var server = flag.String("server", "", "Server to send to (e.g. chat.freenode.net:6667)")

// debug prints a formatted string when the verbose flag is true.
func debug(format string, v ...interface{}) {
	if *verbose {
		fmt.Printf(format, v...)
	}
}

// send sends a formatted string to the IRC server.
func send(conn net.Conn, format string, v ...interface{}) {
	debug("> "+format, v...)

	fmt.Fprintf(conn, format, v...)
	time.Sleep(700 * time.Millisecond) // Wait a bit so we don't flood
}

// sendMessage sends the text in s to dest, splitting by newline if necessary.
func sendMessage(conn net.Conn, dest string, s string) {
	parts := strings.Split(s, "\n")
	for i := 0; i < len(parts); i++ {
		if parts[i] != "" {
			send(conn, "NOTICE %s :%s\r\n", dest, parts[i])
		}
	}
}

func main() {
	flag.Parse()

	// I imagine that there's a far more idiomatic way to do this.
	if *dest == "" {
		fmt.Println("You must specify a destination (see --help)")
		os.Exit(1)
	}

	if *server == "" {
		fmt.Println("You must specify a server (see --help)")
		os.Exit(1)
	}

	// Read STDIN.
	messageBytes, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Println("Error reading STDIN")
		os.Exit(1)
	}

	message := string(messageBytes)

	debug("Connecting to %s\n", *server)

	// "Dial" the connection to the server.
	conn, err := net.Dial("tcp", *server)
	if err != nil {
		fmt.Println("Couldn't connect")
		os.Exit(1)
	}

	// Create a buffered reader for the connection.
	r := bufio.NewReader(conn)

	// Introduce ourselves to the server.
	send(conn, "NICK %s\r\n", *nick)
	send(conn, "USER %s * * :%s\r\n", *nick, *nick)

	for {
		// Read in a line at a time.
		line, err := r.ReadString('\n')

		if err != nil {
			break
		}

		debug("< %s", line)

		// If the response contains 422 or 376 we're ready to rock.
		if strings.Contains(line, "422") || strings.Contains(line, "376") {
			// Is the destination a channel?
			if strings.HasPrefix(*dest, "#") || strings.HasPrefix(*dest, "&") ||
				strings.HasPrefix(*dest, "!") || strings.HasPrefix(*dest, "+") {
				send(conn, "JOIN %s\r\n", *dest)
			}

			sendMessage(conn, *dest, message)

			send(conn, "QUIT :irccat v0.1\r\n")
			break
		}
	}

	conn.Close()
}
