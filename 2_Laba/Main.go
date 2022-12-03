package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

var abort bool                   
var exitchnl = make(chan int, 2) 


func DFileR(path string, filesChan chan string) { 
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".csv" {
			filesChan <- filepath.Join(path, file.Name())
		} else if file.IsDir() {
			DFileR(filepath.Join(path, file.Name()), filesChan)
		}
	}
}


func FFileD(path string, filesChan chan string, inputfileName string, abort bool) {
	if abort {
		return
	}

	files, err := os.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == inputfileName && !abort {
			filesChan <- filepath.Join(path, file.Name())
			abort = true
		} else if file.IsDir() {
			if abort {
				break
			}
			FFileD(filepath.Join(path, file.Name()), filesChan, inputfileName, abort)
		}
	}
}


type tree struct { 
	data  []string 
	left  *tree    
	right *tree    
}


func createTreeV(buffer string) *tree { 
	var vertex = new(tree)
	vertex.data = strings.Split(buffer, ";") 
	return vertex
}


func addBranch(vertex *tree, buffer string, sortByLine *int) {
	var compare = strings.Split(buffer, ";")
	if compare[*sortByLine] > vertex.data[*sortByLine] {
		if vertex.right == nil {
			vertex.right = new(tree)
			vertex.right.data = compare
		} else {
			addBranch(vertex.right, buffer, sortByLine)
		}
	} else if compare[*sortByLine] < vertex.data[*sortByLine] {
		if vertex.left == nil {
			vertex.left = new(tree)
			vertex.left.data = compare
		} else {
			addBranch(vertex.left, buffer, sortByLine)
		}
	} else if compare[*sortByLine] == vertex.data[*sortByLine] {
		if vertex.left == nil {
			vertex.left = new(tree)
			vertex.left.data = compare
		} else {
			insertElem(vertex, compare)
		}
	}
}


func insertElem(vertex *tree, data []string) {
	temp := vertex.left
	vertex.left = new(tree)
	vertex.left = &tree{left: temp, data: data}
}


func outTree(vertex *tree, file *bufio.Writer) { // recursively derive the value in ascending order
	if vertex.left != nil {
		outTree(vertex.left, file)
	}
	writeFile(file, vertex)
	if vertex.right != nil {
		outTree(vertex.right, file)
	}
}

//Output values ​​to the console 
func outCLI(vertex *tree) { // recursively derive the value in ascending order
	if vertex.left != nil {
		outCLI(vertex.left)
	}
	println(strings.Join(vertex.data, ";"))
	if vertex.right != nil {
		outCLI(vertex.right)
	}
}
func outCLIRev(vertex *tree) { // recursively output the value in descending order
	if vertex.right != nil {
		outCLIRev(vertex.right)
	}
	println(strings.Join(vertex.data, ";"))
	if vertex.left != nil {
		outCLIRev(vertex.left)
	}
}

//Outputting values ​​to a file
func outTreeR(vertex *tree, file *bufio.Writer) { // recursively output the value in descending order
	if vertex.right != nil {
		outTreeR(vertex.right, file)
	}
	writeFile(file, vertex)
	if vertex.left != nil {
		outTreeR(vertex.left, file)
	}
}
func writeFile(file *bufio.Writer, temp *tree) { // output the data array of each branch
	file.WriteString(strings.Join(temp.data, ";"))
	file.WriteByte('\n')
}

//Program exit signal handler
func handler(signal os.Signal) {
    if signal == syscall.SIGTERM {
        fmt.Println("Got kill signal. ")
        fmt.Println("Program will terminate now.")
        os.Exit(0)
    } else if signal == syscall.SIGINT {
        fmt.Println("Got CTRL+C signal.")
        fmt.Println("Closing.")
        os.Exit(0)
    } else {
        fmt.Println("Ignoring signal: ", signal)
    }
}




func main() {
	const go_size int = 3 

	var (
		path           = flag.String("d", ".", "Use a file with the name file-name as an input")
		inputFileName  = flag.String("i", "", "Use a file with the name file-name as an input")
		sortByLine     = flag.Int("f", 0, "Sort input lines by value number N")
		outputFileName = flag.String("o", "", "Use a file with the name file-name as an output")
		revSort        = flag.Bool("r", false, "Sort input lines in reverse order")
	)
	flag.Parse()

	sigchnl := make(chan os.Signal, 1)

	filesChan := make(chan string)       
	isProcessed := make(chan struct{})  
	filesContent := make(chan string, 3) 
	buildTree := make(chan *tree)        

	signal.Notify(sigchnl, syscall.SIGINT)
	go func() {
		for {
			s := <-sigchnl
			handler(s, isProcessed, filesContent, filesChan)
		}
	}()
	//Stage one: Directory Reading
	go func() {
		if *inputFileName != "" {
			if *path != "." {
				log.Fatal("You can't use -i and -d options at the same time")
			}
			FFileD(".", filesChan, *inputFileName, abort) // знаходимо один заданий  файл
		} else {
			DFileR(*path, filesChan) // знаходимо усі файли .csv
		}
		close(filesChan)
	}()

	//Stage two: File Reading
	for i := 0; i < go_size; i++ {
		go func() {
			var line string
			var reader *bufio.Reader

			for path := range filesChan { 
				file, err := os.Open(path) 
				if err != nil {
					log.Fatal(err)
				}
				reader = bufio.NewReader(file)
				for {
					line, _ = reader.ReadString('\n') 
					line = strings.Trim(line, "\n")   
					if line == "" {                   
						break 
					}
					filesContent <- line 
				}
				file.Close() 
			}
			isProcessed <- struct{}{} 
		}()
	}
	//Stage three: Sorting
	go func() {
		var vertex *tree = createTreeV(<-filesContent)
		for cont := range filesContent {
			addBranch(vertex, cont, sortByLine)
		}
		buildTree <- vertex 
		close(buildTree)   
	}()
	for i := 0; i < go_size; i++ { 
		<-isProcessed
	}
	close(filesContent) 

	if *outputFileName != "" { 
		outFile, outErr := os.Create(*outputFileName) 

		if outErr != nil {
			log.Fatal(outErr)
		}
		defer outFile.Close() 
		writer := bufio.NewWriter(outFile)

		if !*revSort { 
			outTree(<-buildTree, writer)
		} else {
			outTreeR(<-buildTree, writer)
		}
		writer.Flush() 
		} else { 
		if !*revSort { 
			outCLI(<-buildTree)
		} else {
			outCLIRev(<-buildTree)
		}
	}

	go func() {
		exitchnl <- 0
		close(exitchnl)
	}()
	exitcode := <-exitchnl
}