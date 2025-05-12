package aol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type TomlDataStruct struct {
	FileContent string `json:"fileContent"`
	FileName    string `json:"fileName"`
	FilePath    string `json:"filePath"`
	Server      string `json:"server"`
}

type TomlResponseStruct struct {
	Status string `json:"status"`
	ErrMsg string `json:"errMsg"`
}

func SetTomlValue(w http.ResponseWriter, r *http.Request) {
	log.Println("SetTomlValue Started(+)")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT")
	(w).Header().Set("Access-Control-Allow-Headers", " Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "PUT" {
		var lFetchDataRec TomlDataStruct
		var lFetchResponseRec TomlResponseStruct

		lFetchResponseRec.Status = "S"

		lBody, lErr := ioutil.ReadAll(r.Body)
		if lErr != nil {
			lFetchResponseRec.Status = "E"
			lFetchResponseRec.ErrMsg = "Error in (stv01): " + lErr.Error()
		} else {
			//Unmarshall the body and save into FetchDataStruct Struct
			lErr = json.Unmarshal(lBody, &lFetchDataRec)
			if lErr != nil {
				lFetchResponseRec.Status = "E"
				lFetchResponseRec.ErrMsg = "Error in (stv02): " + lErr.Error()
			} else {
				log.Println("lFetchDataRec", lFetchDataRec)

				filePath := filepath.Join(lFetchDataRec.FilePath, lFetchDataRec.FileName)
				log.Println("filePath", filePath)

				// Check if file exists
				if _, lErr := os.Stat(filePath); os.IsNotExist(lErr) {
					// Create the file if it doesn't exist
					file, err := os.Create(filePath)
					if err != nil {
						lFetchResponseRec.Status = "E"
						lFetchResponseRec.ErrMsg = "Error creating file (stv03): " + err.Error()
					}
					defer file.Close()
				}

				// Write the content to the file
				err := ioutil.WriteFile(filePath, []byte(lFetchDataRec.FileContent), 0644)
				if err != nil {
					lFetchResponseRec.Status = "E"
					lFetchResponseRec.ErrMsg = "Error writing to file: " + err.Error()
				} else {
					lFetchResponseRec.Status = "S"
					lFetchResponseRec.ErrMsg = "File written successfully"
				}
			}
		}
		log.Println("lFetchResponseRec", lFetchResponseRec)
		//marshall the Response to front
		data, lErr := json.Marshal(lFetchResponseRec)
		if lErr != nil {
			fmt.Fprintf(w, "Error taking data"+lErr.Error())
		} else {
			fmt.Fprintf(w, string(data))
		}
	}
	log.Println("SetTomlValue Ended(-)")
}

// package aol

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// )

// type TomlDataStruct struct {
// 	FileContent string `json:"fileContent"`
// 	FileName    string `json:"fileName"`
// 	FilePath    string `json:"filePath"`
// 	Server      string `json:"server"`
// }

// type TomlResponseStruct struct {
// 	Status      string        `json:"status"`
// 	ErrMsg      string        `json:"errMsg"`
// 	FilePath    string        `json:"filePath"`
// 	FileNameArr []fileNamestr `json:"fileNameArr"`
// }

// func SetToml(w http.ResponseWriter, r *http.Request) {
// 	log.Println("SetToml Started(+)")
// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "PUT")
// 	(w).Header().Set("Access-Control-Allow-Headers", " Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	if r.Method == "PUT" {
// 		var lFetchDataRec TomlDataStruct
// 		var lFetchResponseRec TomlResponseStruct

// 		lFetchResponseRec.Status = "S"

// 		lbody, lerr := ioutil.ReadAll(r.Body)
// 		if lerr != nil {
// 			lFetchResponseRec.Status = "E"
// 			lFetchResponseRec.ErrMsg = "Error in (stv01): " + lerr.Error()
// 		} else {
// 			//Unmarshall the body and save into FetchDataStruct Struct
// 			lerr = json.Unmarshal(lbody, &lFetchDataRec)
// 			if lerr != nil {
// 				lFetchResponseRec.Status = "E"
// 				lFetchResponseRec.ErrMsg = "Error in (flf02): " + lerr.Error()
// 			} else {
// 				log.Println("lFetchDataRec", lFetchDataRec)

// 				filePath := lFetchDataRec.FilePath + "/" + lFetchDataRec.FileName
// 				log.Println("filePath", filePath)

// 				// Write the content to the file
// 				err := ioutil.WriteFile(filePath, []byte(lFetchDataRec.FileContent), 0644)
// 				if err != nil {
// 					lFetchResponseRec.Status = "E"
// 					lFetchResponseRec.ErrMsg = "Error writing to file: " + err.Error()
// 				} else {
// 					lFetchResponseRec.Status = "S"
// 					lFetchResponseRec.ErrMsg = "File written successfully"
// 				}
// 			}
// 		}
// 		log.Println("lFetchResponseRec", lFetchResponseRec)
// 		// json.NewEncoder(w).Encode(lFetchResponseRec)

// 		//marshall the Response to front
// 		data, lerr := json.Marshal(lFetchResponseRec)
// 		if lerr != nil {
// 			fmt.Fprintf(w, "Error taking data"+lerr.Error())
// 		} else {
// 			fmt.Fprintf(w, string(data))
// 			log.Println("data", string(data))
// 		}
// 		//lFetchResponseRec.fileNameArr is converting at marshal time i want exact content to the front

// 	}
// 	log.Println("SetToml Ended(-)")
// }
