package hw10programoptimization

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

var (
	errorCountDomains = errors.New("invalid input: domain is empty")
	errorReader       = errors.New("invalid input: reader is empty")
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

// оставил на будущее, как пример мучений с горутинами
// func oldGetUsers(r io.Reader) (result users, err error) {
//	if err != nil {
//		return
//	}
//	resCh := make(chan User, 10)
//	ch := make(chan []byte, 10)
//
//	numWorkers := runtime.GOMAXPROCS(0)
//	var wg sync.WaitGroup
//	for i := 0; i < numWorkers; i++ {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			for line := range ch {
//				var user User
//				if err := json.Unmarshal(line, &user); err != nil {
//					return
//				}
//				resCh <- user // Отправляем каждого пользователя сразу
//			}
//		}()
//	}
//	scanner := bufio.NewScanner(r)
//	scanner.Split(bufio.ScanLines)
//	go func() {
//		defer close(resCh)
//		for scanner.Scan() {
//			ch <- scanner.Bytes()
//		}
//	}()
//	wg.Wait()
//	//close(resCh)
//	for user := range resCh {
//		result = append(result, user)
//	}
//	return
//}

func getUsers(r io.Reader) (chan User, error) {
	if r == nil {
		return nil, errorReader
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	decoder := json.NewDecoder(r)
	userCh := make(chan User, runtime.GOMAXPROCS(0))

	go func() {
		defer close(userCh)
		var user User
		for {
			if err := decoder.Decode(&user); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				return
			}
			if user.Email != "" {
				userCh <- user
			}
		}
	}()
	return userCh, nil
}

func countDomains(u chan User, domain string) (DomainStat, error) {
	if domain == "" {
		return nil, errorCountDomains
	}

	result := make(DomainStat)
	numWorkers := runtime.GOMAXPROCS(0)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range u {
				if strings.Contains(user.Email, "."+domain) {
					parts := strings.SplitN(user.Email, "@", 2)
					if len(parts) != 2 {
						continue
					}
					domainPart := strings.ToLower(parts[1])

					mu.Lock()
					result[domainPart]++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	return result, nil
}
