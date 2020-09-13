package dbexport

import (
	"fmt"
	"os"

	"github.com/subosito/gotenv"
)

func GetConfig() {
	err := gotenv.Load(".env.local")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Println(os.Getenv("DB_DATABASE"))
}
