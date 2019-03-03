// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
package main

import (
	"bufio"
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)


func main() {
	files, err := fileListCurrentDirectory()

	if err!=nil {
		log.Fatalf("Failed getting Current directory: %v",err)
	}
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	filelist := widgets.NewList()
	filelist.Title = "Archivos disponibles"
	filelist.SetRect(0, 0, 40, 12)
	filelist.Rows = files
	filelist.TextStyle = ui.NewStyle(ui.ColorYellow)
	filelist.WrapText = false

	headline := widgets.NewParagraph()
	headline.Border = false
	headline.Text = "Presione (q) para salir."
	headline.SetRect(0, 13, 60, 16)

	ui.Render(headline)
	ui.Render(filelist)
	status := widgets.NewParagraph()
	status.Title = "Status del programa"
	status.Text = "Programa Iniciado. Seleccionar archivo NASTRAN DECK."
	status.SetRect(40, 0, 60, 12)
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
			filelist.HalfPageDown()
		case "<C-u>":
			filelist.HalfPageUp()
		case "<C-f>":
			filelist.PageDown()
		case "<C-b>":
			filelist.PageUp()
		case "g":
			if previousKey == "g" {
				filelist.ScrollTop()
			}
		case "<Home>":
			filelist.ScrollTop()
		case "G", "<End>":
			filelist.ScrollBottom()
		case "<Enter>":
			filedir := files[filelist.SelectedRow]
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
		// Little paragraph
		ui.Render(status)
		ui.Render(filelist)
	}
	//q := widgets.NewParagraph()
	//q.Text = "another one!"
	//q.SetRect(25, 0, 50, 5)

}

func fileListCurrentDirectory() ([]string,error) {
	var files []string
	root,err:=filepath.Abs("./")
	if err!=nil{
		return nil,err
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		files = append(files,path)
		return nil
	})
	if err != nil {
		return nil,err
	}
	// Ahora lo que hago es excluir la parte reduntante del dir
	// C:/Go/mydir/foo/myfile.exe  ----> se convierte a ---> foo/myfile.exe
	//const numberOfFiles
	var fileLength int
	maxFileLength  := 0
	minFileLength := 2047
	i:= 0
	//imin:=-1
	//imax:=-1
	var newfiles = files[0:len(files)-1]

	for _,file := range files {
		fileLength = len(file)
		if fileLength>maxFileLength {
			maxFileLength = fileLength
			//imax = i
		}
		if fileLength<minFileLength {
			minFileLength = fileLength
			//imin = i
		}
		i++
	}
	i=0
	//stringer:= strconv.Itoa(minFileLength)
	for _,file := range files {
		if len(file) <= minFileLength {
			continue
		}
		newfiles[i] ="~"+ file[minFileLength:]
		i++
	}
	return newfiles,nil
	}

