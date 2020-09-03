package parse

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/pierelucas/atlantr-extreme/data"
)

// Hosters --
func Hosters(path string) (map[string]*data.Host, error) {
	var hosterData = make(map[string]*data.Host)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		temp := strings.Split(scanner.Text(), ":")
		host := data.NewHost(temp[1], temp[2])
		hosterData[temp[0]] = host
	}

	if len(hosterData) < 1 {
		return nil, fmt.Errorf("no hoster settings found")
	}

	return hosterData, nil
}

// Matchers --
func Matchers(path string) ([]string, error) {
	var matchers []string

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		matchers = append(matchers, strings.TrimSpace(scanner.Text()))
	}

	if len(matchers) < 1 {
		return nil, fmt.Errorf("no matchers settings found")
	}

	return matchers, nil
}

// UserPass --
func UserPass(str string) (string, string, error) {
	str2 := strings.Split(str, ":")
	lenStr2 := len(str2)

	if lenStr2 < 2 {
		return "", "", fmt.Errorf("no in [:] in %s", str)
	}

	user := strings.ToLower(str2[0])

	var pass string
	if lenStr2 > 2 {
		pass = strings.Join(str2[1:], ":")
	} else {
		pass = str2[1]
	}

	//RFC 2821
	if len(user) > 254 {
		return "", "", fmt.Errorf("user too long")
	}

	//validating username
	if !strings.Contains(user, "@") {
		return "", "", fmt.Errorf("no [@] in %s", str)
	}

	//validating username
	if !strings.Contains(user, ".") {
		return "", "", fmt.Errorf("no in [.] in %s", str)
	}

	return user, pass, nil
}

// LastLineLog --
func LastLineLog(path string) (int, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	nStr := strings.TrimSpace(strings.Trim(string(data), "\n"))

	return strconv.Atoi(nStr)
}
