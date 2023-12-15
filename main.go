package main

import (
	"log"
	"DOWNLOAD_MANAGER/utils"
	"DOWNLOAD_MANAGER/http"
	"path"
	"sync"
)

func main() {
	const max_chunk_size  = 10000
	url := utils.Get_url()
	total_data_size := http.Init_download(url)
	file_name := path.Base(url)

	file , err := utils.Create_file( file_name , int(total_data_size) )
	if( err != nil ) {
		log.Fatal(err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	worker_id := 0
	for temp:=0 ; ; worker_id++ {

		if( temp > (int(total_data_size)+max_chunk_size) ) {
			break
		}

		var worker http.Worker_info
		worker.Start_byte = temp
		temp += max_chunk_size -1
		worker.End_byte = temp
		worker.URL = url
		worker.WG = &wg
		worker.Worker_id = worker_id
		worker.File_name = file_name
		worker.File = file
		temp++

		wg.Add(1)
		go http.Download_chunk(worker)
	}

	wg.Wait()
}