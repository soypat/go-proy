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

var Days = [...]string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"}
var DaysNoAccent = [...]string{"Lunes", "Martes", "Miercoles", "Jueves", "Viernes", "Sabado", "Domingo"}
var DaysBadParse = [...]string{"Lúnes", "Mártes", "Mi�rcoles", "Jueves", "Viernes", "Sébado", "Domingo"}
var DaysEnglish = [...]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

func main() {
	criteria := NewScheduleCriteria()
	criteria.maxSuperposition = 5.2 // en horas
	criteria.maxTotalSuperposition = 90
	criteria.minFreeDays = 0
	criteria.maxNumberOfSuperpositions = 7
	//criteria.freeDays[4] = true
	Classes, err := GatherClasses("data_superpos.dat")
	if err != nil {
		panic("Big baddy")
	}

	//fmt.Printf("%+v", *Classes)
	ScheduleList := GetSchedules(Classes, &criteria)
	if ScheduleList != nil {
		for _, v := range *ScheduleList {
			fmt.Printf("\n\n%+v", v)
		}
	} else {
		fmt.Printf("No schedules found.")
	}

}



type comision struct {
	label     string
	schedules []schedule
	teachers  []string
}
type schedule struct {
	day   int // from 0 to 6
	start time
	end   time
}

func (mySchedule schedule) Duration() float32 {
	return float32(mySchedule.end.hour - mySchedule.start.hour + (mySchedule.end.minute-mySchedule.start.minute)/60)

}

func NewSchedule() schedule {
	return schedule{}
}

type time struct {
	hour   int
	minute int
}

func NewTime() time {
	return time{}
}
func NewComision() comision {
	return comision{}
}

type scheduleCriteria struct {
	maxSuperposition          float32
	maxTotalSuperposition     float32
	maxNumberOfSuperpositions int
	freeDays                  [len(Days)]bool
	minFreeDays               int
}

// Class as in school class
type Class struct {
	num1       int
	num2       int
	name       string
	comisiones []comision
}

type Cursada []comision // Cursada is a the group of courses a student attends during the year/semester

func NewScheduleCriteria() scheduleCriteria {
	return scheduleCriteria{}
}

func NewCursada() Cursada {
	return Cursada{}
}

func NewCursadaList() []Cursada {
	return []Cursada{}
}

func NewClass() Class {
	return Class{}
}



func searcher(classes *[]Class, currentCursada *Cursada, classNumber int, criteria *scheduleCriteria) *[]Cursada {
	nextClass := (*classes)[classNumber]
	cursadaListMaster := NewCursadaList()
	for _, v := range nextClass.comisiones {
		cursadaInstance := append(*currentCursada, v)

		if classNumber == len(*classes)-1 { //llegue a la ultima clase
			isValid := verifyCursada(&cursadaInstance, criteria)
			if isValid { //El schedule es bueno, lo devuelvo como lista no nula
				cursadaListMaster = append(cursadaListMaster, cursadaInstance)
				continue
			} else {
				continue
			}

		} else { // Si no es la ultima clase, sigo por aca
			cursadaList := searcher(classes, &cursadaInstance, classNumber+1, criteria) // Awesome recursion baby
			if len(*cursadaList) == 0 {
				continue
			}
			cursadaListMaster = append(cursadaListMaster, *cursadaList...)
		}
	}
	return &cursadaListMaster
}

func findCollision(schedule1 *schedule, schedule2 *schedule) float32 {
	if (schedule1.start.hour >= schedule2.start.hour && schedule1.start.hour < schedule2.end.hour) && schedule1.Duration() >= 0.5 {
		if schedule1.end.hour <= schedule2.end.hour {
			return schedule1.Duration()
		} else {
			return schedule1.Duration() + float32(-schedule1.end.hour+schedule2.end.hour)
		}
	}
	if (schedule2.start.hour >= schedule1.start.hour && schedule2.start.hour < schedule1.end.hour) && schedule2.Duration() >= 0.5 {
		if schedule2.end.hour <= schedule1.end.hour {
			return schedule2.Duration()
		} else {
			return schedule2.Duration() + float32(-schedule2.end.hour+schedule1.end.hour)
		}
	}
	return 0.0
}

