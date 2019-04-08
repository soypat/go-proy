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

var Days = [...]string{"Lunes", "Martes", "Miercoles", "Jueves", "Viernes", "Sabado", "Domingo"}

type comision struct {
	label     string
	schedules []schedule
	teachers  []string
}
type schedule struct {
	//TODO rewrite schedule/horario/cursada. Confusing structs.
	day   int // from 0 to 6
	start time
	end   time
}

func NewSchedule() schedule {
	return schedule{}
}

type time [2]int

func NewTime() time {
	return time{}
}
func NewComision() comision {
	return comision{}
}

type scheduleCriteria struct {
	maxSuperposition float32
	freeDays         [len(Days)]bool
}

// Class as in school class
type Class struct {
	num1       int
	num2       int
	name       string
	comisiones []comision
}

type Cursada []comision // Cursada is a the group of courses a student attends during the year/semester
// TODO Solve issue if CURSADA should have flag or if search should recieve criteria for schedule
func NewCursada() Cursada {
	return Cursada{}
}

func NewCursadaList() []Cursada {
	return []Cursada{}
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
	fmt.Printf("\n\n%+v", (*ScheduleList)[0])
}

func searcher(classes *[]Class, currentCursada *Cursada, classNumber int) *[]Cursada {
	nextClass := (*classes)[classNumber]
	cursadaListMaster := NewCursadaList()
	for _, v := range nextClass.comisiones {
		cursadaInstance := append(*currentCursada, v)

		if classNumber == len(*classes)-1 { //llegue a la ultima clase
			isValid := verifyCursada(&cursadaInstance)
			if isValid { //El schedule es bueno, lo devuelvo como lista no nula
				cursadaListMaster = append(cursadaListMaster, cursadaInstance)
				continue
			} else {
				continue
			} //Return ends

		} else { // Si no es la ultima clase, sigo por aca
			cursadaList := searcher(classes, &cursadaInstance, classNumber+1) // Awesome recursion baby
			if *cursadaList == nil {
				continue
			}
			cursadaListMaster = append(cursadaListMaster, *cursadaList...)
		}
	}
	return &cursadaListMaster
}

func verifyCursada(currentCursada *Cursada) bool {
	// TODO hard part coming ahead. Actual verification
	numberOfMaterias := len(*currentCursada)

	for i := 0; i < numberOfMaterias-1; i++ {
		firstComision := (*currentCursada)[i]
		for j := numberOfMaterias - 1; j > i; j-- {
			secondComision := (*currentCursada)[j]
			firstCursada := firstComision.schedules
			secondCursada := secondComision.schedules
			for _, schedule1 := range firstCursada { //TODO change iteration after storing horario correctly
				for _, schedule2 := range secondCursada {
					//matches1 := reWeek.FindAllString(horario1,-1)
					//matches2 := reWeek.FindAllString(horario2,-1)
					//fmt.Printf("%s -- %s",matches1,matches2)
					if schedule1.day != schedule2.day { // si no coinciden los dias, verifica ese horario, continuo buscando colisiones
						continue
					} else { // en el caso que sean el mismo día:

					}
				}
			}
			//TODO see TODO for schedule sniffing in GatherClasses
			//fmt.Printf("%s -- %s\n\n", firstSchedule, secondCursada)
		}
	}
	return true
}

func stringToTime(scheduleString string) (int, []time, error) {
	reWeek := regexp.MustCompile(`(?i)Lunes|Martes|Miercoles|Jueves|Viernes|Sabado|Domingo`)
	reSchedule := regexp.MustCompile(`[0-9:0-9]{5}[\s-]{1,5}[0-9:0-9]{5}`)
	reScheduleStart := regexp.MustCompile(`^[0-9:0-9]{5}`)
	reScheduleFinish := regexp.MustCompile(`[0-9:0-9]{5}$`)
	reHours := regexp.MustCompile(`^[0-9]{2}`)
	reMinutes := regexp.MustCompile(`[0-9]{2}$`)
	diaString := reWeek.FindString(scheduleString)
	diaInt := -1
	for i, v := range Days {
		if strings.Contains(strings.Title(diaString), v) {
			diaInt = i
		}
	}
	if diaInt == -1 {
		return diaInt, nil, fmt.Errorf("Failed to match string  %s  with a day.", diaString)
	}
	timeString := reSchedule.FindString(scheduleString)
	startTime := reScheduleStart.FindString(timeString)
	endTime := reScheduleFinish.FindString(timeString)
	startTimeHours := reHours.FindString(startTime)
	startTimeMinutes := reMinutes.FindString(startTime)
	endTimeHours := reHours.FindString(endTime)
	endTimeMinutes := reMinutes.FindString(endTime)
	timeStart := NewTime()
	timeEnd := NewTime()
	number1, _ := strconv.Atoi(startTimeHours)
	number2, _ := strconv.Atoi(startTimeMinutes)
	timeStart[0] = number1
	timeStart[1] = number2
	number3, _ := strconv.Atoi(endTimeHours)
	number4, _ := strconv.Atoi(endTimeMinutes)

	timeEnd[0] = number3
	timeStart[1] = number4
	return diaInt, []time{timeStart, timeEnd}, nil
}

func GatherSchedules(classes *[]Class) *[]Cursada {
	//numberOfClasses := len(*classes)
	//verifiedScheduleNumber := 0
	//scheduleListMaster := NewCursadaList()
	currentSchedule := NewCursada()
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
	reSchedule := regexp.MustCompile(`^[\s]{0,99}[A-Za-zéá]{5,9}[\s][0-9:0-9]{5}[\s-]{1,5}[0-9:0-9]{5}`)
	reComisionLabel := regexp.MustCompile(`(?:^[\s]{0,99})[A-Z]{1,8}(?:[\s]{0,99}$)`)
	reEndComision := regexp.MustCompile(`^[\s]{0,99}[0-9]{1,4}[\/\s]{1,3}[0-9]{1,4}[\s]{0,99}$`)

	var (
		currentClass          Class
		allClasses            []Class
		numberString          string
		currentStringSchedule string
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

				currentStringSchedule = reSchedule.FindString(textLine)

				if currentStringSchedule != "" {
					//slicedSchedules := sliceSchedules(&currentStringSchedule)
					diaInt, theTime, err := stringToTime(currentStringSchedule)
					if err != nil {
						return nil, err
					}

					if debug {
						fmt.Printf("[DEBUG] DIA: %s STRUCT TIME: %+v\n\n", Days[diaInt], theTime)
					}
					currentSchedule := NewSchedule()
					currentSchedule.start = theTime[0]
					currentSchedule.start = theTime[0]
					currentSchedule.day = diaInt

					currentComision.schedules = append(currentComision.schedules, currentSchedule)
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









