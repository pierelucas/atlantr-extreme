package proxy

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

type Proxy struct {
	URL     string
	Port    int
	Counter int
}

type socks5 struct {
	Prox  string
	Valid bool
}

//restProxy rests a proxy for a period in to mitigate rate limiting
func restProxy(proxy Proxy, proxies chan Proxy, restTime time.Duration, cURL string, httpTimeout time.Duration) {
	start := time.Now()
	err := checkSocks(proxy.URL+":"+strconv.Itoa(proxy.Port), cURL, httpTimeout)

	if err != nil {
		log.Println("socks5 " + proxy.URL + " was evicted because it does not work anymore")
		return
	}

	elapsed := time.Since(start)
	restTimeNew := restTime - elapsed
	time.Sleep(restTimeNew)

	proxy.Counter = 0
	proxies <- proxy
}

func GetProxy(proxies chan Proxy, restTime time.Duration, interval int, cURL string,
	httpTimeout time.Duration, breaker int) (Proxy, error) {
	counter := 0

	for {
		select {
		case p := <-proxies:
			if p.Counter >= interval {
				go restProxy(p, proxies, restTime, cURL, httpTimeout)
				break
			}
			p.Counter++
			proxies <- p

			return p, nil
		default:
			time.Sleep(restTime)
			counter++

			if counter > breaker {
				return Proxy{}, fmt.Errorf("probalby all proxies are dead")
			}
		}
	}
}

func removeDuplicatesFromSlice(s []string) []string {
	m := make(map[string]bool)

	for _, item := range s {
		if _, ok := m[item]; ok {
		} else {
			m[item] = true
		}
	}

	result := make([]string, 0, len(m))

	for item := range m {
		result = append(result, item)
	}

	return result
}

func checkSocks(socks5Address string, url string, timeOut time.Duration) error {
	dialer, err := proxy.SOCKS5("tcp", socks5Address, nil, proxy.Direct)
	if err != nil {
		return err
	}

	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport, Timeout: timeOut}

	//dial is deprecated but the proxy package does not support it
	httpTransport.Dial = dialer.Dial
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if strings.Contains(socks5Address, string(b)) {
		return nil
	}

	return fmt.Errorf("ip not matching")
}

func validateIPport(str1 string) error {
	str := strings.Split(str1, ":")
	if len(str) != 2 {
		return fmt.Errorf("bad format")
	}

	ip := net.ParseIP(str[0])

	if ip.To4() != nil {
		port, err := strconv.Atoi(str[1])
		if err != nil {
			return fmt.Errorf("bad Port")
		}

		if port > 65535 || port < 0 {
			return fmt.Errorf("bad Port: %d", port)
		}

		return nil
	}

	return fmt.Errorf("bad ip: %s", str[0])
}

func readInfile(path string) ([]string, error) {
	var lines []string
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())
		err := validateIPport(l)

		if err != nil {
			continue
		}

		lines = append(lines, l)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("no socks to read in file")
	}

	return lines, nil
}

func worker(jobs <-chan string, results chan<- socks5, ipCheckURL string, timeOut time.Duration) {
	for j := range jobs {
		err := checkSocks(j, ipCheckURL, timeOut)
		if err != nil {
			results <- socks5{Prox: j, Valid: false}
			continue
		}
		results <- socks5{Prox: j, Valid: true}
	}
}

// InitSocks returns read-only channel where we can read our valid proxies from and len(n)
func InitSocks(pathSocks5 string, socksCheckWorker int, cURL string, checkTimeout time.Duration) (<-chan Proxy, []string, int) {
	lll, err := readInfile(pathSocks5)
	if err != nil {
		log.Fatal(err)
	}

	l := removeDuplicatesFromSlice(lll)
	lenL := len(l)
	jobs := make(chan string, lenL)
	results := make(chan socks5, lenL)

	for w := 1; w <= socksCheckWorker; w++ {
		go worker(jobs, results, cURL, checkTimeout)
	}

	for _, v := range l {
		jobs <- v
	}

	close(jobs)

	var validSocks []string

	for a := 1; a <= lenL; a++ {
		r := <-results
		if r.Valid {
			validSocks = append(validSocks, r.Prox)
		}
	}

	// delete non-unique socks from validSocks slice
	validSocks = removeDuplicatesFromSlice(validSocks)
	lenValid := len(validSocks)

	if lenValid < 1 {
		log.Fatal("no Valid proxies found")
	}

	proxies := make(chan Proxy, lenValid)

	for _, v := range validSocks {
		temp := strings.Split(v, ":")
		url := temp[0]
		port, err := strconv.Atoi(temp[1])
		if err != nil {
			log.Println(err)
			continue
		}
		proxies <- Proxy{URL: url, Port: port, Counter: 0}
	}

	return proxies, validSocks, lenValid
}