func verifyCursada(currentCursada *Cursada, criteria *scheduleCriteria) bool {
	// TODO hard part coming ahead. Actual verification
	numberOfMaterias := len(*currentCursada)
	superpositionCounter := 0
	totalSuperpositions := float32(0)
	var busyDays = []bool{false, false, false, false, false, false, false}
	for i := 0; i < numberOfMaterias-1; i++ {
		firstComision := (*currentCursada)[i]
		for j := numberOfMaterias - 1; j > i; j-- {
			secondComision := (*currentCursada)[j]
			firstCursada := firstComision.schedules
			secondCursada := secondComision.schedules
			for _, schedule1 := range firstCursada {
				for _, schedule2 := range secondCursada {
					busyDays[schedule1.day] = true
					busyDays[schedule2.day] = true
					if criteria.freeDays[schedule1.day] || criteria.freeDays[schedule2.day] { // Criterio absoluto. Si quiero un free day, entonces se va hacer un free day
						return false
					}
					if schedule1.day != schedule2.day { // si no coinciden los dias, verifica ese horario, continuo buscando colisioneschedule1s
						continue
					} else { // en el caso que sean el mismo día:
						superpositions := findCollision(&schedule1, &schedule2)
						if superpositions == 0.0 { // NO se encuentran superposiciones
							continue
						} else {
							superpositionCounter++
							totalSuperpositions += superpositions
							if totalSuperpositions > criteria.maxTotalSuperposition || superpositionCounter > criteria.maxNumberOfSuperpositions || superpositions > criteria.maxSuperposition {
								return false
							}
						}
					}
				}
			}

		}
	}
	// Verificación final de dias ocupados
	for i, b := range busyDays {
		if criteria.freeDays[i] && b {
			return false
		}
	}
	return true
}

func stringToTime(scheduleString string) (int, []time, error) {
	reWeek := regexp.MustCompile(`(?i)Lunes|Martes|Miércoles|Jueves|Viernes|Sábado|Domingo|Miercoles|Sabado|Sébado`)
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
		} else if strings.Contains(strings.Title(diaString), DaysNoAccent[i]) || strings.Contains(strings.Title(diaString), DaysBadParse[i]) || strings.Contains(strings.Title(diaString), DaysEnglish[i]) {
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
	timeStart.hour = number1
	timeStart.minute = number2
	number3, _ := strconv.Atoi(endTimeHours)
	number4, _ := strconv.Atoi(endTimeMinutes)

	timeEnd.hour = number3
	timeStart.minute = number4
	return diaInt, []time{timeStart, timeEnd}, nil
}

func GetSchedules(classes *[]Class, criteria *scheduleCriteria) *[]Cursada {
	//numberOfClasses := len(*classes)
	//verifiedScheduleNumber := 0
	//scheduleListMaster := NewCursadaList()
	currentSchedule := NewCursada()
	scheduleListMaster := searcher(classes, &currentSchedule, 0, criteria)

	// Search function: cada instancia de recursividad busca verificar un schedule (lo va fabricando a medida que avanza) y devuelve una lista de schedules verificados y los va juntando en cada instancia
	if len(*scheduleListMaster) == 0{
		return nil
	}
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
	reSchedule := regexp.MustCompile(`^[\s]{0,99}[A-Za-zéá�]{5,9}[\s][0-9:0-9]{5}[\s-]{1,5}[0-9:0-9]{5}`)
	reComisionLabel := regexp.MustCompile(`(?:^[\s]{0,99})[A-Z]{1,8}(?:[\s]{0,99}$)`)
	reEndComision := regexp.MustCompile(`^[\s]{0,99}[0-9]{1,4}[\/\s]{1,3}[0-9]{1,4}[\s]{0,99}$`)
	reEAccent := regexp.MustCompile(`[�]{1}`)

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
		//Sanitize unicode disgusting badness
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
					currentStringSchedule = reEAccent.ReplaceAllString(currentStringSchedule, "é")
					diaInt, theTime, err := stringToTime(currentStringSchedule)
					if err != nil {
						return nil, err
					}

					if debug {
						fmt.Printf("[DEBUG] DIA: %s STRUCT TIME: %+v\n\n", Days[diaInt], theTime)
					}

					currentSchedule := NewSchedule()
					currentSchedule.start = theTime[0]
					currentSchedule.end = theTime[1]
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
