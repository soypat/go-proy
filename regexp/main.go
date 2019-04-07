package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const debug = true

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

type Schedule []Class

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

	fmt.Printf("%+v", *Classes)
	ScheduleList := GatherSchedules(Classes)
	fmt.Printf("%+v", ScheduleList)
}

func searcher(classes *[]Class, currentSchedule *Schedule, classNumber int) *[]Schedule{
	nextClass := (*classes)[classNumber]
	scheduleListMaster := NewScheduleList()
	for _, v := range nextClass.comisiones {
		if classNumber == len(*classes)-1 {
			isValid := verifySchedule(currentSchedule)
		} else {
			scheduleList := searcher(classes,currentSchedule,1)
			if *scheduleList == nil { //TODO Make sure to dereference all schedule pointers during work
				continue
			}
			scheduleListMaster = append(scheduleListMaster,*scheduleList...)
		}
	}

}
func verifySchedule(currentSchedule *Schedule) bool {
	// TODO hard part coming ahead. Actual verification
}

func GatherSchedules(classes *[]Class) []Schedule {
	numberOfClasses := len(*classes)
	verifiedScheduleNumber := 0
	nextClass := (*classes)[0]
	scheduleListMaster := NewScheduleList()
	for _, v := range nextClass.comisiones {//TODO Include this in recursive function
		currentSchedule := NewSchedule()
		scheduleList := searcher(classes,&currentSchedule,1)
		scheduleListMaster = append(scheduleListMaster,*scheduleList...)
	}

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

				if currentSchedule != "" {
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
