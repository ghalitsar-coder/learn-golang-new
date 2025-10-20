package main

import (
	"fmt"
	"sync"
	"time"
)

// Simulasi real-world scenario: E-commerce price comparison
func simulateEcommercePriceComparison() {
	fmt.Println("\n🛒 SIMULASI: E-COMMERCE PRICE COMPARISON")
	fmt.Println("======================================")

	fmt.Println("Scenario: 1 user search 'iPhone 15'")
	fmt.Println("Response time target: < 2 seconds")

	start := time.Now()

	// Channel untuk jobs
	jobs := make(chan string, 100)
	results := make(chan string, 100)

	// Worker pool (3 workers untuk demo)
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				// Simulasi API call ke marketplace
				time.Sleep(200 * time.Millisecond) // API latency
				result := fmt.Sprintf("Worker-%d: %s completed", workerID, job)
				results <- result
			}
		}(i)
	}

	// Jobs yang dibutuhkan untuk 1 user search
	searchJobs := []string{
		"Scrape Tokopedia",
		"Scrape Shopee",
		"Scrape Blibli",
		"Scrape Bukalapak",
		"Scrape Lazada",
		"Parse Tokopedia results",
		"Parse Shopee results",
		"Parse Blibli results",
		"Parse Bukalapak results",
		"Parse Lazada results",
		"Image processing - product 1",
		"Image processing - product 2",
		"Image processing - product 3",
		"Image processing - product 4",
		"Image processing - product 5",
		"Price comparison analysis",
		"Sort by relevance",
		"Filter by rating",
		"Calculate best deals",
		"Update cache",
	}

	fmt.Printf("Mengirim %d jobs untuk 1 user search...\n", len(searchJobs))

	// Send jobs
	for _, job := range searchJobs {
		jobs <- job
	}
	close(jobs)

	// Wait dan collect results
	wg.Wait()
	close(results)

	// Count results
	var resultCount int
	for range results {
		resultCount++
	}

	duration := time.Since(start)
	fmt.Printf("✅ Selesai: %d jobs dalam %v\n", resultCount, duration)
	fmt.Printf("💡 1 user search = %d concurrent jobs\n", len(searchJobs))
}

// Simulasi notification system
func simulateNotificationSystem() {
	fmt.Println("\n🔔 SIMULASI: NOTIFICATION SYSTEM")
	fmt.Println("==============================")

	fmt.Println("Scenario: 1 instructor post announcement ke 200 students")

	start := time.Now()

	students := 200
	notificationTypes := []string{"email", "push", "sms", "in-app"}

	// Hitung total jobs
	totalJobs := students * len(notificationTypes)

	fmt.Printf("Total jobs yang dibutuhkan: %d\n", totalJobs)
	fmt.Printf("Breakdown: %d students × %d notification types = %d jobs\n",
		students, len(notificationTypes), totalJobs)

	// Worker pool untuk handle notifications
	jobs := make(chan string, totalJobs)
	var wg sync.WaitGroup

	// 5 workers untuk notification
	numWorkers := 5
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				_ = job // Use the variable
				// Simulasi send notification
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}

	// Generate jobs
	for studentID := 1; studentID <= students; studentID++ {
		for _, notifType := range notificationTypes {
			job := fmt.Sprintf("Send %s to student-%d", notifType, studentID)
			jobs <- job
		}
	}
	close(jobs)

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("✅ Selesai: %d notifications dalam %v\n", totalJobs, duration)
	fmt.Printf("💡 1 instructor action = %d concurrent jobs\n", totalJobs)
}

