// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
//+build termuiver

package main

import (
	"bufio"
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)


func main() {
	displayedFileNames,fileNames, err := fileListCurrentDirectory()


	if err!=nil {
		log.Fatalf("Failed getting Current directory: %v",err)
	}
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	filelist := widgets.NewList()
	filelist.Title = "Archivos disponibles"
	filelist.SetRect(0, 0, 60, 12)
	filelist.Rows = displayedFileNames
	filelist.TextStyle = ui.NewStyle(ui.ColorYellow)
	filelist.WrapText = false

	headline := widgets.NewParagraph()
	headline.Border = false
	headline.Text = "Presione (q) para salir."
	headline.SetRect(0, 12, 60, 15)

	ui.Render(headline)
	ui.Render(filelist)
	status := widgets.NewParagraph()
	status.Title = "Status del programa"
	status.Text = `Programa Iniciado. Seleccionar archivo NASTRAN DECK.

Programado en Go.
Patricio Whittingslow 2019`
	status.SetRect(60, 0, 80, 12)
	ui.Render(status)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			filelist.ScrollDown()
		case "k", "<Up>":
			filelist.ScrollUp()
		case "<C-d>":
			filelist.ScrollHalfPageDown()
		case "<C-u>":
			filelist.ScrollHalfPageUp()
		case "<C-f>":
			filelist.ScrollPageDown()
		case "<C-b>":
			filelist.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				filelist.ScrollTop()
			}
		case "<Home>":
			filelist.ScrollTop()
		case "G", "<End>":
			filelist.ScrollBottom()
		case "<Enter>":
			filedir := fileNames[filelist.SelectedRow]
			status.Text = filedir+ "\nArchivo Seleccionado! Espere por favor..."
			status.TextStyle = ui.NewStyle(ui.ColorWhite)
			ui.Render(status)
			err := NasRead(filedir[2:])
			if err!=nil {
				status.Text = fmt.Sprintf("Ups, hubo un error: %s",err)
				status.TextStyle = ui.NewStyle(ui.ColorRed)
			} else {

				status.Text = "Archivo "+ filedir +" leido. Ver archivos nodos.txt y elementos.txt"
				status.TextStyle = ui.NewStyle(ui.ColorGreen)
			}

		}
		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
		ui.Render(status)
		ui.Render(filelist)
	}


}

func fileListCurrentDirectory() ([]string,[]string,error) {
	var files []string
	root,err:=filepath.Abs("./")
	if err!=nil{
		return nil,nil,err
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		files = append(files,path)
		return nil
	})
	if err != nil {
		return nil,nil,err
	}
	// Ahora lo que hago es excluir la parte reduntante del dir
	// C:/Go/mydir/foo/myfile.exe  ----> se convierte a ---> foo/myfile.exe
	//const numberOfFiles
	var fileLength int
	maxFileLength  := 0
	minFileLength := 2047

	i:= 0
	var shortFileNames,actualFileNames []string
	shortFileNames = append(files[:0:0], files...)
	actualFileNames = append(files[:0:0], files...)

	for _,file := range files {
		fileLength = len(file)
		if fileLength>maxFileLength {
			maxFileLength = fileLength
		}
		if fileLength<minFileLength {
			minFileLength = fileLength
		}
		i++
	}
	permittedStringLength := 54
	i=0

	for _,file := range files {
		if len(file) <= minFileLength {
			files = remove(files,i)
			shortFileNames = remove(shortFileNames,i)
			actualFileNames = remove(actualFileNames,i)
			continue
		}
		if len(file)>permittedStringLength+minFileLength {

			shortFileNames[i] = `~\…`+ file[len(file)-permittedStringLength:]

		} else {
			shortFileNames[i] ="~"+ file[minFileLength:]

		}
		actualFileNames[i] ="~"+ file[minFileLength:]
		i++
	}
	return shortFileNames,actualFileNames,nil
	}

const whitespace string = "\n\r\t "
const spacedInteger string = "%d\t"
const spacedDimension string = "%e\t%e\t%e\t%e\r\n"

type dims struct {
	x float64
	y float64
	z float64
	t float64
}

type node struct {
	number int
	dims
}

type element struct {
	number    int
	nodeIndex []int
	Type      string

	group int // No creo que lo usaría
}

