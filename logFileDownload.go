package aol

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fcs47pkg/common"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileObj struct {
	FileName string `json:"fileName"`
}

type ZipReqStruct struct {
	FileName []FileObj `json:"fileName"`
	Path     string    `json:"path"`
}

type ZipResponseStruct struct {
	Status  string `json:"status"`
	ErrMsg  string `json:"errMsg"`
	ZipFile string `json:"zipFile"`
}

/*
Pupose: This API is getting request as Filename,Path from front and Creating the zipfile then filter the user given records
		and save into one zip folder then send as response type of "application/zip" for blob object

MethodName:LogFileDownload,zipFiles

Response:

        On Sucess
        =========
		ZipFile: Blob Object
        "Status": "S",
        "ErrMsg": " ",
        }
        On Error
        ========
		ZipFile: Error
        "Status": "E",
        "ErrMsg": "Error in Converting zip file",

Author: DINESH KUMAR K U
Date: 04APRIL2024
*/

/*	Method 1: LogFileDownload
	Purpose : In this Method  im unmarshalling the request and save in lZipReqRec struct then call the method get the zip file as response
			once file is available , marshall and attached the lzip file folder as send
	Author: DINESH KUMAR K U
	Date: 04APRIL2024
*/

func LogFileDownload(w http.ResponseWriter, r *http.Request) {
	log.Println("LogFileDownload Started(+)")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT")
	(w).Header().Set("Access-Control-Allow-Headers", " Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "PUT" {
		var lZipResponeRec ZipResponseStruct
		var lZipReqRec ZipReqStruct
		lZipResponeRec.Status = "S"
		//Read toml from zip_file path
		lconfigFile := common.ReadTomlConfig("./toml/amconfig.toml")
		lZipFilePath := fmt.Sprintf("%v", lconfigFile.(map[string]interface{})["ZipFilePath"])

		lbody, lerr := ioutil.ReadAll(r.Body)
		if lerr != nil {
			log.Println("Error in reading time", lerr)
			lZipResponeRec.Status = "E"
			lZipResponeRec.ErrMsg = "Error in (LFD01): " + lerr.Error()
		} else {
			lerr = json.Unmarshal(lbody, &lZipReqRec)
			if lerr != nil {
				lZipResponeRec.Status = "E"
				lZipResponeRec.ErrMsg = "Error in (LFD02): " + lerr.Error()
				log.Println("unmarshall time", lerr)
			} else {
				lerr := zipFiles(lZipFilePath, lZipReqRec)
				if lerr != nil {
					lZipResponeRec.ErrMsg = lerr.Error() + "(LFD03)"
					lZipResponeRec.Status = "E"
					http.Error(w, "Error creating zip file", http.StatusInternalServerError)
				} else {
					//lZipFilePath assingning to the struct
					lZipResponeRec.ZipFile = lZipFilePath
				}
			}
			//Sending the Response to the front as zip file (i.e)blob file
			w.Header().Set("Content-Disposition", "attachment; filename=logFile.zip")
			w.Header().Set("Content-Type", "application/zip")
			http.ServeFile(w, r, lZipFilePath)
			data, err := json.Marshal(lZipResponeRec)
			if err != nil {
				fmt.Fprintf(w, "Error taking data"+err.Error())
			} else {
				fmt.Fprintf(w, string(data))
			}
			//removing after sending lzip file
			err = os.Remove(lZipFilePath)
			if err != nil {
				log.Println("Error deleting zip file from local path:", err)
			}
		}
	}
	log.Println("LogFileDownload Ended (-)")
}

/*	Method 2: ZipFiles
	Purpose : In this Method  im Creating zip path and read the local log file & compare with user give file header name
			after filter read & write the content in respesctive file then finally conver the zip file
	Author: DINESH KUMAR K U
	Date: 04APRIL2024
*/

func zipFiles(pZipFilePath string, pInputData ZipReqStruct) error {
	log.Println("zipFiles (+)")

	lZipFile, lerr := os.Create(pZipFilePath)
	if lerr != nil {
		log.Println("ERROR: (0Zf01) ", lerr)
		return lerr
	} else {
		defer lZipFile.Close()
		//zip.NewWriter is a function provided by the archive/zip package,
		lZipWriter := zip.NewWriter(lZipFile)
		defer lZipWriter.Close()
		//Check the file info to iterate through all files and directories under a specific directory
		lerr = filepath.Walk(pInputData.Path, func(filePath string, info os.FileInfo, lerr error) error {
			if lerr != nil {
				log.Println("ERROR: (0Zf02) ", lerr)
				return lerr
			} else {
				//Read the File InfoHeader
				lHeader, lerr := zip.FileInfoHeader(info)
				if lerr != nil {
					log.Println("ERROR: (0Zf03) ", lerr)
					return lerr
				} else {
					//Here file extracting with the name
					lHeader.Name = strings.TrimPrefix(lHeader.Name, "/log")

					//Iterating the User Selected Name & Compare with lHeader.Name
					for _, lfileName := range pInputData.FileName {
						if lHeader.Name == lfileName.FileName {
							lHeader.Method = zip.Deflate
							if info.IsDir() {
								lHeader.Name += "/"
								_, lerr = lZipWriter.CreateHeader(lHeader)
								if lerr != nil {
									log.Println("ERROR: (0Zf04) ", lerr)
									return lerr
								}
							}
							lfile, lerr := os.Open(filePath)
							if lerr != nil {
								log.Println("ERROR: (0Zf05) ", lerr)
								return lerr
							} else {
								defer lfile.Close()
								lreader := bufio.NewReader(lfile)
								lwriter, lerr := lZipWriter.CreateHeader(lHeader)
								if lerr != nil {
									log.Println("ERROR: (0Zf06) ", lerr)
									return lerr
								} else {
									_, lerr = io.Copy(lwriter, lreader)
									if lerr != nil {
										log.Println("ERROR: (0Zf07) ", lerr)
										return lerr
									}
								}
							}
						}
					}
				}
			}
			return lerr
		})
		log.Println("zipFiles (-)")
		return lerr
	}
}
