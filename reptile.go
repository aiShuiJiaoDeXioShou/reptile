package main

import (
	"io/ioutil"
	"reptile/service"
)

func main() {
	// 读取指定文件并且打印
	b, err := ioutil.ReadFile("static/json/logo.txt")
	if err != nil {
		println(" ▄▄▄▄▄▄▄▄▄▄▄  ▄    ▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄       ▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄        ▄  ▄         ▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄    ▄  ▄▄▄▄▄▄▄▄▄▄▄            ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄ \n▐░░░░░░░░░░░▌▐░▌  ▐░▌▐░░░░░░░░░░░▌▐░░▌     ▐░░▌▐░░░░░░░░░░░▌▐░░▌      ▐░▌▐░▌       ▐░▌▐░░░░░░░░░░░▌▐░▌  ▐░▌▐░░░░░░░░░░░▌          ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░▌\n▐░█▀▀▀▀▀▀▀█░▌▐░▌ ▐░▌ ▐░█▀▀▀▀▀▀▀█░▌▐░▌░▌   ▐░▐░▌▐░█▀▀▀▀▀▀▀▀▀ ▐░▌░▌     ▐░▌▐░▌       ▐░▌ ▀▀▀▀█░█▀▀▀▀ ▐░▌ ▐░▌  ▀▀▀▀█░█▀▀▀▀           ▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌▐░▌\n▐░▌       ▐░▌▐░▌▐░▌  ▐░▌       ▐░▌▐░▌▐░▌ ▐░▌▐░▌▐░▌          ▐░▌▐░▌    ▐░▌▐░▌       ▐░▌     ▐░▌     ▐░▌▐░▌       ▐░▌               ▐░▌          ▐░▌       ▐░▌▐░▌\n▐░█▄▄▄▄▄▄▄█░▌▐░▌░▌   ▐░▌       ▐░▌▐░▌ ▐░▐░▌ ▐░▌▐░█▄▄▄▄▄▄▄▄▄ ▐░▌ ▐░▌   ▐░▌▐░▌   ▄   ▐░▌     ▐░▌     ▐░▌░▌        ▐░▌               ▐░▌ ▄▄▄▄▄▄▄▄ ▐░▌       ▐░▌▐░▌\n▐░░░░░░░░░░░▌▐░░▌    ▐░▌       ▐░▌▐░▌  ▐░▌  ▐░▌▐░░░░░░░░░░░▌▐░▌  ▐░▌  ▐░▌▐░▌  ▐░▌  ▐░▌     ▐░▌     ▐░░▌         ▐░▌               ▐░▌▐░░░░░░░░▌▐░▌       ▐░▌▐░▌\n▐░█▀▀▀▀▀▀▀▀▀ ▐░▌░▌   ▐░▌       ▐░▌▐░▌   ▀   ▐░▌▐░█▀▀▀▀▀▀▀▀▀ ▐░▌   ▐░▌ ▐░▌▐░▌ ▐░▌░▌ ▐░▌     ▐░▌     ▐░▌░▌        ▐░▌               ▐░▌ ▀▀▀▀▀▀█░▌▐░▌       ▐░▌▐░▌\n▐░▌          ▐░▌▐░▌  ▐░▌       ▐░▌▐░▌       ▐░▌▐░▌          ▐░▌    ▐░▌▐░▌▐░▌▐░▌ ▐░▌▐░▌     ▐░▌     ▐░▌▐░▌       ▐░▌               ▐░▌       ▐░▌▐░▌       ▐░▌ ▀ \n▐░▌          ▐░▌ ▐░▌ ▐░█▄▄▄▄▄▄▄█░▌▐░▌       ▐░▌▐░█▄▄▄▄▄▄▄▄▄ ▐░▌     ▐░▐░▌▐░▌░▌   ▐░▐░▌ ▄▄▄▄█░█▄▄▄▄ ▐░▌ ▐░▌  ▄▄▄▄█░█▄▄▄▄           ▐░█▄▄▄▄▄▄▄█░▌▐░█▄▄▄▄▄▄▄█░▌ ▄ \n▐░▌          ▐░▌  ▐░▌▐░░░░░░░░░░░▌▐░▌       ▐░▌▐░░░░░░░░░░░▌▐░▌      ▐░░▌▐░░▌     ▐░░▌▐░░░░░░░░░░░▌▐░▌  ▐░▌▐░░░░░░░░░░░▌          ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░▌\n ▀            ▀    ▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀         ▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀        ▀▀  ▀▀       ▀▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀    ▀  ▀▀▀▀▀▀▀▀▀▀▀            ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀  ▀")
	}else {
		println(string(b))
	}
	service.StartUp()
}