func assignElement(element *element, integerString []string) (err error) {
	// El primer integer es el NUMERO del elemento, no es un nodo... y el segundo es el grupo del elemento!
	element.number, err = strconv.Atoi(strings.TrimLeft(integerString[0], whitespace))
	if err != nil {
		return err
	}
	for _, currentString := range integerString[2:] {
		integer, err := strconv.Atoi(strings.TrimLeft(currentString, whitespace))
		if err != nil {
			return err
		}
		element.nodeIndex = append(element.nodeIndex, integer)

	}
	return nil
}
func assignDims(dimensions *dims, floatString []string) {
	for i := range floatString {
		switch i {
		case 0:
			dimensions.x, _ = strconv.ParseFloat(floatString[i], 64)
		case 1:
			dimensions.y, _ = strconv.ParseFloat(floatString[i], 64)
		case 2:
			dimensions.z, _ = strconv.ParseFloat(floatString[i], 64)
		case 3:
			dimensions.t, _ = strconv.ParseFloat(floatString[i], 64)
		default:
			panic("Unreachable")
		}
	}
}

func writeNode(node *node, writer *bufio.Writer) (err error) {
	_, err = writer.WriteString(fmt.Sprintf(spacedInteger, node.number))
	if err != nil {
		return err
	}
	_, err = writer.WriteString(fmt.Sprintf(spacedDimension, node.x, node.y, node.z, node.t))
	if err != nil {
		return err
	} else {
		writer.Flush() // Let it be, mother mary say to me
		return nil
	}


}

func writeElement(element *element, writer *bufio.Writer) (err error) {
	_, err = writer.WriteString(fmt.Sprintf(spacedInteger, element.number))
	for _, v := range element.nodeIndex {
		if v == element.nodeIndex[len(element.nodeIndex)-1] {
			_, err = writer.WriteString(fmt.Sprintf("%d\r\n", v))
		} else {
			_, err = writer.WriteString(fmt.Sprintf(spacedInteger, v))
		}
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func NasRead(filedir string) error {

	dataFile, err := os.Open(filedir)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	nodeFile, err := os.Create("nodos.txt")
	if err != nil {
		return err
	}
	nodeWriter := bufio.NewWriter(nodeFile)

	elementFile, err := os.Create("elementos.txt")
	if err != nil {
		return err
	}
	defer elementFile.Sync()
	defer nodeFile.Sync()
	defer elementFile.Close()
	defer nodeFile.Close()

	elementWriter := bufio.NewWriter(elementFile)

	reGridStart := regexp.MustCompile(`GRID\*`)
	reElementStart := regexp.MustCompile(`[A-Z]{2,7}[\s\d]+[\+\n]`)
	reElementType := regexp.MustCompile(`^[A-Z]{2,7}`)
	reInteger := regexp.MustCompile(`\s+\d+`)
	reNodeNumber := regexp.MustCompile(`(?:GRID\*\s+)([\d]+)`)
	reNonNumerical := regexp.MustCompile(`[A-Za-z\*\+\-\s,]+`)
	reFloat := regexp.MustCompile(`\d{1}\.\d{4,16}E{1}[\+\-]{1}\d{2}`)
	// Nastran Decks tienen el flag "+" para indicar que la informacion del objeto sigue en la proxima linea
	reLineContinueFlag := regexp.MustCompile(`\+{1}\n*$`)
	scanner := bufio.NewScanner(dataFile)
	line := 1
	var nodeNumberString, currentText string
	var floatStrings, integerStrings []string

	for scanner.Scan() {
		line++

		if reGridStart.MatchString(scanner.Text()) {
			var currentNode node

			currentText = reLineContinueFlag.ReplaceAllString(scanner.Text(), "")
			for reLineContinueFlag.MatchString(scanner.Text()) {
				scanner.Scan()
				line++
				currentText = currentText + scanner.Text()
			}
			nodeNumberString = reNodeNumber.FindString(currentText)
			currentNode.number, err = strconv.Atoi(reNonNumerical.ReplaceAllString(nodeNumberString, ""))

			floatStrings = reFloat.FindAllString(currentText, 4)
			assignDims(&currentNode.dims, floatStrings)
			err = writeNode(&currentNode, nodeWriter)
			if err != nil {
				return err
			} else {
				continue
			}
		}
		if reElementStart.MatchString(scanner.Text()) {
			var currentElement element
			currentText = reLineContinueFlag.ReplaceAllString(scanner.Text(), "")
			for reLineContinueFlag.MatchString(scanner.Text()) {
				scanner.Scan()
				line++
				currentText = currentText + scanner.Text()
			}
			currentElement.Type = reElementType.FindString(currentText)
			integerStrings = reInteger.FindAllString(currentText, -1)
			err = assignElement(&currentElement, integerStrings)
			if err != nil {
				return err
			}
			err = writeElement(&currentElement, elementWriter)
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}