package heputils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/metrico/promcasa/model"
	"github.com/valyala/bytebufferpool"
)

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorBlue         = "\u001b[34m"
	ColorReset        = "\u001b[0m"
)

// Version holds a loghttp version
type Version int

// Valid Version values
const (
	VersionLegacy = Version(iota)
	VersionV1
)

/* colorize message */
func Colorize(color Color, message string) {
	fmt.Println(string(color), message, string(ColorReset))
}

func MakeJson(lbs []model.Label, sb *bytebufferpool.ByteBuffer) string {
	sb.WriteString(`{`)
	for i, v := range lbs {
		sb.WriteString(`"`)
		sb.WriteString(v.Key)
		sb.WriteString(`":"`)
		sb.WriteString(v.Value)
		sb.WriteString(`"`)
		if i < len(lbs)-1 {
			sb.WriteString(`,`)
		}
	}
	sb.WriteString(`}`)
	return sb.String()
}

func AppendTwoSlices(a []string, b []string) []string {

	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter, _ := range check {
		res = append(res, letter)
	}

	return res
}

func UniqueSlice(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func PartTime(gTime string) int64 {
	parseTime, _ := strconv.Atoi(gTime)
	sec := time.Duration(parseTime) * time.Nanosecond
	return int64(sec.Milliseconds())
}

// GetVersion returns the loghttp version for a given path.
func GetVersion(uri string) Version {
	if strings.Contains(strings.ToLower(uri), "/loki/api/v1") {
		return VersionV1
	}

	return VersionLegacy
}

func SplitDelimiter(r rune) bool {
	return r == ':' || r == '='
}
