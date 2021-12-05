package cli

import (
	"flag"
	"fmt"
	"os"

	"./github.com/bento1/cloneCoin/explorer"
	"./github.com/bento1/cloneCoin/rest"
)

func usage() {
	fmt.Printf("Welcom to Dong Coin\n")
	fmt.Printf("Please use the following flags\n")
	fmt.Printf("-port=4000 explorer : Start the HTML Explorer\n")
	fmt.Printf("-mode=rest : Start the REST API (recommenended)\n")
	os.Exit(0) // 0은 에러가 없음
}

//flag 패키지 사용
// go run main.go rest -port=4000 이런식으로 원하는 것을 실행 시키게해줌 -port 이걸 flag라고한다.
func Start() {
	// fmt.Println(os.Args)
	// if len(os.Args) < 2 {
	// 	usage()
	// }
	// restCommand := flag.NewFlagSet("rest", flag.ExitOnError)
	// portFlag := restCommand.Int("port", 4000, "Sets the port of the server (default 4000)")
	// switch os.Args[1] {
	// case "explorer":
	// 	fmt.Println("Start Explorer")
	// case "rest":
	// 	fmt.Println("Start REST API")
	// 	restCommand.Parse(os.Args[2:]) //"restCommand에서 Int에서 "port"인지 체크하고 찾는다 "

	// default:
	// 	usage()
	// }
	// if restCommand.Parsed() { //한번이라도 파싱됐다 (정상적으로)
	// 	//프로그램 수행
	// 	fmt.Println(*portFlag)
	// 	fmt.Println("Start Server")
	// }
	if len(os.Args) == 1 {
		usage()
	}
	port := flag.Int("port", 4000, "Set port of the server") //flag이름, 기본값, 에러시 출력문구
	mode := flag.String("mode", "rest", "Choose between html and rest")
	flag.Parse()
	fmt.Println(*port, *mode) // go run main.go -port 4000 -mode html
	switch *mode {
	case "html":
		fmt.Println("Start Explorer")
		explorer.Start(*port)
	case "rest":
		fmt.Println("Start REST API")
		rest.Start(*port)

	default:
		usage()
	}
}
