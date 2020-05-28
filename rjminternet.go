package gwcommon

import (
	"fmt"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
)

func DownloadFile(filepath string, url string) error {

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(filepath, url)

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			sProg := fmt.Sprintf("%.1f", 100*resp.Progress())
			//fmt.Println(sProg + "% complete...")
			fmt.Printf("\r" + sProg + "%% complete...")

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		return err
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)

	return nil
}

// func DownloadFile(url string, dest string) {

// 	file := path.Base(url)

// 	log.Printf("Downloading file %s from %s\n", file, url)

// 	var path bytes.Buffer
// 	path.WriteString(dest)
// 	path.WriteString("/")
// 	path.WriteString(file)

// 	start := time.Now()

// 	out, err := os.Create(path.String())

// 	if err != nil {
// 		fmt.Println(path.String())
// 		panic(err)
// 	}

// 	defer out.Close()

// 	headResp, err := http.Head(url)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer headResp.Body.Close()

// 	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

// 	if err != nil {
// 		panic(err)
// 	}

// 	done := make(chan int64)

// 	go PrintDownloadPercent(done, path.String(), int64(size))

// 	resp, err := http.Get(url)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer resp.Body.Close()

// 	n, err := io.Copy(out, resp.Body)

// 	if err != nil {
// 		panic(err)
// 	}

// 	done <- n

// 	elapsed := time.Since(start)
// 	log.Printf("Download completed in %s", elapsed)
// }

// func PrintDownloadPercent(done chan int64, path string, total int64) {

// 	var stop bool = false

// 	for {
// 		select {
// 		case <-done:
// 			stop = true
// 		default:

// 			file, err := os.Open(path)
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			fi, err := file.Stat()
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			size := fi.Size()

// 			if size == 0 {
// 				size = 1
// 			}

// 			var percent float64 = float64(size) / float64(total) * 100

// 			fmt.Printf("%.0f", percent)
// 			fmt.Println("%")
// 		}

// 		if stop {
// 			break
// 		}

// 		time.Sleep(time.Second)
// 	}
// }
