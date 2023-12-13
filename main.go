package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	url := Get_url()
	total_data_size := Init_download(url)
}

func Get_url() string {
	fmt.Print( "Enter the URL:" )
	var url string
	_ , err := fmt.Scan(&url)
	if( err != nil ) {
		log.Fatal(err)
	}
	
	return url
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