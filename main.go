package main

import (
	"DOWNLOAD_MANAGER/http"
	"DOWNLOAD_MANAGER/utils"
	"log"
	"path"
	"sync"
)

func main() {
	const max_chunk_size  = 1024 * 1024 * 35 
	url := utils.Get_url()
	total_data_size := http.Init_download(url)
	file_name := path.Base(url)

	file_fd , err := utils.Create_file( file_name , int(total_data_size) )
	if( err != nil ) {
		log.Fatal(err)
	}
	defer file_fd.Close()

	var wg sync.WaitGroup
	iterrates := total_data_size / max_chunk_size
	count := 0
	for id :=0 ; id < int(iterrates) ; id++ {
		var worker http.Worker_info
		worker.Start_byte = max_chunk_size * id
		worker.End_byte = (max_chunk_size * (id+1)) -1
		worker.Worker_id = id
		worker.File_fd = file_fd
		worker.URL = url
		worker.WG = &wg
		count++

		wg.Add(1)
		go http.Download_chunk(worker)
	}

	last_downloaded_byte := (max_chunk_size*(count+1)) -1
	if( last_downloaded_byte > int(total_data_size) ) {
		var worker http.Worker_info
		worker.Start_byte = max_chunk_size * count
		worker.End_byte = (max_chunk_size * (count+2)) -1
		worker.Worker_id = count
		worker.File_fd = file_fd
		worker.URL = url
		worker.WG = &wg

		wg.Add(1)
		go http.Download_chunk(worker)
	}

	wg.Wait()
}