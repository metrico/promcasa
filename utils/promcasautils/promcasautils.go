package promcasautils

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/metrico/promcasa/utils/logger"
)

type Color string
type UUID [16]byte

var timeBase = time.Date(1582, time.October, 15, 0, 0, 0, 0, time.UTC).Unix()
var hardwareAddr []byte
var clockSeq uint32

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorBlue         = "\u001b[34m"
	ColorReset        = "\u001b[0m"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var UUIDArray = [...]string{"callid", "session_id", "uuid", "sid", "cid", "correlation_id"}

// import  checkFloatValue
func CheckFloatValue(val interface{}) float64 {
	if val != nil {
		myType := reflect.TypeOf(val)
		switch myType.Kind() {
		case reflect.Int:
			return float64(val.(int))
		case reflect.String:
			tmp, _ := strconv.Atoi(val.(string))
			return float64(tmp)
		case reflect.Int32:
			return float64(val.(int32))
		case reflect.Float64:
			return val.(float64)
		case reflect.Int64:
			return float64(val.(int64))
		default:
			return float64(0)
		}
	}
	return float64(0)
}

// import  checkFloatValue
func CheckStringValue(val interface{}) string {
	if val != nil {
		myType := reflect.TypeOf(val)
		switch myType.Kind() {
		case reflect.Int:
			return fmt.Sprintf("%d", val.(int))
		case reflect.String:
			return val.(string)
		case reflect.Int32:
			return fmt.Sprintf("%d", val.(int32))
		case reflect.Float64:
			return fmt.Sprintf("%f", val.(float64))
		case reflect.Int64:
			return fmt.Sprintf("%d", val.(int64))
		default:
			return ""
		}
	}
	return ""
}

// import  checkFloatValue
func CheckIntValue(val interface{}) int {
	if val != nil {
		myType := reflect.TypeOf(val)
		switch myType.Kind() {
		case reflect.String:
			tmp, _ := strconv.Atoi(val.(string))
			return tmp
		case reflect.Int:
			return val.(int)
		case reflect.Float64:
			return int(val.(float64))
		case reflect.Int64:
			return int(val.(int64))
		default:
			return int(0)
		}
	}
	return int(0)
}

// import  checkFloatValue
func CheckBooleanValue(val interface{}) bool {
	retVal := 0
	if val != nil {
		myType := reflect.TypeOf(val)
		switch myType.Kind() {
		case reflect.String:
			tmp, _ := strconv.ParseBool(val.(string))
			return tmp
		case reflect.Int:
			retVal = val.(int)
		case reflect.Float64:
			retVal = int(val.(float64))
		case reflect.Int64:
			retVal = int(val.(int64))
		case reflect.Bool:
			return val.(bool)
		default:
			return false
		}
	}

	if retVal == 0 {
		return false
	} else {
		return true
	}
}

// import  CheckTypeValue
func CheckTypeValue(val interface{}) reflect.Kind {
	if val != nil {
		return reflect.TypeOf(val).Kind()
	}
	return reflect.Invalid
}

