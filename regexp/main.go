package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const debug = false

type comision struct {
	label    string
	schedule []string
	teachers []string
}

func NewComision() comision {
	return comision{}
}

// Class as in school class
type Class struct {
	num1       int
	num2       int
	name       string
	comisiones []comision
}

type Schedule []comision

func NewSchedule() Schedule {
	return Schedule{}
}

func NewScheduleList() []Schedule {
	return []Schedule{}
}

func NewClass() Class {
	return Class{}
}

func main() {
	Classes, err := GatherClasses("data_smol.dat")
	if err != nil {
		panic("Big baddy")
	}

	//fmt.Printf("%+v", *Classes)
	ScheduleList := GatherSchedules(Classes)
	fmt.Printf("%+v",(*ScheduleList)[0])
}

func searcher(classes *[]Class, currentSchedule *Schedule, classNumber int) *[]Schedule {
	nextClass := (*classes)[classNumber]
	scheduleListMaster := NewScheduleList()
	for _, v := range nextClass.comisiones {
		scheduleInstance := append(*currentSchedule, v)

		if classNumber == len(*classes)-1 { //llegue a la ultima clase
			isValid := verifySchedule(&scheduleInstance)
			if isValid { //El schedule es bueno, lo devuelvo como lista no nula
				scheduleListMaster = append(scheduleListMaster, scheduleInstance)
				continue
			} else {
				continue
			} //Return ends

		} else { // Si no es la ultima clase, sigo por aca
			scheduleList := searcher(classes, &scheduleInstance, classNumber+1) // Awesome recursion baby
			if *scheduleList == nil {
				continue
			}
			scheduleListMaster = append(scheduleListMaster, *scheduleList...)
		}
	}
	return &scheduleListMaster
}



func verifySchedule(currentSchedule *Schedule) bool {
	// TODO hard part coming ahead. Actual verification
	numberOfMaterias := len(*currentSchedule)
	reWeek := regexp.MustCompile(`(?i)Lunes|Martes|Miercoles|Jueves|Viernes|Sabado|Domingo`)
	for i:=0; i<numberOfMaterias-1;i++  {
		firstComision := (*currentSchedule)[i]
		for j:=numberOfMaterias-1;j>i;j-- {
			secondComision := (*currentSchedule)[j]
			firstSchedule := firstComision.schedule
			secondSchedule := secondComision.schedule
			matches := reWeek.FindAllStringIndex(firstSchedule, -1) //TODO see TODO for schedule sniffing in GatherClasses
			fmt.Printf("%s -- %s\n\n",firstSchedule,secondSchedule)
		}
	}
	return true
}

func GatherSchedules(classes *[]Class) *[]Schedule {
	//numberOfClasses := len(*classes)
	//verifiedScheduleNumber := 0
	//scheduleListMaster := NewScheduleList()
	currentSchedule := NewSchedule()
	scheduleListMaster := searcher(classes, &currentSchedule, 0)

	// Search function: cada instancia de recursividad busca verificar un schedule (lo va fabricando a medida que avanza) y devuelve una lista de schedules verificados y los va juntando en cada instancia

	return scheduleListMaster
}

func GatherClasses(filedir string) (*[]Class, error) {
	f, err := os.Open(filedir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	line := 0

	reClassNumber := regexp.MustCompile(`[0-9]{2}\.[0-9]{2}`)
	reSchedule := regexp.MustCompile(`^[\s]{0,99}[A-Za-zéá]{5,9}[\s][0-9:0-9]{5}[\s-]{1,3}[0-9:0-9]{5}`)
	reComisionLabel := regexp.MustCompile(`(?:^[\s]{0,99})[A-Z]{1,8}(?:[\s]{0,99}$)`)
	reEndComision := regexp.MustCompile(`^[\s]{0,99}[0-9]{1,4}[\/\s]{1,3}[0-9]{1,4}[\s]{0,99}$`)
	var (
		currentClass    Class
		allClasses      []Class
		numberString    string
		currentSchedule string
	)

	for scanner.Scan() {
		line++
		textLine := scanner.Text()

		if len(textLine) == 0 {
			continue
		}

		// Si encuentro una clase
		for reClassNumber.MatchString(textLine) {
			if debug {
				fmt.Printf("[DEBUG] Nueva class hallada (%d)\n", line)
			}
			currentClass = NewClass()
			numberString = reClassNumber.FindString(textLine)
			currentClass.num1, err = strconv.Atoi(numberString[0:2])
			if err != nil {
				break
			}
			currentClass.num2, err = strconv.Atoi(numberString[3:5])
			if err != nil {
				break
			}
			currentClass.name = textLine[8:]

			currentComision := NewComision()
			// Entro en el for loop de las comisiones
			for scanner.Scan() {
				line++
				textLine = scanner.Text()
				if reClassNumber.MatchString(textLine) {
					allClasses = append(allClasses, currentClass)
					if debug {
						fmt.Printf("[DEBUG] Fin de class y comienzo de otra hallada (%d)\n", line)
					}
					break
				}
				// Si es una comision:
				if reComisionLabel.MatchString(textLine) {
					if debug {
						fmt.Printf("[DEBUG] Nueva Comision %s encontrada (%d)\n", textLine, line)
					}
					if currentComision.label != "" {
						if debug {
							fmt.Printf("[DEBUG] Comision %s append a class (%d)\n", currentComision.label, line)
						}
						currentClass.comisiones = append(currentClass.comisiones, currentComision)
						currentComision = NewComision()
					}

					currentComision.label = reComisionLabel.FindString(textLine)

				}

				if reEndComision.MatchString(textLine) {
					currentClass.comisiones = append(currentClass.comisiones, currentComision)
					if debug {
						fmt.Printf("[DEBUG] Fin de una comision. (%d)\n", line)
					}
					continue
				}

				currentSchedule = reSchedule.FindString(textLine)

				if currentSchedule != "" { //TODO Possible bugfix. All schedules for one comision are put in same string. Better to split.
					currentComision.schedule = append(currentComision.schedule, currentSchedule)
				}
				if strings.Contains(textLine, ",") {
					currentComision.teachers = append(currentComision.teachers, textLine)
				}
			}
		}
	}
	if debug {
		fmt.Printf("[DEBUG] Se termino de buscar Class. GatherClass Over (%d)\n", line)
	}
	allClasses = append(allClasses, currentClass)
	if err != nil {
		err = fmt.Errorf("Hubo un error (%d): %s\n", line, err)
	}
	return &allClasses, err
}
