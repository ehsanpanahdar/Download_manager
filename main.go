package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

func main() {
	const max_chunk_size  = 10000
	url := Get_url()
	total_data_size := Init_download(url)
	
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