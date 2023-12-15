package utils

import (
	"fmt"
	"log"
	"net/url"
	"os"
)

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

func Create_file( file_name string , file_size int ) (*os.File , error) {
	file_handle , err := os.Create(file_name) 
	if( err != nil ) {
		return &os.File{} , err
	}

	err = file_handle.Truncate(int64(file_size))
	if( err != nil ) {
		return &os.File{} , err
	}

	return file_handle , nil 
} 