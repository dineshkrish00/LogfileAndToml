package aol

import (
	"encoding/json"
	"fcs47pkg/common"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type FetchDataStruct struct {
	Path string `json:"path"`
}

type fileNamestr struct {
	FileName string      `json:"fileName"`
	Content  interface{} `json:"content"`
}

type FetchResponseStruct struct {
	Status      string        `json:"status"`
	ErrMsg      string        `json:"errMsg"`
	FilePath    string        `json:"filePath"`
	FileNameArr []fileNamestr `json:"fileNameArr"`
}

func FetchingLogFile(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchingLogFile Started(+)")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT")
	(w).Header().Set("Access-Control-Allow-Headers", " Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "PUT" {
		var lFetchDataRec FetchDataStruct
		var lFetchResponseRec FetchResponseStruct

		lFetchResponseRec.Status = "S"

		lbody, lerr := ioutil.ReadAll(r.Body)
		if lerr != nil {
			lFetchResponseRec.Status = "E"
			lFetchResponseRec.ErrMsg = "Error in (flf01): " + lerr.Error()
		} else {
			//Unmarshall the body and save into FetchDataStruct Struct
			lerr = json.Unmarshal(lbody, &lFetchDataRec)
			if lerr != nil {
				lFetchResponseRec.Status = "E"
				lFetchResponseRec.ErrMsg = "Error in (flf02): " + lerr.Error()
			} else {
				//assigning filepath to struct variable
				lFetchResponseRec.FilePath = lFetchDataRec.Path
				//Calling CheckingFileName method

				// path := fmt.Sprintf("%v", configFile.(map[string]interface{})["Path"])
				lFetchResponseRec.FileNameArr, lerr = CheckingFileName(lFetchDataRec)
				if lerr != nil {
					lFetchResponseRec.Status = "E"
					lFetchResponseRec.ErrMsg = "Error in (flf03): " + lerr.Error()
				}
			}
		}
		log.Println("lFetchResponseRec", lFetchResponseRec)

		//marshall the Response to front
		data, lerr := json.Marshal(lFetchResponseRec)
		if lerr != nil {
			fmt.Fprintf(w, "Error taking data"+lerr.Error())
		} else {
			fmt.Fprintf(w, string(data))
			log.Println("data", string(data))
		}
	}
	log.Println("FetchingLogFile Ended(-)")
}

func CheckingFileName(pInputData FetchDataStruct) ([]fileNamestr, error) {
	log.Println("CheckingFileName (+)")

	var lfileNamesRec fileNamestr
	var lfileNamesArr []fileNamestr

	//Read the User Given File Path
	lfiles, lErr := ioutil.ReadDir(pInputData.Path)
	if lErr != nil {
		log.Println("Error in ReadDir (CFN01): ", lErr.Error())
		return lfileNamesArr, lErr
	} else {
		//Iterating and Sepearting the given (path ,from date,todate) based on file names and save into lFetchResponseRec.FileNameArr
		for _, lfile := range lfiles {
			//Checking the Directory
			if !lfile.IsDir() {
				lfileNamesRec.FileName = lfile.Name()
				if filepath.Ext(lfile.Name()) == ".toml" && len(lfileNamesRec.FileName) >= 3 {
					//Read the toml file and return the values
					configFile := common.ReadTomlConfig(filepath.Join(pInputData.Path, lfile.Name()))

					lfileNamesRec.Content = configFile

					lfileNamesArr = append(lfileNamesArr, lfileNamesRec)
				}
			} else {
				return lfileNamesArr, fmt.Errorf("Error in (CFN05) Directory Path")
			}
		}

	}

	log.Println("CheckingFileName (-)")
	return lfileNamesArr, nil
}

// package aol

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"time"
// )

// type FetchDataStruct struct {
// 	Path     string `json:"path"`
// 	FromDate string `json:"fromDate"`
// 	ToDate   string `json:"toDate"`
// }

// type fileNamestr struct {
// 	FileName string `json:"fileName"`
// }

// type FetchResponseStruct struct {
// 	Status      string        `json:"status"`
// 	ErrMsg      string        `json:"errMsg"`
// 	FilePath    string        `json:"filePath"`
// 	FileNameArr []fileNamestr `json:"fileNameArr"`
// }

// /*
// Pupose: This API is getting request as path,fromDate,todate based on value it can taken the log file from the path and send it as response

// MethodName:FetchingLogFile,CheckingFileName

// Response:

//         On Sucess
//         =========
// 	  {"status":"S","errMsg":"","filePath":"/home/user/Documents/Dinesh/CODE/Spark/01-04-24(program master)/FCS_47_ListenLocalPKG/log","fileNameArr":[
// 		{"fileName":"logfile02042024.14.08.48.578283309.txt"},
// 	  	{"fileName":"logfile02042024.15.03.34.897596350.txt"}
//         "Status": "S",
//         "ErrMsg": " ",
//         }
//         On Error
//         ========
//        {"status":"E","errMsg":"Error: open /home/user/Documents/ODE/Spark/01-04-24(program master)/FCS_47_ListenLocalPKG/log: no such file or directory",
// 		"filePath":"/home/user/Documents/ODE/Spark/01-04-24(program master)/FCS_47_ListenLocalPKG/log",
// 		"fileNameArr":null}

