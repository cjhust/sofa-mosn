package main
import "fmt"
import "runtime"
import "time"
import "os"
import "runtime/pprof"
import "runtime/trace"

var log *os.File

var s string = "test case hello world\n"

func test(id int) {
	i := 0
	for {
		i++
		//     fmt.Printf("%d %v test i = %d\n", id, time.Now(), i)
		//      println("test case hello world")
		log.WriteString(s)
	}
}


func main() {
	runtime.GOMAXPROCS(1)

	f, err := os.Create("test.pprof")
	if err == nil {
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	t, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer t.Close()

	err = trace.Start(t)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()


	log, _ = os.Create("/tmp/1.log")
	defer log.Close()

	fmt.Println("test begin")
	for i := 0; i < 2; i++ {
		go test(i)
	}

	end := time.NewTimer(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)


	for {
		select {
		case <- end.C:
			fmt.Println("test end")
			return
		case <- ticker.C:
			fmt.Println("ticker")
		}
	}
}

