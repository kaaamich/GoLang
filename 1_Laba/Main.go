package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type node struct { // data structure for the iterative method
	data []string // data
	next *node    // addresses of the following nodes
}

type tree struct { // tree structure
	data  []string // data
	left  *tree    // left branch
	right *tree    // right branch
}

//Binary search tree functions
func createTV(buffer *string) *tree { // creation of the top of the branch
	var vertex = new(tree)
	vertex.data = strings.Split(*buffer, ";") // filling the vertex with the initial value
	return vertex
}
func addB(vertex *tree, buffer *string, sortByLine *int) { // add a new branch
	var compare = strings.Split(*buffer, ";")
	if compare[*sortByLine] > vertex.data[*sortByLine] {
		if vertex.right == nil {
			vertex.right = new(tree)
			vertex.right.data = compare
		} else {
			addB(vertex.right, buffer, sortByLine)
		}
	} else if compare[*sortByLine] < vertex.data[*sortByLine] {
		if vertex.left == nil {
			vertex.left = new(tree)
			vertex.left.data = compare
		} else {
			addB(vertex.left, buffer, sortByLine)
		}
	} else if compare[*sortByLine] == vertex.data[*sortByLine] {
		if vertex.left == nil {
			vertex.left = new(tree)
			vertex.left.data = compare
		} else {
			insertE(vertex, compare)
		}
	}
}
func insertE(vertex *tree, data []string) { // add element between branches
	temp := vertex.left
	vertex.left = new(tree)
	vertex.left.data = data
	vertex.left.left = temp
}
func outT(vertex *tree, file *os.File) { // recursively output the value in ascending order
	if vertex.left != nil {
		outT(vertex.left, file)
	}
	writeFile(file, vertex)
	if vertex.right != nil {
		outT(vertex.right, file)
	}
}
func outTR(vertex *tree, file *os.File) { // recursively output values ​​in descending order
	if vertex.right != nil {
		outTR(vertex.right, file)
	}
	writeFile(file, vertex)
	if vertex.left != nil {
		outTR(vertex.left, file)
	}
}
func writeFile(file *os.File, temp *tree) { // output the data array of each branch
	for i := 0; i < len(temp.data); i++ { // output the array of substrings to a file
		file.WriteString(temp.data[i] + ";")
	}
	file.WriteString("\n") // new line in file
}

//Sorting functions using a singly linked list (nodes)
func createNodeHeader(buffer *string) node { // create the beginning of the list
	var _startCell node
	_startCell.data = strings.Split(*buffer, ";") // add a node
	_startCell.next = nil
	return _startCell
}
func addN(temp **node) { // додаємо вузол
	(*temp).next = new(node)
	(*temp) = (*temp).next
	(*temp).next = nil
}
func headLineOptionSet(temp **node, _startCell *node, file *os.File, i *int) { // set header option(-h)
	(*temp) = _startCell.next
	writeIn(file, _startCell)
	(*i)--
}
func nodeBegin(temp **node, _startCell *node) { // return the list to the beginning
	(*temp) = _startCell
}
func nextNode(temp **node, _startCell *node, headOp *bool) { // switch the node depending on the status of the selected option
	if *headOp {
		if (*temp).next == nil { // we make the following conditions until we find all the elements
			(*temp) = _startCell.next
		} else {
			(*temp) = (*temp).next
		}
	} else {
		if (*temp).next == nil { // we make the following conditions until we find all the elements
			(*temp) = _startCell
		} else {
			(*temp) = (*temp).next
		}
	}
}

