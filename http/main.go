package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	user_agent string = "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0"
)

type Worker_info struct {
	Start_byte , End_byte int
	URL string
	WG *sync.WaitGroup
	Worker_id int
	File_fd *os.File
}

func Init_download(url string) (int64 , bool) {
	var client http.Client
	resp , err := client.Head(url)
	if( err != nil ) {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// AR stands for Accept-Ranges
	var Is_AR_supported bool
	accept_ranges := resp.Header.Get("Accept-Ranges")
	if( accept_ranges != "none" ) {
		Is_AR_supported = true
	} else {
		Is_AR_supported = false
	}

	return resp.ContentLength , Is_AR_supported
}

func Download_chunk( info Worker_info ) {
	if( info.WG != nil ) {
		defer info.WG.Done()
	}

	request , _ := http.NewRequest( "GET" , info.URL , nil )
	range_bytes := fmt.Sprintf( "bytes=%d-%d" , info.Start_byte , info.End_byte )
	request.Header.Add( "Range" , range_bytes )
	request.Header.Add( "User-Agent" , user_agent )

	var client http.Client
	resp , err := client.Do(request)
	if( err != nil ) {
		log.Printf( "Worker with ID %d was unable to download:%s" , info.Worker_id , err.Error() )
		return
	}

	buffer , err := io.ReadAll(resp.Body)
	if( err != nil ) {
		log.Printf( "Worker with ID %d was unable to copy from response body to memory:%s" , info.Worker_id , err.Error() )
		return
	}

	_ , err = info.File_fd.WriteAt( buffer , int64(info.Start_byte) )
	if( err != nil ) {
		log.Printf( "Worker with ID %d was unable to copy data from memory to file:%s" , info.Worker_id , err.Error() )
		return 
	}

	log.Printf( "Worker with ID %d finished successfuly" , info.Worker_id )
}