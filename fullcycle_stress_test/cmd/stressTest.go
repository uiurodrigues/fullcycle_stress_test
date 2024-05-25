/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

// stressTestCmd represents the stressTest command
var stressTestCmd = &cobra.Command{
	Use:   "stressTest",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		requests, _ := cmd.Flags().GetInt("requests")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		runStressTest(url, requests, concurrency)
	},
}

func init() {
	rootCmd.AddCommand(stressTestCmd)

	stressTestCmd.Flags().StringP("url", "u", "", "URL to test")
	stressTestCmd.MarkFlagRequired("url")

	stressTestCmd.Flags().IntP("requests", "r", 100, "Number of requests to send")
	stressTestCmd.MarkFlagRequired("requests")

	stressTestCmd.Flags().IntP("concurrency", "c", 10, "Number of concurrent requests")
	stressTestCmd.MarkFlagRequired("concurrency")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stressTestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stressTestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runStressTest(url string, requests int, concurrency int) {
	log.Println("Starting stress test...")
	defer log.Println("Stress test finished...")

	report := &Report{}
	report.StartExecution()

	parallelThreads := make(chan int, concurrency)
	wg := &sync.WaitGroup{}
	wg.Add(requests)

	log.Printf(">>>>>> Starting %d requests with concurrency %d \n", requests, concurrency)
	for i := 0; i < requests; i++ {
		parallelThreads <- 1
		go callEndpoint(url, parallelThreads, wg, report)
	}
	wg.Wait()

	report.EndExecution()
	report.Show()
}

func callEndpoint(url string, ch chan int, wg *sync.WaitGroup, report *Report) {
	defer wg.Done()

	report.TotalRequests.Add(1)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("error creating request. err:%s \n", err.Error())
		report.TotalUndefinedError.Add(1)
		<-ch
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusTooManyRequests:
				report.Total429.Add(1)
			case http.StatusNotFound:
				report.Total404.Add(1)
			case http.StatusInternalServerError:
				report.Total500.Add(1)
			default:
				report.TotalUndefinedError.Add(1)
			}
		} else {
			report.TotalUndefinedError.Add(1)
		}
		<-ch
		return
	}
	defer resp.Body.Close()

	report.Total200.Add(1)
	<-ch
}

type Report struct {
	executionStartTime  time.Time
	executionEndTime    time.Time
	TotalRequests       atomic.Int32
	Total200            atomic.Int32
	Total404            atomic.Int32
	Total429            atomic.Int32
	Total500            atomic.Int32
	TotalUndefinedError atomic.Int32
}

func (r *Report) StartExecution() {
	r.executionStartTime = time.Now()
}

func (r *Report) EndExecution() {
	r.executionEndTime = time.Now()
}

func (r *Report) Show() {
	log.Println(">>>>>> Execution Report <<<<<<")
	log.Printf("Total requests: %d \n", r.TotalRequests.Load())
	log.Printf("Total 200: %d \n", r.Total200.Load())
	log.Printf("Total 404: %d \n", r.Total404.Load())
	log.Printf("Total 429: %d \n", r.Total429.Load())
	log.Printf("Total 500: %d \n", r.Total500.Load())
	log.Printf("Total undefined errors: %d \n", r.TotalUndefinedError.Load())
	log.Printf("Execution time: %.f seconds \n", r.executionEndTime.Sub(r.executionStartTime).Seconds())
}
