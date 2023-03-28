package cloud

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func Run() {
}

// Provide a Terraform external datasource
// by communicating in JSON via stdin/stdout
func TerraformExternalDataSource() {
	var data []byte
	var err error
	if data, err = io.ReadAll(os.Stdin); err != nil {
		log.Fatal(err)
	}
	var params map[string]string
	if err = json.Unmarshal(data, &params); err != nil {
		log.Fatal(err)
	}
	params["HELLO"] = "TF INTERFACE WORKS!"
	if data, err = json.Marshal(params); err != nil {
		log.Fatal(err)
	}
	if _, err = os.Stdout.Write(data); err != nil {
		log.Fatal(err)
	}
}
