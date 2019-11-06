package helpers

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/denisbrodbeck/machineid"
	"github.com/globalsign/mgo/bson"
	uuid "github.com/satori/go.uuid"
)

// Payload -- used to pass variables into fn.Execute
type Payload map[string]interface{}

// FailOnError will log error and reply if required when err object is nil
func FailOnError(err error, msg string) {
	if err != nil {
		fmt.Println(err, msg)
		os.Exit(3)
	}
}

// Log will log msg to terminal
func Log(msg string, mode int) {
	if AppMode >= mode {
		fmt.Println("Log --> ", msg)
	}
}

// GetUUID returns UUID
func GetUUID() (string, error) {
	u2, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u2.String(), nil
}

// GetUUIDLite returns UUID without error
func GetUUIDLite() string {
	u2, _ := uuid.NewV4()
	return u2.String()
}

// InArray checks if given string is inside array
func InArray(val string, _records interface{}) (exists bool, index int) {
	records := _records.([]string)
	for i, b := range records {
		if b == val {
			return true, i
		}
	}
	return false, -1

	// switch reflect.TypeOf(array).Kind() {
	// case reflect.Slice:
	// 	s := reflect.ValueOf(array)

	// 	for i := 0; i < s.Len(); i++ {
	// 		if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
	// 			index = i
	// 			exists = true
	// 			return
	// 		}
	// 	}
	// }

	// return
}

// ObjInArray checks if given bson objectid is inside array
func ObjInArray(val bson.ObjectId, records []bson.ObjectId) (exists bool, index int) {
	for i, b := range records {
		if b == val {
			return true, i
		}
	}
	return false, -1

	// switch reflect.TypeOf(array).Kind() {
	// case reflect.Slice:
	// 	s := reflect.ValueOf(array)

	// 	for i := 0; i < s.Len(); i++ {
	// 		if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
	// 			index = i
	// 			exists = true
	// 			return
	// 		}
	// 	}
	// }

	// return
}

// DownloadFile -- downloads file from url to given path
func DownloadFile(url string, localPath string) error {

	out, err1 := os.Create(localPath)
	if err1 != nil {
		return err1
	}

	// download from S3
	data, err2 := http.Get(url)
	if err2 != nil {
		return err2
	}

	if data.StatusCode >= 200 && data.StatusCode <= 299 {
		// add downloaded content to file
		_, err3 := io.Copy(out, data.Body)
		if err3 != nil {
			return err3
		}

		defer data.Body.Close()
		defer out.Close()

		return nil
	}

	return errors.New("got status code : " + GetString(data.StatusCode))

}

// CopyFile -- copy file from one path to another
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// GetMD5Hash -- get MD5 of given string
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GetString -- get string of given variable
func GetString(val interface{}) string {
	return fmt.Sprint(val)
}

// IsEmpty -- returns boolean of whether given struct is empty or not
func IsEmpty(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

// KillProcess -- will kill any linux process -- used in killing yolo training process
func KillProcess(processID int) {
	syscall.Kill(processID, syscall.SIGKILL)
}

// GetGPUStats -- returns GPU stats
func GetGPUStats() (map[string]float64, error) {
	// check existing processes + GPU usage here
	cmd := exec.Command("/usr/bin/nvidia-smi")
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	// get GPU usage
	r, _ := regexp.Compile("(\\d*)MiB")
	matches := r.FindAllStringSubmatch(string(data), 50)
	usedGPU, _ := strconv.ParseFloat(matches[0][1], 32)
	totalGPU, _ := strconv.ParseFloat(matches[1][1], 32)

	stats := make(map[string]float64)
	stats["used"] = usedGPU
	stats["total"] = totalGPU

	return stats, nil
}

// GetMachineUID -- returns machine unique ID
func GetMachineUID() bson.M {
	secret := "testing123"
	id, _ := machineid.ID()
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(secret))

	machineID := fmt.Sprintf("%x", mac.Sum(nil))
	// ip, _ := externalip.DefaultConsensus(nil, nil).ExternalIP()

	return bson.M{
		"machine_id": machineID,
		// "ip":         ip,
	}
}

// Convert -- converts one type to another
func Convert(input interface{}, output *interface{}) error {
	var err error
	switch input.(type) {
	case []byte:
		byteArray := input.([]byte)
		err = json.Unmarshal(byteArray, output)
	default:
		bytesArray, err2 := json.Marshal(input)
		err2 = json.Unmarshal(bytesArray, output)
		err = err2
	}
	return err
}

