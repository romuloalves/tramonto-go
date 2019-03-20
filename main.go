package main

import (
	"encoding/json"
	"fmt"
	"time"

	tramonto "gitlab.com/tramonto-one/go-tramonto/tramonto"
)

func main() {
	repoPath := "/Users/romuloalves/.gvm/pkgsets/go1.11.4/global/src/gitlab.com/tramonto-one/go-tramonto/.temp-ipfs"

	done := make(chan bool, 1)

	go func() {
		time.Sleep(20 * time.Minute)
		done <- true
	}()

	one, err := tramonto.NewTramontoOne(repoPath)
	if err != nil {
		panic(err)
	}

	if err := one.Start(); err != nil {
		panic(err)
	}

	// test, err := one.NewTest("TR0003", "My cool description!")
	// if err != nil {
	// 	panic(err)
	// }

	// jsonData, err := json.Marshal(test)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("\n\n%s\n\n", string(jsonData))

	getResult, err := one.GetTest("QmTJMBr9dTqkw321ne3fUPYfQRSgs8Ut2UqUjDQsv9DUmr", "")

	// getResult, err := one.GetTest(test.IpnsHash, test.Secret)
	if err != nil {
		panic(err)
	}

	jsonDataTwo, err := json.Marshal(getResult)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n>>>>> \n\n%s\n\n", string(jsonDataTwo))

	<-done
}
