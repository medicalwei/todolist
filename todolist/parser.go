package todolist

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/now"
)

type Parser struct{}

func (p *Parser) ParseNewTodo(input string) *Todo {
	r, _ := regexp.Compile(`^(add|a)(\\ |)`)
	input = r.ReplaceAllString(input, "")
	if input == "" {
		return nil
	}

	todo := NewTodo()
	todo.Subject = p.Subject(input)
	todo.Projects = p.Projects(input)
	todo.Contexts = p.Contexts(input)
	if p.hasDue(input) {
		todo.Due = p.Due(input, time.Now())
	}
	return todo
}

func (p *Parser) ParseNewTodos(input string) []string {
	return strings.Split(input, ",")
}

func (p *Parser) Subject(input string) string {
	if strings.Contains(input, " due") {
		index := strings.LastIndex(input, " due")
		return input[0:index]
	} else {
		return input
	}
}

func (p *Parser) ExpandProject(input string) string {
	r, _ := regexp.Compile(`\+\w+[^:]`)
	return p.matchWords(input, r)[0]
}

func (p *Parser) Projects(input string) []string {
	r, _ := regexp.Compile(`\+\w+`)
	return p.matchWords(input, r)
}

func (p *Parser) Contexts(input string) []string {
	r, err := regexp.Compile(`\@\w+`)
	if err != nil {
		fmt.Println("regex error", err)
	}
	return p.matchWords(input, r)
}

func (p *Parser) hasDue(input string) bool {
	r1, _ := regexp.Compile(`due \w+$`)
	r2, _ := regexp.Compile(`due \w+ \d+$`)
	return (r1.MatchString(input) || r2.MatchString(input))
}

func (p *Parser) Due(input string, day time.Time) string {
	r, _ := regexp.Compile(`due .*$`)

	res := r.FindString(input)
	res = res[4:len(res)]
	switch res {
	case "none":
		return ""
	case "today", "tod":
		return now.BeginningOfDay().Format("2006-01-02")
	case "tomorrow", "tom":
		return now.BeginningOfDay().AddDate(0, 0, 1).Format("2006-01-02")
	case "monday", "mon":
		return p.monday(day)
	case "tuesday", "tue":
		return p.tuesday(day)
	case "wednesday", "wed":
		return p.wednesday(day)
	case "thursday", "thu":
		return p.thursday(day)
	case "friday", "fri":
		return p.friday(day)
	case "saturday", "sat":
		return p.saturday(day)
	case "sunday", "sun":
		return p.sunday(day)
	case "next week":
		n := now.BeginningOfDay()
		return now.New(n).Monday().AddDate(0, 0, 7).Format("2006-01-02")
	}
	return p.parseArbitraryDate(res, time.Now())
}

func (p *Parser) parseArbitraryDate(_date string, pivot time.Time) string {
	d1 := p.parseArbitraryDateWithYear(_date, pivot.Year())

	var diff1 time.Duration
	if d1.After(time.Now()) {
		diff1 = d1.Sub(pivot)
	} else {
		diff1 = pivot.Sub(d1)
	}
	d2 := p.parseArbitraryDateWithYear(_date, pivot.Year()+1)
	if d2.Sub(pivot) > diff1 {
		return d1.Format("2006-01-02")
	} else {
		return d2.Format("2006-01-02")
	}
}

func (p *Parser) parseArbitraryDateWithYear(_date string, year int) time.Time {
	res := strings.Join([]string{_date, strconv.Itoa(year)}, " ")
	if date, err := time.Parse("Jan 2 2006", res); err == nil {
		return date
	}

	if date, err := time.Parse("2 Jan 2006", res); err == nil {
		return date
	}
	panic(fmt.Errorf("Could not parse the date you gave me: '%s'", _date))
}

func (p *Parser) monday(day time.Time) string {
	mon := now.New(day).Monday()
	return p.thisOrNextWeek(mon, day)
}

func (p *Parser) tuesday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 1)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) wednesday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 2)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) thursday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 3)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) friday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 4)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) saturday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 5)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) sunday(day time.Time) string {
	tue := now.New(day).Monday().AddDate(0, 0, 6)
	return p.thisOrNextWeek(tue, day)
}

func (p *Parser) thisOrNextWeek(day time.Time, pivotDay time.Time) string {
	if day.Before(pivotDay) {
		return day.AddDate(0, 0, 7).Format("2006-01-02")
	} else {
		return day.Format("2006-01-02")
	}
}

func (p *Parser) matchWords(input string, r *regexp.Regexp) []string {
	results := r.FindAllString(input, -1)
	ret := []string{}

	for _, val := range results {
		ret = append(ret, val[1:len(val)])
	}
	return ret
}
