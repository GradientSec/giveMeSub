package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// --- Ranglar va Dizayn uchun Konstantalar ---
const (
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorReset  = "\033[0m"
)

// --- "giveMeSub" NOMINI ANIQ AKS ETTIRUVCHI BANNER ---
const banner = `
 ██████╗ ██╗██╗   ██╗███████╗    ███╗   ███╗███████╗  ███████╗ ██╗   ██╗██████╗ 
██╔════╝ ██║██║   ██║██╔════╝    ████╗ ████║██╔════╝  ██╔════╝ ██║   ██║██╔══██╗
██║  ███╗██║██║   ██║█████╗      ██╔████╔██║█████╗    ███████╗ ██║   ██║██████╔╝
██║   ██║██║╚██╗ ██╔╝██╔══╝      ██║╚██╔╝██║██╔══╝    ╚════██║ ██║   ██║██╔══██╗
╚██████╔╚██║ ╚████╔╝ ███████╗    ██║ ╚═╝ ██║███████╗  ███████╔╝╚██████╔╝██████╔╝
╚═════╝  ╚═╝  ╚═══╝  ╚══════╝    ╚═╝     ╚═╝╚══════╝   ╚═════╝  ╚═════╝ ╚═════╝ 

      giveMeSub - Discovering the surface at the speed of Go.
                       by GradientSec @ethica @awds
`

// Topilgan natijani saqlash uchun struktura
type ScanResult struct {
	Domain    string
	IPs       []net.IP
	OpenPorts []string
}

func main() {
	// --- Banner va Kirish Ma'lumotlari ---
	fmt.Println(ColorBlue + banner + ColorReset)
	fmt.Printf("%s[*]%s Dastur ishga tushirildi: %s\n", ColorYellow, ColorReset, time.Now().Format("2006-01-02 15:04:05"))

	// --- Foydalanuvchi Kiritadigan Parametrlar (Flags) ---
	targetDomain := flag.String("d", "", "Maqsad domen (masalan, example.com)")
	wordlistFile := flag.String("w", "", "Subdomenlar ro'yxati joylashgan fayl (masalan, subdomains.txt)")
	outputFile := flag.String("o", "found_subdomains.txt", "Natijalarni saqlash uchun fayl nomi")
	threads := flag.Int("t", 100, "Parallel ishchilar (potoklar) soni")
	flag.Parse()

	if *targetDomain == "" || *wordlistFile == "" {
		fmt.Println("\nNoto'g'ri buyruq! -d va -w parametrlari majburiy.")
		flag.Usage()
		return
	}
	fmt.Printf("%s[*]%s Maqsad: %s\n", ColorYellow, ColorReset, *targetDomain)
	fmt.Printf("%s[*]%s Wordlist: %s\n", ColorYellow, ColorReset, *wordlistFile)
	fmt.Printf("%s[*]%s Natijalar fayli: %s\n", ColorYellow, ColorReset, *outputFile)
	fmt.Printf("%s[*]%s Potoklar soni: %d\n\n", ColorYellow, ColorReset, *threads)


	// --- Fayllarni O'qish va Yozish ---
	file, err := os.Open(*wordlistFile)
	if err != nil {
		log.Fatalf("Wordlist faylni ochib bo'lmadi: %s", err)
	}
	defer file.Close()

	outFile, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Natija faylini ochib bo'lmadi: %s", err)
	}
	defer outFile.Close()

	// --- Parallel Vazifalarni Boshqarish ---
	var wg sync.WaitGroup
	jobs := make(chan string, *threads)
	results := make(chan ScanResult)

	// Natijalarni faylga va ekranga chiqarish uchun alohida goroutine
	go handleResults(results, outFile, &wg)

	// Ishchilarni (worker) ishga tushirish
	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go worker(&wg, jobs, results)
	}
	
	// Fayldan o'qib, vazifalarni ishchilarga yuborish
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subdomain := scanner.Text()
		fullDomain := fmt.Sprintf("%s.%s", subdomain, *targetDomain)
		jobs <- fullDomain
	}

	close(jobs)
	wg.Wait()
	close(results)

	fmt.Printf("\n%s[*]%s Skanerlash yakunlandi. Natijalar '%s' fayliga saqlandi.\n", ColorYellow, ColorReset, *outputFile)
}

// Natijalarni boshqaruvchi funksiya
func handleResults(results <-chan ScanResult, file *os.File, wg *sync.WaitGroup) {
	for result := range results {
		// Natijani formatlash
		ipString := strings.Trim(fmt.Sprintf("%v", result.IPs), "[]")
		portString := strings.Join(result.OpenPorts, ", ")
		if portString == "" {
			portString = "N/A"
		}

		// Ekranga chiqarish
		fmt.Printf("%s[+]%s Topildi: %-40s -> IP: [%-15s] -> Portlar: [%s]\n",
			ColorGreen, ColorReset, result.Domain, ipString, ColorYellow+portString+ColorReset)

		// Faylga yozish
		fileLine := fmt.Sprintf("%s,%s,%s\n", result.Domain, ipString, portString)
		file.WriteString(fileLine)

		wg.Done()
	}
}


// Asosiy ishchi funksiyasi
func worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- ScanResult) {
	for domain := range jobs {
		ips, err := net.LookupIP(domain)
		if err == nil {
			openPorts := checkWebServer(domain)
			
			result := ScanResult{
				Domain:    domain,
				IPs:       ips,
				OpenPorts: openPorts,
			}
			wg.Add(1) 
			results <- result
		}
	}
	wg.Done()
}


// Web-serverni tekshiruvchi funksiya
func checkWebServer(domain string) []string {
	var openPorts []string
	ports := []string{"80", "443"}

	for _, port := range ports {
		address := fmt.Sprintf("%s:%s", domain, port)
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err == nil {
			protocol := "HTTP"
			if port == "443" {
				protocol = "HTTPS"
			}
			openPorts = append(openPorts, fmt.Sprintf("%s:%s", protocol, port))
			conn.Close()
		}
	}
	return openPorts
}
