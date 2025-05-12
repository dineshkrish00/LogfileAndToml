package aol

import (
	"database/sql"
	"encoding/json"
	"fcs47pkg/ftdb"
	"fmt"
	"log"
	"net/http"
)

type CodeDataStruct struct {
	ProductCode  string `json:"productCode"`
	Descripition string `json:"descripition"`
}

type CodeResponseStruct struct {
	ProductCodeArr []CodeDataStruct `json:"productCodeArr"`
	Status         string           `json:"status"`
	ErrMsg         string           `json:"errMsg"`
}

/*
Pupose: This API is used to Fetch the ProductCode and  Descripition(FilePath,Server,ProgramCode,ProgramName) from program and send it at response

MethodName:FetchProgramCode,SelectProgramCode

Response:

        On Sucess
        =========
	   {"productCodeArr":[{"productCode":"FCS_156","descripition":"/home/user/Documents/Dinesh/CODE/Flow/21-03-24 logfile download/FCS_94_FlowAPI/log,192.168.2.5,FCS_156_Testing"},
	   {"productCode":"FCS_151","descripition":"/home/user/Documents/Dinesh/CODE/Spark/01-04-24(program master)/FCS_47_ListenLocalPKG/log,192.168.150.12,FCS_151_CallListingingReport"}}
        "Status": "S",
        "ErrMsg": " ",
        }
        On Error
        ========
        {
	   "Status": "E",
        "ErrMsg": "Error",
        }

Author: DINESH KUMAR K U
Date: 04APRIL2024
*/

/*	Method 1: FetchProgramCode
	Purpose : To Establish the MariaFtPrd db connection and Passing another method to do task then send response to front
	Author: DINESH KUMAR K U
	Date: 04APRIL2024*/

func FetchProgramCode(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchProgramCode(+)")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", " Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(200)
	if r.Method == "GET" {
		var lFetchResponseRec CodeResponseStruct
		lFetchResponseRec.Status = "S"
		db, err := ftdb.LocalDbConnect(ftdb.MariaFTPRD)
		if err != nil {
			log.Println(err.Error() + "(FPD01)")
			lFetchResponseRec.ErrMsg = err.Error() + "(FPD01)"
			lFetchResponseRec.Status = "E"
		} else {
			defer db.Close()
			lFetchResponseRec.ProductCodeArr, err = SelectProgramCode(db)
			if err != nil {
				lFetchResponseRec.Status = "E"
				lFetchResponseRec.ErrMsg = err.Error() + "(FPD02)"
			}
		}
		data, err := json.Marshal(lFetchResponseRec)
		if err != nil {
			fmt.Fprintf(w, "Error taking data"+err.Error())
		} else {

			fmt.Fprintf(w, string(data))
		}
	}
	log.Println("FetchProgramCode(-)")
}

/*	Method 2: SelectProgramCode
	Purpose : To Fetch the Data from program table and return the value , error
	Author: DINESH KUMAR K U
	Date: 04APRIL2024*/

func SelectProgramCode(db *sql.DB) ([]CodeDataStruct, error) {
	log.Println("SelectProgramCode(+)")
	var lFetchDataArr []CodeDataStruct
	var lCodeDataRec CodeDataStruct

	lCoreString := `select concat(nvl(ProgramCode,''),'_',nvl(ProgramName,'')), 
					concat(nvl(TomlPath, ''), ',', nvl(Server, ''),',',concat(nvl(ProgramCode,''),'_',nvl(ProgramName,''))) as FileServer
					from programs`

	lrows, lerr := db.Query(lCoreString)
	if lerr != nil {
		log.Println("Error: (PC001)", lerr)
		return nil, lerr
	} else {
		defer lrows.Close()
		for lrows.Next() {
			lerr := lrows.Scan(&lCodeDataRec.ProductCode, &lCodeDataRec.Descripition)
			if lerr != nil {
				log.Println("Error: (PC002)", lerr)
				return nil, lerr
			} else {
				if lCodeDataRec.ProductCode != "" {
					lFetchDataArr = append(lFetchDataArr, lCodeDataRec)
				}
			}
		}
	}
	log.Println("SelectProgramCode(-)")
	return lFetchDataArr, nil
}