// import YesNo
func YesNo(table string) bool {
	prompt := promptui.Select{
		Label: "Force to populate table [" + table + "]  [Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

// import YesNo
func YesNoPromptText(promptText string) bool {
	prompt := promptui.Select{
		Label: promptText + "  [Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

/* colorize message */
func Colorize(color Color, message string) {
	fmt.Println(string(color), message, string(ColorReset))
}

/* colorize message */
func ColorizeSprintf(color Color, message string) string {
	return fmt.Sprintf(string(color), message, string(ColorReset))
}

func Sanitize(text string) string {

	//if strings.HasPrefix(text, "!=") {
	//	text = strings.TrimPrefix(text, "!=")
	//}

	return strings.NewReplacer(
		`"`, `\"`,
		`&`, "&amp;",
	).Replace(text)
}

func SanitizeTextArray(valArray []string) []string {

	for key, val := range valArray {
		valArray[key] = Sanitize(val)
	}

	return valArray
}

func QueryCheck(query string) error {

	valArray := [...]string{"DELETE", "CREATE", "DROP", "CREATE", "UPDATE"}

	for _, val := range valArray {
		if strings.Contains(strings.ToLower(query), strings.ToLower(val)) {
			return errors.New("bad value inside: " + val)
		}
	}
	return nil
}

func SanitizeIntArray(valArray []string) []int {

	var intArray = []int{}
	for key, val := range valArray {
		intArray[key] = CheckIntValue(val)
	}
	return intArray
}

// import  convertPayloadTypeToString
func ConvertPayloadTypeToString(val float64) (string, string) {

	var Method, Text string

	switch val {
	case 81:
		Method = "TDR"
		Text = "TDR"
		break
	case 100:
		Method = "LOG"
		Text = "LOG"
		break
	case 5:
		Method = "RTCP"
		Text = "RTCP"
		break
	case 34:
		Method = "Report RTP"
		Text = "Report RTP"
		break
	case 35:
		Method = "Report RTP"
		Text = "Report RTP"
		break
	case 200:
		Method = "Loki Data"
		Text = "Loki Data"
		break
	case 54:
		Method = "ISUP"
		Text = "ISUP message"
		break
	default:
		Method = "Generic"
		Text = "generic"
		break
	}

	return Method, Text
}

// import  convertProtoTypeToString
func ConvertProtoTypeToString(val float64) string {

	var protoText string

	switch val {
	case 6:
		protoText = "TCP"
		break
	case 17:
		protoText = "UDP"
		break
	case 132:
		protoText = "SCTP"
		break
	default:
		protoText = "UDP"
		break
	}

	return protoText
}

// import  convertProtoTypeToString
func ConvertProtoStringToInt(val string) int {

	var proto int

	switch val {
	case "TCP":
		proto = 6
		break
	case "UDP":
		proto = 17
		break
	case "SCTP":
		proto = 132
		break
	default:
		proto = 6
		break
	}

	return proto
}

func SanitizeProto(text string) string {

	return strings.NewReplacer(
		`UDP`, `17`,
		`udp`, `17`,
		`TCP`, `6`,
		`tcp`, `6`,
		`SCTP`, `132`,
		`sctp`, `132`,
	).Replace(text)
}

/* isup to HEX */
func IsupToHex(s string) string {
	p1 := strings.Index(s, "/isup")
	if p1 == -1 {
		if p1 = strings.Index(s, "/ISUP"); p1 == -1 {
			return s
		}
	}

	if p2 := strings.Index(s[p1:], "\r\n\r\n"); p2 > -1 {
		p2 = p1 + p2 + 4
		if p3 := strings.Index(s[p2:], "\r\n"); p3 > -1 {
			p3 = p2 + p3
			return injectHex(s, p2, p3)
		} else {
			return injectHex(s, p2, len(s)-1)
		}
	}
	return s
}

func injectHex(s string, start, end int) string {
	return s[:start] + fmt.Sprintf("% X", s[start:end]) + s[end+1:]
}

/* check if the element exists */
func ItemExists(arr []string, elem string) bool {

	for index := range arr {
		if arr[index] == elem {
			return true
		}
	}
	return false
}

/* check if the element exists */
func ElementExists(arr []string, elem string) bool {

	if len(arr) == 0 {
		return true
	}

	if len(arr) == 1 && arr[0] == "" {
		logger.Debug("empty ELements")
		return true
	}

	for index := range arr {
		if arr[index] == elem {
			logger.Debug("Equal 1: ", arr[index])
			logger.Debug("Equal 2: ", elem)
			return true
		}
	}
	return false
}

/* check if the element exists */
func ElementRealExists(arr []string, elem string) bool {

	if len(arr) == 0 {
		return false
	}

	if len(arr) == 1 && arr[0] == "" {
		logger.Debug("empty ELements")
		return false
	}

	for index := range arr {
		if arr[index] == elem {
			logger.Debug("Real exists: ", arr[index], " vs ", elem)
			return true
		}
	}
	return false
}

/* check if the element exists */
func ElementExistsPosition(arr []string, elem string) int {

	if len(arr) == 0 {
		return -1
	}

	if len(arr) == 1 && arr[0] == "" {
		logger.Debug("empty ELements")
		return -1
	}

	for index := range arr {
		if arr[index] == elem {
			logger.Debug("Real exists: ", arr[index], " vs ", elem)
			return index
		}
	}
	return -1
}

/* check if the element exists */
func KeyExists(arr []uint32, elem uint32) bool {

	if len(arr) == 0 {
		return true
	}

	for index := range arr {
		if arr[index] == elem {
			return true
		}
	}
	return false
}

func GenerateToken() string {

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 80)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func MakeTimestamp() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}

// LookupIP looks up host using the local resolver.
// It returns a slice of that host's IPv4 and IPv6 addresses.
func LookupIP(ip string) (string, error) {

	names, err := net.LookupAddr(ip)
	if err != nil {
		return "", err
	}

	if len(names) == 0 {
		return "", nil
	}

	for _, name := range names {
		return name, nil
	}

	return "", nil
}

// UUIDFromTime generates a new time based UUID (version 1) as described in
// RFC 4122. This UUID contains the MAC address of the node that generated
// the UUID, the given timestamp and a sequence number.
func UUIDFromTime(aTime time.Time) UUID {
	var u UUID

	utcTime := aTime.In(time.UTC)
	t := uint64(utcTime.Unix()-timeBase)*10000000 + uint64(utcTime.Nanosecond()/100)
	u[0], u[1], u[2], u[3] = byte(t>>24), byte(t>>16), byte(t>>8), byte(t)
	u[4], u[5] = byte(t>>40), byte(t>>32)
	u[6], u[7] = byte(t>>56)&0x0F, byte(t>>48)

	clock := atomic.AddUint32(&clockSeq, 1)
	u[8] = byte(clock >> 8)
	u[9] = byte(clock)

	copy(u[10:], hardwareAddr)

	u[6] |= 0x10 // set version to 1 (time based uuid)
	u[8] &= 0x3F // clear variant
	u[8] |= 0x80 // set to IETF variant

	return u
}

// String returns the UUID in it's canonical form, a 32 digit hexadecimal
// number in the form of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	var offsets = [...]int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34}
	const hexString = "0123456789abcdef"
	r := make([]byte, 36)
	for i, b := range u {
		r[offsets[i]] = hexString[b>>4]
		r[offsets[i]+1] = hexString[b&0xF]
	}
	r[8] = '-'
	r[13] = '-'
	r[18] = '-'
	r[23] = '-'
	return string(r)

}

// make a hash
func Hash32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// make a hash
func Hash64(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// make a genCodeChallengeS256
func GenCodeChallengeS256(s string) string {

	/*
		hash := hmac.New(sha256.New, []byte(s))
		hex.EncodeToString(hash.Sum(nil))
		return base64.StdEncoding.EncodeToString(hash.Sum(nil))
	*/

	/* we have to remove the = because RFC is not allowed */
	//TrimRight - remove all ==
	s256 := sha256.Sum256([]byte(s))
	return strings.TrimSuffix(base64.URLEncoding.EncodeToString(s256[:]), "=")
}

func HashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
