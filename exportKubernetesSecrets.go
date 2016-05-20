package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("usage exportKubernetesSecrets <k8s secret>")
	}
	secret := os.Args[1]
	os.Mkdir(secret, 0700)
	cmd := exec.Command("kubectl", "get", "secret", secret, "-o", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("An error ocurred while running kubectl: %s\n", err)
		log.Fatal(string(out))
	} else {
		var f map[string]interface{}
		err := json.Unmarshal(out, &f)
		if err != nil {
			fmt.Println("Cant read json from kubectl")
		} else {
			files := f["data"].(map[string]interface{})
			for k, v := range files {
				b64str := v.(string)
				data, err := base64.StdEncoding.DecodeString(b64str)
				if err != nil {
					log.Fatal("error:", err)
				}
				ioutil.WriteFile(secret+string(os.PathSeparator)+k, data, 0755)
			}
		}
	}
}
