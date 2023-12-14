package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
)

type Worker_info struct {
	Start_byte , End_byte int
	URL , File_name string
	WG *sync.WaitGroup
	Worker_id int
}

func main() {
	const max_chunk_size  = 10000
	url := Get_url()
	total_data_size := Init_download(url)
	file_name := path.Base(url)

	var wg sync.WaitGroup
	worker_id := 0
	for temp:=0 ; ; worker_id++ {

		if( temp > (int(total_data_size)+max_chunk_size) ) {
			break
		}

		var worker Worker_info
		worker.Start_byte = temp
		temp += max_chunk_size -1
		worker.End_byte = temp
		worker.URL = url
		worker.WG = &wg
		worker.Worker_id = worker_id
		worker.File_name = file_name
		temp++

		wg.Add(1)
		go Download_chunk(worker)
	}

	wg.Wait()
}

func Get_url() string {
	fmt.Print( "Enter the URL: " )
	var URL string
	_ , err := fmt.Scan(&URL)
	if( err != nil ) {
		log.Fatal(err)
	}

	_ , err = url.ParseRequestURI(URL)
	if( err != nil ) {
		URL = "https://" + URL
		_ , err = url.ParseRequestURI(URL)
		if( err != nil ) {
			log.Fatal(err)
		}
	}
	
	return URL
}

func Init_download(url string) int64 {
	client := http.Client{}
	resp , err := client.Head(url)
	if( err != nil ) {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if( resp.StatusCode != 200 ) {
		log.Printf( "HTTP HEAD request was unsuccessful with status code %d" , resp.StatusCode )
		os.Exit(1)
	}

	return resp.ContentLength
}

func Download_chunk( info Worker_info ) {
	defer info.WG.Done()

	request , _ := http.NewRequest( "GET" , info.URL , nil )
	range_bytes := fmt.Sprintf( "byte=%d-%d" , info.Start_byte , info.End_byte )
	request.Header.Add( "Range" , range_bytes )

	var client http.Client
	resp , err := client.Do(request)
	if( err != nil ) {
		log.Printf( "Worker with ID %d was unable to download:%s" , info.Worker_id , err.Error() )
		return
	}

	file_name := fmt.Sprintf( "%s-%d" , info.File_name , info.Worker_id )
	file , err := os.Create(file_name)
	if( err != nil ) {
		log.Printf( "Worker with ID %d was unable to create temp file:%s" , info.Worker_id , err.Error() ) 
		return
	}

	_ , err = io.Copy( file , resp.Body )
	if( err != nil ) {
		log.Printf( "Worker with ID %d was unable to save data to temp file:%s" , info.Worker_id , err.Error() )
		return
	}

	log.Printf( "Worker with ID %d finished successfuly" , info.Worker_id )
}