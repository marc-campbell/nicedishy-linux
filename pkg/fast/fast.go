package fast

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	fastcli "github.com/gesquive/fast-cli/fast"
	fastformat "github.com/gesquive/fast-cli/format"
	"github.com/gesquive/fast-cli/meters"
	"github.com/marc-campbell/nicedishy-linux/pkg/version"
)

/*
https://github.com/gesquive/fast-cli/blob/master/main.go
(MIT licensed), but not exported code
*/

func Run() error {
	count := uint64(3)
	fastcli.UseHTTPS = true
	urls := fastcli.GetDlUrls(count)

	if len(urls) == 0 {
		urls = append(urls, fastcli.GetDefaultURL())
	}

	client := &http.Client{}
	count = uint64(len(urls))

	primaryBandwidthReader := meters.BandwidthMeter{}
	bandwidthMeter := meters.BandwidthMeter{}
	ch := make(chan *copyResults, 1)
	completed := uint64(0)

	for i := uint64(0); i < count; i++ {
		// Create the HTTP request
		request, err := http.NewRequest("GET", urls[i], nil)
		if err != nil {
			return err
		}
		request.Header.Set("User-Agent", version.Version())

		// Get the HTTP Response
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		// Set information for the leading index
		if i == 0 {
			tapMeter := io.TeeReader(response.Body, &primaryBandwidthReader)
			go asyncCopy(i, ch, &bandwidthMeter, tapMeter)
		} else {
			// Start reading
			go asyncCopy(i, ch, &bandwidthMeter, response.Body)
		}

	}

	for {
		select {
		case results := <-ch:
			if results.err != nil {
				fmt.Fprintf(os.Stdout, "\n%v\n", results.err)
				os.Exit(1)
			}

			completed++

			fmt.Printf("%s\n", fastformat.BitsPerSec(bandwidthMeter.Bandwidth()))
			return nil
		case <-time.After(100 * time.Millisecond):

		}
	}

}

type copyResults struct {
	index        uint64
	bytesWritten uint64
	err          error
}

func asyncCopy(index uint64, channel chan *copyResults, writer io.Writer, reader io.Reader) {
	bytesWritten, err := io.Copy(writer, reader)
	channel <- &copyResults{index, uint64(bytesWritten), err}
}

func sumArr(array []uint64) (sum uint64) {
	for i := 0; i < len(array); i++ {
		sum = sum + array[i]
	}
	return
}
