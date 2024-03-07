package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"sscprovider/soltedev.pro/internal/libs/certificats"
	v2 "sscprovider/soltedev.pro/internal/libs/jks/v2"
)

var (
	kafka    = false
	kafkaSrv []string
	cleanup  = false
	dest     = ""
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "App is a CLI tool",
		Long:  `This is a CLI tool to demonstrate how to create a CLI application using Cobra library in Go.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if len(kafkaSrv) == 0 {
				kafka = false
			} else {
				kafka = true
			}

			fmt.Println("Starting the certificate generation process")
			run()
		},
	}

	rootCmd.PersistentFlags().StringSliceVarP(&kafkaSrv, "kafka", "k", []string{}, "list of Kafka servers that required certs")
	rootCmd.PersistentFlags().BoolVarP(&cleanup, "cleanup", "c", false, "cleanup data at the source")
	rootCmd.PersistentFlags().StringVarP(&dest, "dest", "d", "./", "specify destination folder")
	rootCmd.Execute()
}

func run() {

	if dest != "./" {
		if !strings.HasSuffix(dest, "/") {
			dest = dest + "/"
		}
	}

	_, err := os.Stat(dest)
	if err != nil {
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(dest, 0755)
			if errDir != nil {
				log.Fatal(errDir)
			}
		}
	}

	files, err := os.ReadDir(dest)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) != 0 {
		log.Fatalf("Directory is not empty: %v", dest)
	}

	//if cleanup{
	//	if filepath.Walk(dest == "./") {
	//
	//	}
	//	os.RemoveAll()
	//}

	//logger := logger.New()
	certificats.RootCA(dest, "sandbox")
	certificats.Client(dest, "sandbox", "client")
	certificats.Server(dest, "sandbox", "server")

	if kafka {
		fmt.Print("Enter password: ")

		for i := range kafkaSrv {
			kafkaSrv[i] = strings.TrimSpace(kafkaSrv[i])
		}

		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("\nPassword entered: ", string(password))

		v2.GenJKS(kafkaSrv, dest, "sandbox", string(password))
	}

	//cg := generate.New(logger)
	//jks := v2.New(logger)

	//if cleanup {
	//	utility.ClearDestFolder(dest)
	//}
	//
	//err := cg.RootCA("SandboxCA")
	//if err != nil {
	//	return
	//}
	//err = cg.Client("SandboxCA", "SandboxClient")
	//if err != nil {
	//	return
	//}
	//err = cg.Server("SandboxCA", "SandboxServer")
	//if err != nil {
	//	return
	//}
	//
	//if kafka {
	//	jks.CreateJKSCerts(kafkaSrv, "long-test-pass")
	//}
	//
	//if dest != "./" {
	//	utility.Move(dest)
	//}
	//
	//if cleanup {
	//	utility.Cleanup()
	//}
}
