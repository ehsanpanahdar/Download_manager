package main

import (
	"DOWNLOAD_MANAGER/http"
	"DOWNLOAD_MANAGER/manager"
	"DOWNLOAD_MANAGER/utils"
	"log"
	"path"
)

func main() {
	const max_chunk_size  = 1024 * 1024  
	url := utils.Get_url()
	total_data_size , Is_AR_supported := http.Init_download(url) // AR stands for Accept-Ranges
	file_name := path.Base(url)

	file_fd , err := utils.Create_file( file_name , int(total_data_size) )
	if( err != nil ) {
		log.Fatal(err)
	}
	defer file_fd.Close()

	if( Is_AR_supported == false ) {
		var info http.Worker_info
		info.Start_byte = 0
		info.End_byte = int(total_data_size) - 1
		info.File_fd = file_fd
		info.URL = url
		info.Worker_id = 0 

		http.Download_chunk(info)
	} else {
		manager.Manage_goroutines( int(total_data_size) , max_chunk_size , file_fd , url )
	}
}