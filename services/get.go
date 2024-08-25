package services

import (
	"net/http"
)

func Get(writer http.ResponseWriter, request *http.Request) error {
	storage := Storage{
		RootFolder: request.Context().Value(KeyDataFolder).(string),
	}

	bucketName, objectName := bucketNameAndObjectKey(request.URL.Path, request.Context().Value(KeyUrlContext).(string))

	data, err := storage.GetData(bucketName, objectName, "")
	if err != nil {
		// Handle statistics application
		statisticsApplicationFolder := request.Context().Value(KeyStatisticsApplicationFolder).(string)
		fs := http.FileServer(http.Dir(statisticsApplicationFolder))
		fs.ServeHTTP(writer, request)
		return err
	}

	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}
