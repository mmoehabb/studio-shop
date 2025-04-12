package ancillaries

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Aws_Access_Key_Id     string
	Aws_Secret_Access_Key string
}

var envLoaded = false
var env Env

func LoadEnv() {
	if envLoaded == true {
		return
	}
	if err := godotenv.Load(); err != nil {
		panic(".env file missing!")
	}
	envLoaded = true
	env = Env{
		Aws_Access_Key_Id:     os.Getenv("AWS_ACCESS_KEY_ID"),
		Aws_Secret_Access_Key: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
}

func GetEnv() *Env {
	if envLoaded == false {
		LoadEnv()
	}
	return &env
}
