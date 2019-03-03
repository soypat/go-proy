package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"math"
	"regexp"
)

func userInput(m string) string {
	var out string
	fmt.Println(m)
	reader := bufio.NewReader(os.Stdin)
	out, _ = reader.ReadString('\n')
	return out
}


func NasRead(filedir string) error {
	//fmt.Println("Patricio Whittingslow's miraculous .dat reader\n")
	//fmt.Println("Programado con Go usando GoLand.\n\n")
	//fmt.Println("Guardar Nastran Deck como datos.dat!\n")
	//enterContinue("")

	//eleType:=userInput("Que elementos quiere buscar?")
	//if eleType=="" {
	//	eleType="CHEXA"
	//}

	f, err :=os.Open(filedir)
	if err!= nil {
		//enterContinue( fmt.Sprintf("No se encontró %s\nEl programa cerrará\n\n",filedir))
		//panic(err)
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

		//fmt.Sprintf("%e", nodex)
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
		//nodeline := fmt.Sprintf(strconv.Itoa(nodeNum),spacer,fmt.Sprintf("%e", nodex),spacer,fmt.Sprintf("%e", nodey),spacer,fmt.Sprintf("%e", nodez),"\n")
		//_, err := writer.WriteString(nodeline)
		if err!=nil {
			return err
		}
		fmt.Println(nodeNum,spacer,nodex,spacer,nodey,spacer,nodez)
		if math.Mod(float64(line),1000)==0 {
			d.Sync()
		}




		//fmt.Println("Node#: ",nodeNum)
		//fmt.Println(nodex)
		//fmt.Println(nodey)
		//fmt.Println(nodez)
		//fmt.Println("\n")
		//fmt.Println(chunk)
		//fmt.Println(chunkz)
		//enterContinue("found chunks")
	}
	d.Sync()
	d.Close()

	//fmt.Println("Nodos escritos. Buscando Elementos tipo CHEXA por defecto (20 nodos).")
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
		}
		if strings.Contains(scanner.Text(), "MATERIAL CARDS") {
			break
		}

		//enterContinue("ELEMENT FOUND:")
		//fmt.Printf(line1)
	}
	writer2.Flush()
	d2.Sync()
	return nil
}