func NasRead(filedir string) error {
	f, err :=os.Open(filedir)
	if err!= nil {
		return err
	}
	defer f.Close()
	//enterContinue("Se encontró datos.dat!\n")
	scanner := bufio.NewScanner(f)
	line:=1
	for scanner.Scan() {
		line++
		if strings.Contains(scanner.Text(), "GRID CARDS") {
			scanner.Scan()
			line++
			break
		}
	}
	var nodeNum int
	var nodex,nodey,nodez float64
	var mantx,manty,mantz float64
	var expx,expy,expz int
	NodesFound:=false
	ElementsFound := false

	d, err := os.Create("nodos.txt")
	if err!=nil {
		return err
	}
	writer := bufio.NewWriter(d)
	spacer := "\t"


	for scanner.Scan() { // BUSQUEDA DE NODOS
		if strings.Contains(scanner.Text(),"$") {
			break
		}
		line1:=strings.Fields(scanner.Text())
		line++
		scanner.Scan()
		line++
		line2:=strings.Fields(scanner.Text())
		nodeNum,err = strconv.Atoi(line1[1])
		if err!=nil {
			return err
		}
		chunk := line1[2] // el chunk es un string con dos numeros (x, y)
		chunkz:=line2[1]

		mantx, err = strconv.ParseFloat(chunk[1:13], 64)
		if err!=nil {
			return err
		}
		expx,err=strconv.Atoi(chunk[14:17])
		if err!=nil {
			return err
		}
		manty,err=strconv.ParseFloat(chunk[17:29],64)
		if err!=nil {
			return err
		}
		expy,err=strconv.Atoi(chunk[30:33])
		if err!=nil {
			return err
		}
		mantz,err=strconv.ParseFloat( chunkz[0:12],64)
		if err!=nil {
			return err
		}
		expz,err=strconv.Atoi(chunkz[13:16])
		if err!=nil {
			return err
		}
		nodex = mantx*math.Pow10(expx)
		nodey = manty*math.Pow10(expy)
		nodez = mantz*math.Pow10(expz)


		_,err = writer.WriteString(strconv.Itoa(nodeNum))

		if err!=nil {
			return err
		}
		_,err = writer.WriteString(spacer)
		if err!=nil {
			return err
		}
		_,err = writer.WriteString(fmt.Sprintf("%e", nodex))
		if err!=nil {
			return err
		}
		_,_ = writer.WriteString(spacer)
		_,err = writer.WriteString(fmt.Sprintf("%e", nodey))
		if err!=nil {
			return err
		}
		_,_ = writer.WriteString(spacer)
		_,err = writer.WriteString(fmt.Sprintf("%e", nodez))
		if err!=nil {
			return err
		}
		_,err = writer.WriteString("\r\n")
		writer.Flush()

		if err!=nil {
			return err
		}

		if math.Mod(float64(line),1000)==0 {
			d.Sync()
		}
		NodesFound = true
	}
	d.Sync()
	d.Close()

	d2, err := os.Create("elementos.txt")
	if err!=nil {
		return err
	}
	writer2 := bufio.NewWriter(d2)
	defer d2.Close()
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	var elementNum string
	for scanner.Scan() {
		line++
		if math.Mod(float64(line),30)==0 {
			d2.Sync()
		}
		if strings.Contains(scanner.Text(), "CHEXA") {
			line1:=strings.Fields(scanner.Text())
			line++
			scanner.Scan()
			line2:=strings.Fields(scanner.Text())
			line++
			scanner.Scan()
			line3:=strings.Fields(scanner.Text())
			line1[8]= reg.ReplaceAllString(line1[8], "")
			line2[8] = reg.ReplaceAllString(line2[8], "")
			elementNum = line1[1]
			writer2.Flush()
			_,err = writer2.WriteString(elementNum)
			//fmt.Println(elementNum)
			//enterContinue("Did i print elementNum?")
			if err!=nil {
				return err
			}
			for i:=3;i<9;i++ {
				writer2.Flush()
				_,err = writer2.WriteString(spacer)
				_, err = writer2.WriteString(line1[i])
				//fmt.Print(line1[i]," ")
				if err!=nil {
					return err
				}
			}
			for i:=1;i<9;i++ {
				writer2.Flush()
				_,err = writer2.WriteString(spacer)
				_, err = writer2.WriteString(line2[i])
				//fmt.Print(line2[i]," ")
				if err!=nil {
					return err
				}
			}
			for i:=1;i<7;i++ {
				writer2.Flush()
				_,err = writer2.WriteString(spacer)
				_, err = writer2.WriteString(line3[i])
				//fmt.Printf(line3[i]," ")
				if err!=nil {
					return err
				}
			}

			_,err = writer2.WriteString("\r\n")
			writer2.Flush()
			ElementsFound=true
		}
		if strings.Contains(scanner.Text(), "MATERIAL CARDS") {
			break
		}

		//enterContinue("ELEMENT FOUND:")
		//fmt.Printf(line1)
	}
	writer2.Flush()
	d2.Sync()
	if ElementsFound && NodesFound {
		return nil
	} else {
		return fmt.Errorf("No se encontraron nodos o elementos.")
	}

}