// Simulasi batch processing
func simulateBatchProcessing() {
	fmt.Println("\n📊 SIMULASI: DAILY REPORT GENERATION")
	fmt.Println("==================================")

	fmt.Println("Scenario: Generate daily sales report untuk startup kecil")

	start := time.Now()

	// Data untuk diproses
	orders := 500    // 500 orders per day
	sellers := 50    // 50 sellers
	customers := 200 // 200 customers

	// Jobs breakdown
	jobs := []struct {
		name  string
		count int
	}{
		{"Process order", orders},
		{"Calculate seller commission", sellers},
		{"Generate customer invoice", customers},
		{"Update analytics dashboard", 20},
		{"Send seller reports", sellers},
		{"Backup transaction data", 10},
		{"Update search indexes", 30},
		{"Generate charts", 15},
	}

	var totalJobs int
	fmt.Println("Job breakdown:")
	for _, job := range jobs {
		totalJobs += job.count
		fmt.Printf("  • %s: %d jobs\n", job.name, job.count)
	}

	fmt.Printf("\nTotal: %d jobs untuk daily report\n", totalJobs)

	// Simulasi processing dengan worker pool
	jobChan := make(chan string, totalJobs)
	var wg sync.WaitGroup

	// 8 workers untuk batch processing
	numWorkers := 8
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobChan {
				_ = job // Use the variable
				// Simulasi processing time
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Send all jobs
	for _, jobType := range jobs {
		for i := 0; i < jobType.count; i++ {
			jobChan <- fmt.Sprintf("%s-%d", jobType.name, i+1)
		}
	}
	close(jobChan)

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("✅ Daily report selesai dalam %v\n", duration)
	fmt.Printf("💡 Background job = %d concurrent jobs\n", totalJobs)
}

// Simulasi API integration burst
func simulateAPIIntegration() {
	fmt.Println("\n🔗 SIMULASI: API INTEGRATION BURST")
	fmt.Println("================================")

	fmt.Println("Scenario: 50 users refresh bank balance bersamaan")

	start := time.Now()

	users := 50
	apisPerUser := []string{
		"Bank BCA API",
		"Bank BRI API",
		"Bank Mandiri API",
		"Credit Card API",
		"E-wallet OVO API",
		"E-wallet GoPay API",
		"Investment API",
		"Crypto API",
	}

	totalJobs := users * len(apisPerUser)

	fmt.Printf("%d users × %d APIs = %d concurrent API calls\n",
		users, len(apisPerUser), totalJobs)

	// Worker pool untuk API calls
	jobs := make(chan string, totalJobs)
	var wg sync.WaitGroup

	// 10 workers untuk API calls (I/O bound)
	numWorkers := 10
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				_ = job // Use the variable
				// Simulasi API call latency
				time.Sleep(300 * time.Millisecond)
			}
		}(i)
	}

	// Generate API calls
	for userID := 1; userID <= users; userID++ {
		for _, api := range apisPerUser {
			job := fmt.Sprintf("User-%d call %s", userID, api)
			jobs <- job
		}
	}
	close(jobs)

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("✅ All API calls completed dalam %v\n", duration)
	fmt.Printf("💡 Peak concurrent jobs: %d\n", totalJobs)
}

func main() {
	fmt.Println("🎯 REAL WORLD: KAPAN STARTUP KECIL BUTUH 1000+ JOBS?")
	fmt.Println("===================================================")

	fmt.Println(`
🏢 SETTING: Startup dengan 5-10 developer, 100-500 active users
💡 TERNYATA: Sangat mudah mencapai 1000+ concurrent jobs!
`)

	// Berbagai simulasi real-world scenarios
	simulateEcommercePriceComparison()
	simulateNotificationSystem()
	simulateBatchProcessing()
	simulateAPIIntegration()

	fmt.Println(`
📈 SUMMARY - SUMBER 1000+ JOBS DI STARTUP KECIL:
===============================================

🛒 E-COMMERCE FEATURES:
  • Price comparison: 20+ jobs per search
  • Product sync: 1000+ jobs per vendor
  • Image processing: 10+ jobs per upload
  • Inventory updates: 500+ jobs per sync

🔔 NOTIFICATION SYSTEMS:  
  • Email campaigns: 1 email = 1000+ jobs
  • Push notifications: 1 broadcast = 500+ jobs
  • SMS alerts: 1 trigger = 200+ jobs

📊 ANALYTICS & REPORTS:
  • Daily reports: 800+ jobs per day
  • Real-time dashboards: 100+ jobs per update
  • Data processing: 1000+ jobs per batch

🔗 API INTEGRATIONS:
  • Payment gateways: 50+ jobs per checkout
  • Social media sync: 100+ jobs per post
  • External data: 200+ jobs per refresh

⚡ BACKGROUND TASKS:
  • Cache warming: 500+ jobs
  • Database cleanup: 300+ jobs  
  • File processing: 1000+ jobs per batch
  • Search indexing: 200+ jobs per update

🚨 PEAK SCENARIOS:
  • Black Friday: 10,000+ jobs per minute
  • Viral content: 5,000+ jobs per share
  • System maintenance: 20,000+ jobs per cleanup
  • Data migration: 100,000+ jobs per operation

💡 CONCLUSION:
=============
Even dengan 100 concurrent users, startup mudah mencapai:
• 1,000+ jobs dalam normal operation  
• 10,000+ jobs dalam peak times
• 100,000+ jobs dalam maintenance windows

🏆 WORKER POOL BUKAN LUXURY, TAPI NECESSITY!
   Tanpa worker pool = System crash guaranteed! 💥
`)
}
