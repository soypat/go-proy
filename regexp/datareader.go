package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)



type comision struct {
	 label string
	 schedule []string
	 teachers []string
}

func NewComision() comision {
	return comision{}
}
// Class as in school class
type Class struct {
	num1 int
	num2 int
	name string
	comisiones []comision
}
func main() {
	f,err := os.Open("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	line:=1

	reClassNumber := regexp.MustCompile( `[0-9]{2}\.[0-9]{2}`)
	reSchedule := regexp.MustCompile(`^[\s]{0,99}[A-Za-zéá]{5,9}[\s][0-9:0-9]{5}[\s-]{1,3}[0-9:0-9]{5}`)
	reComisionLabel := regexp.MustCompile(`(?:^[\s]{0,99})[A-Z]{1,8}(?:[\s]{0,99}$)`)
	reEndClass := regexp.MustCompile(`^[\s]{0,99}[0-9]{1,4}[\/\s]{1,3}[0-9]{1,4}[\s]{0,99}$`)
	var (
		allClasses []Class
		currentClass Class
		numberString string
		currentSchedule string
		)

	for scanner.Scan() {
		line++
		textLine:= scanner.Text()

		if len(textLine)==0 {
			continue
		}

		// Si encuentro una clase
		if reClassNumber.MatchString(textLine) {
			fmt.Printf("[DEBUG] Nueva class hallada (%d)\n",line)
			numberString = reClassNumber.FindString(textLine)
			currentClass.num1,err =strconv.Atoi(numberString[0:2])
			if err!=nil {
				break
			}
			currentClass.num2,err =strconv.Atoi(numberString[3:5])
			if err!=nil {
				break
			}
			currentClass.name = textLine[8:]

			currentComision := NewComision()
			// Entro en el for loop de las comisiones
			for scanner.Scan() {
				line++
				textLine = scanner.Text()

				// Si es una comision:
				if reComisionLabel.MatchString(textLine) {
					fmt.Printf(reComisionLabel.FindString(textLine)+ "\n\n")
					if currentComision.label!="" {
						fmt.Printf("[DEBUG] Comision %s append a class (%d)\n",currentComision.label,line)
						currentClass.comisiones = append(currentClass.comisiones,currentComision)
						currentComision = NewComision()
					}

					currentComision.label = reComisionLabel.FindString(textLine)

				}

				if   reEndClass.MatchString(textLine) {
					allClasses= append(allClasses, currentClass)
					fmt.Printf(textLine+"\n")
					fmt.Printf("[DEBUG] Fin de comisiones de una class. (%d)\n",line)
					break
				}

				currentSchedule = reSchedule.FindString(textLine)

				if currentSchedule!="" {
					currentComision.schedule = append(currentComision.schedule,currentSchedule)
				}
				if strings.Contains(textLine, ",") {
					currentComision.teachers = append(currentComision.teachers,textLine)
				}


			}
			fmt.Printf("%+v\n",currentComision)

		}

	}
	fmt.Printf("%+v\n",allClasses)
	if err!=nil {
		fmt.Printf("Hubo un error (%d): %s\n",line,err)
	}
}