func writeIn(file *os.File, temp *node) { // output the array of substrings to a file
	for i := 0; i < len(temp.data); i++ {
		file.WriteString(temp.data[i] + ";")
	}
	file.WriteString("\n")
}
func sortUp(temp *node, counter int, outFile *os.File, headOp *bool, _startCell *node, sortByLine *int, str []string) { // sorting in ascending order
	if *headOp {
		temp = _startCell.next
	} else {
		temp = _startCell
	}
	for i := 0; i < counter; {
		if temp.data[*sortByLine] == str[i] {
			writeIn(outFile, temp)
			i++
		}
		nextNode(&temp, _startCell, headOp)
	}
}
func sortRev(temp *node, counter int, outFile *os.File, headOp *bool, _startCell *node, sortByLine *int, str []string) { // sort in descending order
	if *headOp { // assign to temp the address of the node of the list
		temp = _startCell.next
	} else {
		temp = _startCell
	}
	for i := counter - 1; i >= 0; {
		if temp.data[*sortByLine] == str[i] { // search for the i element of the sorted array in the list (sorting method)
			writeIn(outFile, temp)
			i--
		}
		nextNode(&temp, _startCell, headOp)
	}
}
func main() {
	var (
		_startCell node                 // start list node
		temp       *node  = &_startCell // variable to store the address of the temporary node
		buffer     string               // buffer variable for the entered lines
		counter    int    = 1           // counter of the number of elements of a singly linked list
	)
	var (
		inputFileName  = flag.String("i", "input.csv", "Use a file with the name file-name as an input")
		outputFileName = flag.String("o", "output.csv", "Use a file with the name file-name as an output")
		headOp         = flag.Bool("h", true, "The first line is a header that must be ignored during sorting but included in the output")
		sortByLine     = flag.Int("f", 0, "Sort input lines by value number N")
		revSort        = flag.Bool("r", false, "Sort input lines in reverse order")
		treeSort       = flag.Int("a", 1, "Sorty by tree or default algorithm")
	)
	flag.Parse()

	inpFile, inpErr := os.Create(*inputFileName)  // create a file for input
	outFile, outErr := os.Create(*outputFileName) // create a file for output

	if inpErr != nil { // if we received an error, we terminate the program
		fmt.Println("Unable to create input file", inpErr)
		os.Exit(1)
	} else if outErr != nil { // if we received an error, we terminate the program
		fmt.Println("Unable to create output file", outErr)
		os.Exit(1)
	}
	defer inpFile.Close() // finish working with the file
	defer outFile.Close() // finish working with the file

	fmt.Println("Input CSV data line by line:")
	n, _ := fmt.Fscanln(os.Stdin, &buffer) // enter a line, save it in a variable
	inpFile.WriteString(buffer + "\n")

	switch *treeSort { // depending on the selected type of sorting (-a flag)

	//SORTING BY EXCESS
	case 1:
		if n != 0 { //if string is NOT empty
			_startCell = createNodeHeader(&buffer)
		} else { // otherwise, we close the program with a message about the absence of entered data
			fmt.Println("You input no data")
			os.Exit(1)
		}

		for n, _ = fmt.Fscanln(os.Stdin, &buffer); n != 0; counter++ { // create new list elements (nodes)
			addN(&temp)
			temp.data = strings.Split(buffer, ";") // fill the nodes with the entered values ​​in the console
			writeIn(inpFile, temp)
			n, _ = fmt.Fscanln(os.Stdin, &buffer) // scan the next entered line
		}

		if *headOp { // enable the -h header option
			headLineOptionSet(&temp, &_startCell, outFile, &counter)
		} else { // otherwise, we return the list to the beginning
			nodeBegin(&temp, &_startCell)
		}

		str := make([]string, counter) // we create an array based on the number of elements of a singly linked list
		for i := 0; i < counter; i++ { // we fill the array with the first values ​​of the array of substrings from each element of the list
			str[i] = temp.data[*sortByLine]
			temp = temp.next // scroll through the list
		}
		sort.Strings(str) // we sort the elements

		switch *revSort { // sort in ascending or descending order depending on the option
		case true:
			sortRev(temp, counter, outFile, headOp, &_startCell, sortByLine, str)
		case false:
			sortUp(temp, counter, outFile, headOp, &_startCell, sortByLine, str)
		}

	//SORTING BY TREE
	case 2:
		if n != 0 { // if string is NOT empty
			var vertex *tree
			if *headOp {
				outFile.WriteString(buffer + "\n")
				n, _ = fmt.Fscanln(os.Stdin, &buffer)
				inpFile.WriteString(buffer + "\n")
				vertex = createTV(&buffer)
			} else {
				vertex = createTV(&buffer)
			}
			for n, _ = fmt.Fscanln(os.Stdin, &buffer); n != 0; { // create new list elements (nodes)
				inpFile.WriteString(buffer + "\n")
				addB(vertex, &buffer, sortByLine)
				n, _ = fmt.Fscanln(os.Stdin, &buffer) // scan the next entered line
			}
			if *revSort {
				outTR(vertex, outFile)
			} else {
				outT(vertex, outFile)
			}
		} else { // otherwise, we close the program with a message about the absence of entered data
			fmt.Println("You input no data")
			os.Exit(1)
		}
	}
}