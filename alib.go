package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func StrSliceDifference(slice1 []string, slice2 []string) []string {
	var diff []string
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}
	return diff
}

func StrContains(input, substring string) bool {
	return strings.Contains(input, substring)
}

func StrContainsInList(list []string, substring string) bool {
	for _, item := range list {
		if strings.Contains(item, substring) {
			return true
		}
	}
	return false
}

func StrJoin(input []string, seperator string) string {
	return strings.Join(input, seperator)
}

func StrSplit(input, seperator string) []string {
	return strings.Split(input, seperator)
}

func StrStrim(s string) string {
	return strings.Trim(s, " ")
}

func StrToUpper(s string) string {
	return strings.ToUpper(s)
}

func StrToLower(s string) string {
	return strings.ToLower(s)
}

func StrRjust(s string, l int) string {
	if len(s) < l {
		rt := l - len(s)
		return strings.Repeat(" ", rt) + s
	}

	return s
}

func StrLjust(s string, l int) string {
	if len(s) < l {
		rt := l - len(s)
		return s + strings.Repeat(" ", rt)
	}

	return s
}

func StrIsDigint(s string) bool {
	c := "0123456789."
	for _, pos := range s {
		if !strings.Contains(c, string(pos)) {
			return false
		}
	}
	return true
}

func StrRepeat(s string, count int) string {
	return strings.Repeat(s, count)
}

func StrIndex(s string, sub string) int {
	return strings.Index(s, sub)
}

func StrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func StrToInt64(str string) (int64, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func StrToFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func FloatToStr(fl float64) string {
	return strconv.FormatFloat(fl, 'f', 8, 64)
}

func FloatToStrWithLeng(fl float64, l int) string {
	fl_str := FloatToStr(fl)
	fl_ls := StrSplit(fl_str, ".")
	if len(fl_ls) < 2 {
		return fl_str
	}
	var second string
	second_all := fl_ls[1]
	if len(second_all) > l {
		second = second_all[0:l]
	} else {
		second = second_all
	}
	return fl_ls[0] + "." + second
}

func HTTPSendRequest(method, path string, headers map[string]string, body io.Reader) (string, error) {
	result := strings.ToUpper(method)

	if result != "POST" && result != "GET" && result != "DELETE" {
		return "", errors.New("Invalid HTTP method specified.")
	}

	req, err := http.NewRequest(method, path, body)

	if err != nil {
		return "", err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		return "", err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func HTTPSendGetRequest(url string, jsonDecode bool, result interface{}) (err error) {
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		fmt.Printf("HTTP status code: %d\n", res.StatusCode)
		return errors.New("Status code was not 200.")
	}

	contents, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if jsonDecode {
		err := JSONDecode(contents, &result)
		if err != nil {
			fmt.Println(string(contents[:]))
			return err
		}
	} else {
		result = &contents
	}

	return nil
}

func HTTPSendGetRequestTimeout(url string, jsonDecode bool, result interface{}, seconds time.Duration) (err error) {
	timeout := time.Duration(seconds * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Get(url)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("HTTP status code: %d\n", res.StatusCode)

		return errors.New("Status code was not 200.")
	}

	contents, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if jsonDecode {
		err := JSONDecode(contents, &result)
		if err != nil {
			fmt.Println(string(contents[:]))
			return err
		}
	} else {
		result = &contents
	}

	return nil
}

func JSONEncode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func JSONDecode(data []byte, to interface{}) error {
	return json.Unmarshal(data, to)
}

func EncodeURLValues(url string, values url.Values) string {
	path := url
	if len(values) > 0 {
		path += "?" + values.Encode()
	}
	return path
}

func TimeStampToTime(timeint64 int64) time.Time {
	return time.Unix(timeint64, 0)
}

func TimeStampStrToTime(timeStr string) (time.Time, error) {
	i, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(i, 0), nil
}

func TimeStampNow() int64 {
	now := time.Now()
	return now.Unix()
}

func TimeToString(dt time.Time) string {
	tlayout := "2006-01-02 15:04:05.000000000"
	return dt.Format(tlayout)
}

func TimeToStringWithoutNanosec(dt time.Time) string {
	tlayout := "2006-01-02 15:04:05"
	return dt.Format(tlayout)
}

func TimeStrToTime(str_time string) (time.Time, error) {
	if str_time == "" {
		return time.Time{}, nil
	}
	tlayout := "2006-01-02 15:04:05.000000000"
	return time.Parse(tlayout, str_time)
}

func TimeStrWithoutNanosecToTime(str_time string) (time.Time, error) {
	if str_time == "" {
		return time.Time{}, nil
	}
	tlayout := "2006-01-02 15:04:05"
	return time.Parse(tlayout, str_time)
}

func FileRead(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func FileWrite(file string, data []byte) error {
	err := ioutil.WriteFile(file, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func FileList(dir string) []os.FileInfo {
	ls, _ := ioutil.ReadDir(dir)
	return ls
}

// GetURIPath returns the path of a URL given a URL
func GetURIPath(uri string) string {
	urip, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	if urip.RawQuery != "" {
		return fmt.Sprintf("%s?%s", urip.Path, urip.RawQuery)
	}
	return urip.Path
}

func IToStr(i interface{}) string {
	if i == nil {
		return ""
	}
	switch i2 := i.(type) {
	default:
		return fmt.Sprint(i2)
	case []uint8:
		return string(i2)
	case int:
		return strconv.Itoa(i2)
	case int64:
		return strconv.FormatInt(i2, 10)
	case bool:
		if i2 {
			return "true"
		} else {
			return "false"
		}
	case string:
		return i2
	case *bool:
		if i2 == nil {
			return ""
		}
		if *i2 {
			return "true"
		} else {
			return "false"
		}
	case *string:
		if i2 == nil {
			return ""
		}
		return *i2
	case *json.Number:
		return i2.String()
	case json.Number:
		return i2.String()
	}
	return ""
}
