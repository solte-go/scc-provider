package v2

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func GenJKS(roles []string, dest, ca, password string) {
	for _, role := range roles {

		_, err := os.Stat(dest + role)
		if err != nil {
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(dest+role, 0755)
				if errDir != nil {
					log.Fatal(err)
				}
			}
		}

		commands := [][]string{
			{"keytool", "-keystore", dest + role + "/kafka." + role + ".truststore.jks", "-alias", "CARoot", "-import", "-file", dest + ca + ".crt", "-storepass", password, "-keypass", password, "-noprompt"},
			{"keytool", "-keystore", dest + role + "/kafka." + role + ".keystore.jks", "-alias", "localhost", "-validity", "30", "-genkey", "-keyalg", "RSA", "-dname", "CN=localhost, OU=Dev, O=Box, L=Vantaa, ST=Uusimaa, C=FI", "-ext", "SAN=DNS:localhost", "-storepass", password, "-keypass", password, "-noprompt"},
			{"keytool", "-keystore", dest + role + "/kafka." + role + ".keystore.jks", "-alias", "localhost", "-certreq", "-file", dest + role + "/kafka." + role + ".unsigned.crt", "-storepass", password, "-keypass", password, "-noprompt"},

			{"openssl", "x509", "-req", "-CA", dest + ca + ".crt", "-CAkey", dest + ca + ".key", "-in", dest + role + "/kafka." + role + ".unsigned.crt", "-out", dest + role + "/kafka." + role + ".signed.crt", "-days", "30", "-CAcreateserial", "-passin", "pass:" + password},
			{"keytool", "-keystore", dest + role + "/kafka." + role + ".keystore.jks", "-alias", "CARoot", "-import", "-file", dest + ca + ".crt", "-storepass", password, "-keypass", password, "-noprompt"},
			{"keytool", "-keystore", dest + role + "/kafka." + role + ".keystore.jks", "-alias", "localhost", "-import", "-file", dest + role + "/kafka." + role + ".signed.crt", "-storepass", password, "-keypass", password, "-noprompt"},
		}

		for _, command := range commands {
			cmd := exec.Command(command[0], command[1:]...)

			// Stdout and stderr buffers
			cmdOutput := &bytes.Buffer{}
			cmdError := &bytes.Buffer{}

			// Attach stdout and stderr buffers to command
			cmd.Stdout = cmdOutput
			cmd.Stderr = cmdError

			err := cmd.Run()
			if err != nil {
				fmt.Printf("Output: %s\n", cmdOutput.String())
				fmt.Printf("Error: %s\n", cmdError.String())
				log.Fatal(err)
			}
		}

		filenames := []string{role + "_sslkey.creds", role + "_keystore.creds", role + "_truststore.creds"}

		for _, filename := range filenames {
			file, err := os.Create(dest + role + "/" + filename)
			if err != nil {
				log.Fatal(err)
			}

			defer file.Close()

			_, err = file.WriteString(password + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