//ConvertResponse -- converts payload into response type
func ConvertResponse(input interface{}, output *Response) error {
	var err error
	switch input.(type) {
	case []byte:
		byteArray := input.([]byte)
		err = json.Unmarshal(byteArray, output)
	default:
		bytesArray, err2 := json.Marshal(input)
		err2 = json.Unmarshal(bytesArray, output)
		err = err2
	}
	return err
}

// ConvertMap -- converts any type to map of interface
func ConvertMap(input interface{}, output *Payload) error {
	var err error
	switch input.(type) {
	case []byte:
		byteArray := input.([]byte)
		err = json.Unmarshal(byteArray, output)
	default:
		bytesArray, err2 := json.Marshal(input)
		err2 = json.Unmarshal(bytesArray, output)
		err = err2
	}
	return err
}

// ClientKey -- this contains private key of the client application
var ClientKey *rsa.PrivateKey

// ServerKey -- this contains public key of the server application
var ServerKey *rsa.PrivateKey

// VerifyLicense -- verifies whether the device is authorized to access teric.ai and returns a boolean
func VerifyLicense() bool {

	return true

	// load client public + private key ---

	// clientPrivateKey, err := os.Open("client.pem")
	// if err != nil {
	// 	fmt.Println("cannot open client certificate file")
	// 	os.Exit(1)
	// }
	// clientPEMInfo, _ := clientPrivateKey.Stat()
	// clientPEMbytes := make([]byte, clientPEMInfo.Size())
	// clientBuffer := bufio.NewReader(clientPrivateKey)
	// _, err = clientBuffer.Read(clientPEMbytes)

	clientPEMbytes := []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 82, 83, 65, 32, 80, 82, 73, 86, 65, 84, 69, 32, 75, 69, 89, 45, 45, 45, 45, 45, 10, 77, 73, 73, 69, 112, 65, 73, 66, 65, 65, 75, 67, 65, 81, 69, 65, 116, 116, 106, 88, 89, 98, 77, 97, 83, 116, 78, 87, 84, 83, 115, 51, 121, 98, 77, 122, 109, 68, 114, 48, 118, 76, 103, 67, 84, 66, 75, 114, 54, 69, 88, 99, 50, 88, 89, 88, 48, 84, 77, 83, 102, 51, 80, 83, 10, 48, 86, 114, 90, 106, 56, 99, 79, 109, 66, 97, 49, 101, 102, 82, 79, 116, 106, 115, 121, 102, 118, 113, 99, 70, 117, 77, 66, 74, 78, 119, 67, 105, 79, 52, 80, 120, 52, 119, 76, 109, 103, 69, 54, 121, 119, 106, 103, 122, 73, 88, 112, 54, 47, 83, 86, 88, 122, 103, 57, 108, 119, 118, 77, 10, 112, 112, 43, 100, 65, 47, 55, 67, 100, 48, 117, 107, 68, 80, 68, 67, 47, 72, 121, 65, 116, 97, 109, 80, 87, 118, 47, 73, 105, 114, 78, 97, 119, 84, 73, 71, 52, 106, 110, 51, 43, 105, 57, 120, 81, 55, 101, 81, 79, 89, 107, 101, 77, 77, 79, 112, 87, 82, 49, 43, 53, 105, 102, 83, 10, 56, 117, 56, 49, 69, 50, 118, 110, 57, 81, 119, 119, 120, 66, 52, 113, 79, 82, 90, 110, 122, 69, 111, 100, 100, 118, 100, 80, 71, 66, 51, 106, 71, 68, 51, 113, 98, 119, 114, 115, 72, 74, 83, 66, 84, 79, 121, 78, 73, 80, 106, 109, 105, 87, 57, 90, 72, 86, 49, 101, 51, 77, 105, 83, 10, 106, 80, 78, 67, 49, 70, 88, 102, 80, 52, 51, 55, 75, 68, 57, 113, 48, 122, 115, 116, 119, 100, 55, 117, 70, 109, 78, 69, 103, 88, 48, 86, 79, 106, 113, 109, 67, 73, 98, 113, 79, 47, 49, 43, 119, 108, 113, 104, 72, 88, 121, 75, 48, 98, 104, 115, 69, 55, 120, 113, 48, 77, 114, 122, 10, 65, 67, 89, 74, 121, 68, 115, 74, 86, 115, 55, 47, 103, 100, 120, 52, 107, 105, 53, 102, 102, 103, 98, 112, 73, 87, 109, 118, 56, 71, 78, 57, 117, 111, 68, 68, 79, 119, 73, 68, 65, 81, 65, 66, 65, 111, 73, 66, 65, 70, 73, 109, 113, 78, 90, 118, 104, 116, 43, 90, 104, 107, 118, 84, 10, 111, 66, 81, 83, 88, 74, 115, 72, 50, 103, 43, 48, 83, 79, 118, 117, 56, 54, 101, 47, 81, 57, 79, 56, 101, 69, 84, 52, 119, 108, 88, 98, 76, 120, 118, 54, 121, 111, 99, 76, 115, 50, 88, 110, 120, 103, 43, 79, 69, 90, 78, 85, 107, 52, 74, 122, 106, 73, 47, 72, 51, 67, 113, 52, 10, 89, 114, 99, 115, 53, 112, 65, 77, 80, 117, 89, 112, 113, 85, 87, 120, 114, 110, 97, 86, 115, 66, 122, 103, 88, 103, 66, 84, 72, 51, 68, 117, 122, 122, 115, 74, 117, 90, 48, 105, 54, 68, 74, 55, 72, 76, 68, 110, 116, 50, 79, 68, 101, 76, 121, 108, 43, 119, 43, 121, 110, 109, 75, 97, 10, 53, 75, 113, 113, 108, 71, 99, 117, 68, 108, 107, 115, 50, 97, 72, 73, 74, 112, 101, 68, 73, 76, 112, 101, 72, 111, 99, 51, 115, 57, 56, 56, 71, 79, 90, 98, 115, 80, 52, 100, 112, 106, 116, 111, 112, 57, 71, 86, 121, 49, 54, 97, 69, 56, 86, 103, 122, 49, 87, 56, 75, 48, 71, 54, 10, 112, 109, 108, 53, 66, 51, 52, 111, 81, 111, 47, 104, 121, 89, 77, 70, 70, 69, 88, 120, 51, 119, 106, 77, 85, 85, 100, 81, 112, 83, 70, 52, 89, 75, 88, 82, 56, 116, 47, 85, 57, 118, 122, 103, 112, 53, 88, 50, 114, 56, 103, 90, 74, 112, 122, 57, 116, 118, 52, 104, 119, 97, 105, 54, 10, 104, 121, 117, 119, 43, 105, 120, 119, 68, 84, 72, 85, 116, 66, 54, 49, 75, 89, 77, 67, 90, 79, 88, 72, 84, 52, 52, 81, 114, 68, 85, 120, 89, 108, 120, 118, 66, 56, 118, 68, 104, 97, 57, 87, 105, 65, 54, 122, 77, 80, 71, 118, 75, 76, 66, 72, 83, 48, 83, 53, 50, 99, 43, 81, 10, 117, 70, 80, 114, 86, 104, 107, 67, 103, 89, 69, 65, 48, 107, 102, 87, 105, 53, 49, 116, 99, 69, 102, 77, 77, 67, 69, 70, 81, 86, 122, 68, 87, 72, 52, 53, 84, 101, 120, 89, 83, 43, 50, 50, 114, 48, 80, 109, 82, 68, 65, 51, 49, 109, 106, 108, 118, 76, 51, 68, 105, 74, 113, 84, 10, 75, 109, 47, 48, 104, 65, 54, 53, 52, 86, 52, 51, 70, 54, 65, 47, 80, 108, 102, 119, 118, 107, 85, 78, 80, 76, 71, 69, 111, 111, 79, 101, 111, 89, 89, 115, 47, 47, 86, 97, 79, 66, 78, 52, 118, 119, 67, 43, 90, 56, 86, 101, 112, 56, 52, 104, 82, 87, 89, 75, 43, 72, 117, 109, 10, 82, 117, 55, 76, 117, 122, 67, 73, 50, 105, 65, 100, 87, 107, 53, 122, 116, 84, 99, 65, 68, 77, 107, 104, 116, 119, 113, 76, 117, 82, 66, 113, 104, 84, 101, 68, 102, 116, 90, 117, 50, 85, 86, 104, 116, 87, 114, 116, 75, 57, 75, 69, 106, 121, 85, 67, 103, 89, 69, 65, 51, 112, 111, 79, 10, 118, 43, 106, 90, 120, 75, 80, 57, 103, 103, 67, 73, 108, 117, 85, 97, 71, 43, 69, 78, 112, 87, 100, 86, 102, 98, 80, 65, 47, 110, 67, 98, 52, 88, 102, 72, 77, 89, 119, 56, 50, 67, 89, 119, 74, 117, 71, 43, 111, 113, 86, 72, 87, 121, 108, 108, 50, 115, 74, 88, 83, 74, 68, 81, 10, 79, 71, 65, 84, 87, 106, 90, 88, 47, 89, 55, 105, 67, 83, 53, 47, 101, 80, 79, 73, 78, 69, 81, 104, 83, 78, 99, 72, 115, 51, 75, 74, 99, 43, 80, 122, 120, 100, 102, 65, 80, 56, 82, 88, 115, 76, 119, 109, 48, 53, 113, 80, 90, 68, 89, 113, 50, 117, 65, 73, 67, 82, 121, 80, 10, 120, 108, 97, 74, 48, 107, 48, 90, 114, 73, 50, 111, 107, 52, 89, 81, 120, 50, 88, 67, 112, 121, 79, 79, 81, 107, 83, 108, 113, 122, 120, 76, 122, 68, 79, 99, 75, 116, 56, 67, 103, 89, 69, 65, 110, 56, 115, 114, 104, 69, 107, 76, 103, 119, 108, 115, 90, 120, 54, 81, 113, 99, 122, 101, 10, 80, 88, 56, 100, 43, 78, 77, 106, 102, 102, 43, 85, 108, 98, 100, 90, 89, 110, 80, 112, 50, 113, 115, 51, 43, 97, 101, 83, 48, 86, 110, 49, 102, 52, 103, 52, 72, 97, 111, 55, 73, 115, 71, 47, 120, 57, 112, 107, 100, 80, 72, 75, 53, 105, 118, 47, 70, 83, 73, 112, 69, 110, 53, 71, 10, 113, 54, 81, 85, 121, 105, 85, 101, 102, 65, 74, 47, 47, 86, 87, 74, 87, 55, 52, 109, 89, 103, 112, 73, 83, 106, 53, 122, 69, 56, 83, 83, 53, 78, 66, 79, 84, 86, 57, 105, 102, 54, 57, 114, 51, 116, 90, 68, 73, 51, 65, 54, 80, 51, 48, 81, 101, 57, 73, 116, 118, 50, 74, 48, 10, 76, 43, 117, 120, 112, 48, 56, 52, 83, 83, 57, 113, 81, 114, 121, 81, 111, 110, 54, 70, 99, 87, 107, 67, 103, 89, 69, 65, 48, 71, 80, 51, 70, 79, 47, 47, 70, 108, 106, 84, 108, 101, 87, 55, 43, 85, 43, 72, 84, 114, 119, 48, 107, 122, 107, 87, 122, 114, 80, 43, 73, 47, 84, 49, 10, 54, 88, 68, 66, 109, 81, 65, 74, 89, 101, 122, 50, 80, 83, 65, 117, 52, 73, 76, 77, 78, 50, 113, 100, 65, 78, 118, 89, 55, 73, 85, 116, 101, 79, 108, 119, 108, 73, 54, 49, 100, 120, 108, 82, 81, 72, 107, 52, 79, 116, 110, 54, 69, 55, 119, 73, 85, 80, 71, 70, 77, 120, 103, 120, 10, 49, 55, 49, 54, 86, 67, 101, 122, 119, 98, 54, 107, 118, 84, 54, 88, 78, 112, 102, 71, 84, 51, 70, 113, 85, 122, 100, 83, 76, 110, 49, 47, 108, 53, 85, 105, 78, 121, 43, 89, 114, 110, 74, 55, 99, 52, 103, 90, 111, 121, 72, 47, 120, 89, 114, 67, 118, 103, 85, 88, 57, 121, 78, 101, 10, 98, 107, 43, 79, 106, 111, 56, 67, 103, 89, 66, 52, 77, 106, 81, 117, 52, 112, 56, 112, 72, 48, 70, 71, 115, 80, 107, 104, 56, 73, 89, 119, 121, 111, 104, 82, 81, 85, 76, 71, 85, 101, 101, 74, 56, 54, 73, 43, 120, 110, 86, 84, 99, 89, 69, 67, 101, 83, 110, 82, 57, 67, 71, 106, 10, 114, 77, 104, 82, 118, 122, 98, 87, 88, 89, 43, 80, 80, 69, 86, 57, 113, 97, 107, 78, 88, 65, 74, 51, 97, 75, 65, 122, 111, 114, 75, 43, 97, 90, 78, 106, 50, 50, 71, 110, 79, 114, 87, 79, 54, 70, 109, 55, 53, 106, 119, 118, 98, 82, 104, 85, 73, 111, 108, 110, 120, 107, 73, 90, 10, 105, 43, 49, 82, 43, 109, 83, 101, 56, 112, 106, 49, 113, 97, 115, 90, 97, 109, 120, 117, 106, 90, 113, 49, 66, 67, 47, 83, 119, 102, 73, 116, 69, 97, 73, 74, 66, 70, 104, 65, 118, 122, 108, 116, 98, 78, 87, 70, 69, 115, 104, 81, 118, 65, 61, 61, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 82, 83, 65, 32, 80, 82, 73, 86, 65, 84, 69, 32, 75, 69, 89, 45, 45, 45, 45, 45, 10}
	clientData, _ := pem.Decode([]byte(clientPEMbytes))
	// clientPrivateKey.Close()

	// fmt.Println(" Log -- clientData", clientPEMbytes, string(clientPEMbytes))

	ClientKey, err := x509.ParsePKCS1PrivateKey(clientData.Bytes)
	if err != nil {
		fmt.Println("cannot parse client certificate file")
		os.Exit(1)
	}

	// load server public key
	// serverPrivateKey, err := os.Open("server.pem")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// serverPEMInfo, _ := serverPrivateKey.Stat()
	// serverSize := serverPEMInfo.Size()
	// serverPEMbytes := make([]byte, serverSize)
	// serverBuffer := bufio.NewReader(serverPrivateKey)
	// _, err = serverBuffer.Read(serverPEMbytes)

	// serverPrivateKey.Close()

	serverPEMbytes := []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 82, 83, 65, 32, 80, 82, 73, 86, 65, 84, 69, 32, 75, 69, 89, 45, 45, 45, 45, 45, 10, 77, 73, 73, 69, 112, 81, 73, 66, 65, 65, 75, 67, 65, 81, 69, 65, 57, 120, 113, 76, 103, 97, 75, 71, 119, 43, 109, 55, 75, 100, 117, 108, 77, 102, 87, 110, 117, 121, 65, 108, 98, 55, 75, 68, 105, 113, 79, 90, 73, 80, 66, 50, 68, 103, 90, 49, 87, 54, 84, 115, 51, 90, 100, 49, 10, 70, 51, 51, 121, 52, 57, 119, 97, 54, 55, 73, 75, 115, 69, 67, 118, 107, 51, 122, 89, 108, 47, 69, 105, 87, 48, 73, 54, 116, 53, 102, 53, 90, 70, 117, 48, 111, 85, 86, 115, 43, 119, 68, 110, 67, 101, 79, 49, 72, 65, 113, 103, 55, 78, 51, 55, 81, 48, 103, 109, 72, 70, 122, 118, 10, 69, 120, 107, 117, 121, 107, 66, 116, 70, 49, 49, 48, 120, 88, 106, 86, 110, 110, 115, 112, 121, 80, 108, 116, 49, 90, 47, 114, 117, 116, 101, 114, 87, 87, 55, 111, 75, 89, 84, 72, 97, 87, 81, 68, 109, 55, 85, 120, 122, 108, 116, 83, 49, 87, 97, 75, 86, 47, 114, 110, 57, 114, 102, 110, 10, 53, 112, 66, 85, 73, 73, 118, 75, 120, 47, 104, 75, 67, 68, 84, 102, 117, 120, 85, 112, 70, 110, 108, 122, 101, 75, 76, 121, 85, 112, 109, 51, 75, 85, 114, 87, 87, 107, 85, 90, 70, 102, 53, 48, 97, 119, 117, 65, 108, 70, 97, 50, 76, 86, 68, 112, 77, 110, 73, 100, 78, 79, 48, 112, 10, 57, 73, 52, 86, 99, 68, 121, 105, 98, 83, 122, 97, 89, 107, 116, 75, 66, 70, 88, 85, 77, 81, 106, 98, 119, 115, 67, 52, 115, 116, 69, 74, 107, 108, 106, 119, 108, 89, 78, 73, 55, 107, 65, 85, 86, 106, 70, 111, 77, 48, 54, 121, 89, 118, 72, 76, 112, 76, 98, 107, 121, 89, 47, 72, 10, 113, 104, 72, 117, 86, 50, 53, 43, 65, 67, 106, 85, 120, 72, 72, 111, 43, 57, 75, 103, 108, 116, 112, 80, 101, 54, 111, 84, 71, 77, 99, 122, 68, 54, 73, 56, 67, 81, 73, 68, 65, 81, 65, 66, 65, 111, 73, 66, 65, 81, 67, 85, 117, 112, 57, 53, 87, 87, 43, 118, 47, 56, 67, 116, 10, 103, 119, 121, 57, 77, 49, 84, 80, 112, 112, 117, 104, 122, 86, 113, 114, 87, 97, 106, 84, 85, 75, 104, 100, 55, 76, 107, 54, 102, 100, 119, 114, 121, 47, 117, 111, 78, 105, 67, 53, 48, 85, 78, 75, 49, 104, 68, 107, 52, 83, 112, 77, 112, 88, 112, 103, 105, 98, 122, 97, 72, 78, 84, 109, 10, 113, 69, 120, 116, 103, 86, 48, 74, 76, 74, 90, 90, 120, 99, 78, 75, 67, 111, 112, 53, 53, 70, 80, 84, 47, 104, 65, 56, 65, 80, 77, 102, 89, 122, 104, 113, 48, 70, 57, 47, 85, 75, 80, 89, 121, 109, 70, 56, 99, 105, 120, 120, 104, 85, 81, 122, 79, 82, 53, 73, 49, 69, 97, 52, 10, 89, 82, 77, 55, 99, 121, 117, 57, 119, 98, 79, 99, 49, 90, 118, 117, 88, 110, 77, 112, 57, 52, 71, 49, 47, 70, 65, 77, 102, 48, 80, 74, 57, 120, 113, 89, 114, 72, 112, 122, 99, 70, 80, 118, 118, 52, 106, 102, 86, 99, 110, 52, 111, 103, 113, 76, 73, 84, 68, 105, 73, 74, 89, 53, 10, 65, 121, 122, 49, 81, 87, 70, 47, 115, 52, 77, 66, 101, 50, 67, 54, 121, 54, 52, 66, 115, 67, 85, 116, 116, 98, 65, 51, 104, 71, 87, 104, 72, 83, 114, 112, 117, 47, 122, 116, 101, 79, 89, 54, 43, 87, 108, 112, 84, 65, 80, 102, 48, 98, 110, 57, 83, 57, 78, 103, 80, 52, 116, 111, 10, 99, 67, 49, 106, 120, 76, 75, 90, 66, 43, 83, 104, 103, 98, 51, 70, 102, 74, 112, 49, 114, 70, 105, 48, 120, 98, 109, 57, 51, 105, 87, 69, 83, 117, 105, 72, 78, 55, 75, 113, 83, 113, 81, 102, 55, 105, 52, 110, 90, 78, 79, 89, 103, 115, 88, 51, 110, 99, 90, 66, 89, 90, 57, 56, 10, 116, 75, 78, 109, 72, 81, 75, 53, 65, 111, 71, 66, 65, 80, 47, 56, 66, 98, 55, 90, 85, 103, 90, 43, 122, 110, 113, 82, 80, 121, 66, 82, 111, 116, 76, 122, 71, 120, 102, 76, 101, 106, 67, 98, 50, 84, 121, 110, 82, 86, 77, 80, 83, 86, 106, 52, 76, 90, 69, 120, 75, 109, 66, 88, 10, 103, 113, 50, 97, 90, 68, 49, 88, 69, 97, 71, 98, 101, 52, 71, 89, 54, 85, 54, 78, 66, 82, 49, 52, 105, 66, 48, 118, 109, 103, 111, 84, 98, 70, 103, 89, 100, 98, 102, 50, 89, 68, 69, 69, 122, 68, 114, 98, 90, 56, 101, 69, 73, 100, 115, 110, 103, 56, 53, 102, 76, 83, 53, 81, 10, 107, 82, 97, 86, 101, 55, 50, 106, 68, 97, 111, 88, 87, 65, 89, 97, 111, 70, 56, 73, 72, 69, 87, 65, 97, 98, 68, 82, 115, 89, 74, 56, 73, 99, 97, 78, 49, 52, 112, 119, 69, 101, 49, 54, 47, 120, 101, 51, 89, 47, 103, 73, 57, 53, 72, 110, 65, 111, 71, 66, 65, 80, 99, 101, 10, 89, 109, 57, 97, 65, 119, 69, 100, 50, 78, 87, 104, 51, 43, 51, 121, 50, 51, 82, 56, 110, 98, 57, 118, 86, 53, 110, 50, 68, 77, 121, 47, 43, 75, 48, 89, 49, 100, 57, 122, 104, 55, 107, 118, 81, 72, 56, 47, 120, 54, 81, 89, 57, 56, 77, 98, 113, 53, 118, 53, 56, 87, 118, 109, 10, 73, 121, 111, 66, 67, 110, 110, 53, 43, 107, 48, 52, 102, 78, 122, 71, 55, 75, 117, 70, 75, 48, 112, 84, 110, 115, 74, 81, 69, 48, 103, 107, 55, 100, 87, 105, 120, 47, 76, 84, 122, 116, 56, 110, 55, 121, 84, 90, 88, 79, 66, 84, 121, 68, 55, 75, 117, 102, 103, 115, 53, 73, 98, 76, 10, 89, 108, 79, 71, 49, 43, 107, 103, 43, 66, 68, 78, 97, 104, 84, 100, 98, 109, 87, 97, 114, 107, 84, 109, 55, 57, 114, 113, 74, 104, 98, 108, 79, 79, 112, 48, 113, 117, 83, 80, 65, 111, 71, 66, 65, 75, 121, 119, 112, 70, 43, 102, 49, 69, 111, 49, 101, 97, 52, 79, 70, 110, 119, 68, 10, 70, 115, 107, 103, 52, 65, 73, 112, 98, 119, 69, 106, 52, 109, 87, 99, 111, 112, 80, 113, 71, 66, 49, 66, 76, 57, 120, 110, 81, 113, 78, 68, 53, 104, 67, 102, 117, 48, 102, 50, 87, 82, 113, 103, 47, 97, 122, 115, 76, 49, 105, 105, 111, 102, 84, 68, 118, 50, 43, 82, 69, 87, 89, 67, 10, 118, 72, 67, 104, 55, 54, 104, 118, 79, 87, 49, 89, 81, 122, 55, 104, 106, 82, 49, 51, 56, 105, 56, 97, 100, 84, 122, 102, 48, 71, 99, 83, 83, 119, 55, 108, 81, 86, 107, 112, 105, 113, 112, 89, 110, 84, 86, 103, 43, 82, 101, 106, 76, 81, 57, 109, 70, 101, 99, 72, 84, 54, 48, 114, 10, 101, 77, 50, 117, 71, 116, 53, 49, 120, 71, 74, 108, 79, 51, 111, 81, 97, 103, 121, 71, 89, 66, 50, 53, 65, 111, 71, 66, 65, 75, 112, 67, 48, 119, 116, 112, 100, 120, 120, 122, 49, 103, 119, 76, 66, 101, 66, 75, 76, 89, 51, 113, 116, 106, 49, 74, 108, 52, 47, 75, 105, 84, 77, 104, 10, 75, 86, 77, 75, 65, 52, 70, 55, 100, 103, 51, 80, 85, 112, 55, 90, 56, 78, 70, 78, 75, 112, 102, 82, 72, 115, 72, 79, 121, 100, 110, 80, 114, 72, 97, 113, 86, 79, 43, 74, 110, 106, 49, 75, 75, 67, 49, 116, 71, 87, 57, 114, 120, 49, 72, 107, 110, 48, 79, 43, 76, 67, 114, 79, 10, 49, 116, 99, 85, 50, 114, 75, 104, 52, 75, 121, 56, 78, 80, 97, 115, 108, 71, 77, 122, 70, 111, 113, 56, 51, 114, 106, 120, 74, 86, 115, 67, 69, 110, 76, 43, 79, 120, 67, 121, 50, 72, 101, 114, 76, 43, 88, 69, 85, 117, 88, 75, 86, 122, 117, 57, 54, 90, 66, 112, 78, 50, 107, 97, 10, 56, 99, 89, 73, 77, 66, 53, 100, 65, 111, 71, 65, 97, 105, 51, 83, 118, 84, 97, 113, 54, 85, 100, 113, 43, 114, 73, 99, 52, 110, 55, 73, 101, 48, 101, 84, 65, 72, 89, 70, 84, 79, 72, 68, 43, 87, 104, 107, 100, 107, 70, 115, 118, 70, 73, 111, 116, 89, 76, 100, 72, 99, 77, 122, 10, 74, 98, 105, 108, 57, 100, 121, 113, 119, 76, 118, 73, 102, 72, 68, 86, 85, 87, 117, 87, 75, 67, 56, 103, 80, 121, 119, 110, 84, 74, 82, 111, 87, 71, 50, 111, 68, 71, 100, 110, 107, 114, 115, 53, 114, 73, 107, 89, 66, 78, 122, 53, 86, 97, 79, 119, 111, 80, 108, 80, 82, 119, 74, 65, 10, 68, 106, 79, 71, 70, 118, 103, 79, 115, 97, 101, 81, 119, 52, 108, 97, 65, 90, 52, 86, 55, 111, 68, 53, 116, 70, 49, 101, 43, 81, 80, 89, 79, 69, 79, 53, 108, 115, 108, 53, 56, 78, 81, 106, 106, 52, 67, 43, 120, 116, 86, 81, 72, 109, 48, 61, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 82, 83, 65, 32, 80, 82, 73, 86, 65, 84, 69, 32, 75, 69, 89, 45, 45, 45, 45, 45, 10}
	serverData, _ := pem.Decode([]byte(serverPEMbytes))

	// fmt.Println(" Log -- serverData", serverPEMbytes, string(serverPEMbytes))

	serverKey, err := x509.ParsePKCS1PrivateKey(serverData.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ServerKey = serverKey

	// fmt.Println(" Log -- client keys + server pubKey is loaded")
	machineStats := GetMachineUID()
	randomID := GetUUIDLite()
	machineID := machineStats["machine_id"].(string)
	payloadBytes := []byte(randomID + "," + machineID)

	// fmt.Println(" Log -- input bytes array", payloadBytes)

	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, &ServerKey.PublicKey, payloadBytes, nil)
	if err != nil {
		fmt.Println("Log -- error", ciphertext)
	}

	hash = sha512.New()
	_, err9 := rsa.DecryptOAEP(hash, rand.Reader, ServerKey, ciphertext, nil)
	if err9 != nil {
		fmt.Println(" Log -- error decrypting ..")
	}
	// fmt.Println(" Log -- decrypted data = ", string(plaintext))

	sEnc := b64.StdEncoding.EncodeToString(ciphertext)
	// bytesArray, err := b64.StdEncoding.DecodeString(sEnc)

	// fmt.Println(" Log -- input base64 ", sEnc)
	// fmt.Println(" Log -- input bytearray ", ciphertext)
	// fmt.Println(" Log -- input output bytearray", bytesArray)

	// resp, err := http.Get("http://192.168.0.113:8081")

	resp, err := http.Get(os.Getenv("LICENSE_SERVER") + "/validate?token=" + url.QueryEscape(sEnc))
	if err != nil {
		fmt.Println("Cannot connect license server")
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Cannot connect license server")
			return false
		}
		// bodyString := string(bodyBytes)
		// fmt.Println(" Log -- encrypted final solution recieved", bodyString)

		// decrypt final solution and check for OK
		hash = sha512.New()
		plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, ClientKey, bodyBytes, nil)
		if err != nil {
			fmt.Println(" Log -- error decrypting ..")
		}

		finalDecryptedString := string(plaintext)
		decryptedBits := strings.Split(finalDecryptedString, ",")

		if len(decryptedBits) == 3 && decryptedBits[0] == randomID && decryptedBits[1] == machineID && decryptedBits[2] == "OK" {
			return true
		}

		return false

	} else {
		fmt.Println("resp.StatusCode", resp.StatusCode)
	}
	return false

}

//GetValidPayload will send the useful post payload to microservice execution
func GetValidPayload(reqBody map[string]interface{}, validKeys []string) map[string]interface{} {
	cleanedPayload := make(map[string]interface{})
	for _, key := range validKeys {
		if reqBody[key] != nil && reqBody[key] != "" {
			cleanedPayload[key] = reqBody[key]
		}
	}
	return cleanedPayload
}

// PathExists -- check if path exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