// Author: DINESH KUMAR K U
// Date: 04APRIL2024
// */

// /*	Method 1: FetchingLogFile
// 	Purpose : To unmarshall the request and passing to the Checking FileName method the struct after getting response from method
// 			the marshall the response
// 	Author: DINESH KUMAR K U
// 	Date: 04APRIL2024
// */

// func FetchingLogFile(w http.ResponseWriter, r *http.Request) {
// 	log.Println("FetchingLogFile Started(+)")
// 	(w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(w).Header().Set("Access-Control-Allow-Credentials", "true")
// 	(w).Header().Set("Access-Control-Allow-Methods", "PUT")
// 	(w).Header().Set("Access-Control-Allow-Headers", " Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 	if r.Method == "PUT" {
// 		var lFetchDataRec FetchDataStruct
// 		var lFetchResponseRec FetchResponseStruct

// 		lFetchResponseRec.Status = "S"

// 		lbody, lerr := ioutil.ReadAll(r.Body)
// 		if lerr != nil {
// 			lFetchResponseRec.Status = "E"
// 			lFetchResponseRec.ErrMsg = "Error in (flf01): " + lerr.Error()
// 		} else {
// 			//Unmarshall the body and save into FetchDataStruct Struct
// 			lerr = json.Unmarshal(lbody, &lFetchDataRec)
// 			if lerr != nil {
// 				lFetchResponseRec.Status = "E"
// 				lFetchResponseRec.ErrMsg = "Error in (flf02): " + lerr.Error()
// 			} else {
// 				//assigning filepath to struct variable
// 				lFetchResponseRec.FilePath = lFetchDataRec.Path
// 				//Calling CheckingFileName method
// 				lFetchResponseRec.FileNameArr, lerr = CheckingFileName(lFetchDataRec)
// 				if lerr != nil {
// 					lFetchResponseRec.Status = "E"
// 					lFetchResponseRec.ErrMsg = "Error in (flf03): " + lerr.Error()
// 				}
// 			}
// 		}
// 		//marshall the Response to front
// 		data, lerr := json.Marshal(lFetchResponseRec)
// 		if lerr != nil {
// 			fmt.Fprintf(w, "Error taking data"+lerr.Error())
// 		} else {
// 			fmt.Fprintf(w, string(data))
// 		}
// 	}
// 	log.Println("FetchingLogFile Ended(-)")
// }

// /*
// Method 2: CheckingFileName
// Purpose : In this method read the path, convert the date and filtering the file name then append after return the struct as response
// Author: DINESH KUMAR K U
// Date: 04APRIL2024
// */
// func CheckingFileName(pInputData FetchDataStruct) ([]fileNamestr, error) {
// 	log.Println("CheckingFileName (+)")

// 	var lfileNamesRec fileNamestr
// 	var lfileNamesArr []fileNamestr

// 	//Read the User Given File Path
// 	lfiles, lerr := ioutil.ReadDir(pInputData.Path)
// 	if lerr != nil {
// 		log.Println("ERROR: (0X100) ", lerr)
// 		return lfileNamesArr, lerr
// 	} else {
// 		//time.parse a time string in a specific format and return a Time--from Date
// 		lfromTime, lerr := time.Parse("02012006", pInputData.FromDate)
// 		if lerr != nil {
// 			log.Println("fromTime Conversion (0X101): ", lerr)
// 			return lfileNamesArr, lerr
// 		} else {
// 			//time.parse a time string in a specific format and return a Time--to Date
// 			ltoTime, lerr := time.Parse("02012006", pInputData.ToDate)
// 			if lerr != nil {
// 				log.Println("Error in (0X102):", lerr)
// 				return lfileNamesArr, lerr
// 			} else {

// 				//Iterating and Sepearting the given (path ,from date,todate) based on file names and save into lFetchResponseRec.FileNameArr
// 				for _, lfile := range lfiles {
// 					//Checking the Directory
// 					if !lfile.IsDir() {
// 						lfileNamesRec.FileName = lfile.Name()
// 						if len(lfileNamesRec.FileName) >= 15 {
// 							lfileDate, lerr := time.Parse("02012006", lfileNamesRec.FileName[7:15])
// 							if lerr != nil {
// 								log.Println("Error in (0X103):", lerr)
// 								continue
// 							}
// 							//Check the Condition using the values for filter the appending values
// 							if (lfromTime.Before(lfileDate) || lfromTime.Equal(lfileDate)) && (lfileDate.Before(ltoTime) || lfileDate.Equal(ltoTime)) {
// 								lfileNamesArr = append(lfileNamesArr, lfileNamesRec)
// 							}
// 						}
// 					} else {

// 						return lfileNamesArr, fmt.Errorf("Error in (CFN05) Given Directory Path")
// 					}
// 				}
// 			}
// 		}
// 	}
// 	log.Println("CheckingFileName (-)")
// 	return lfileNamesArr, nil
// }
