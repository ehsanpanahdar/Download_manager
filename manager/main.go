package manager

import (
	"os"
	"sync"
	"DOWNLOAD_MANAGER/http"
)

func Manage_goroutines( total_size , max_chunk_size int , file_fd *os.File , url string ) {
	var wg sync.WaitGroup
	iterrates := total_size / max_chunk_size
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
	if( last_downloaded_byte > int(total_size) ) {
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