package aol

import (
	"database/sql"
	"encoding/json"
	"fcs47pkg/ftdb"
	"fcs47pkg/helpers"
	"fcs47pkg/integration/aws/reposity"
	"fcs47pkg/spark/aol/base"
	"fcs47pkg/spark/aol/repository"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type programStruct struct {
	Id               int    `json:"id"`
	ProgramCode      string `json:"programCode"`
	ProgramName      string `json:"programName"`
	AboutTheProgram  string `json:"programDescription"`
	BusinessOwnwer   string `json:"businessOwnwer"`
	DeveloperName    string `json:"developerName"`
	DeveloperManager string `json:"developerManager"`
	User             string `json:"user"`
	FilePath         string `json:"filePath"`
	Server           string `json:"server"`
}

// type createdStruct struct {
// 	CreatedBy   string `json:"createdBy"`
// 	CreatedDate string `json:"createdDate"`
// 	UpdatedBy   string `json:"updatedBy"`
// 	UpdatedDate string `json:"updatedDate"`
// }

type fetchProgramResp struct {
	FetchProgram []programStruct `json:"fetchProgram"`
	Status       string          `json:"Status"`
	ErrMsg       string          `json:"errMsg"`
}

/*
Pupose: This API is used to insert the data of ProgramName, AboutTheProgram, BusinessOwnwer, DeveloperName, and DeveloperManager in the data table

MethodName:InsertDetails,GetNumber

Response:

        On Sucess
        =========
        {
	   {
	   "programStruct:{"programCode":"FCS_01","programName":"FCS","businessOwnwer":"kavya","developerName":"kavya","developerManager":"kavya","programDescription":"tested","user":"kavyadharshani.m@fcsonline.co.in"}
        "Status": "S",
        "ErrMsg": " ",
        }
        On Error
        ========
        {
	   "Status": "E",
        "ErrMsg": "Error",
        }

Author: KAVYA DHARSHANIInsertProgramNumber
Date: 07JUNE2023
*/

func InsertProgramNumber(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT,OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	log.Println("InsertProgramNumber(+)")

	if r.Method == "PUT" {

		// this variable is used for unmarshal of the value and then pass to front end it store the value in DB
		var lfetchRecp programStruct

		// this variable is used for get the details of createdby, created date, updatedby, updateddate for the database
		// var ldetails createdStruct

		// this variable is used for Error and then marshals the value and come back to front end
		var lresp fetchProgramResp

		lresp.Status = "S"
		body, err := ioutil.ReadAll(r.Body)
		//log.Println("body", string(body))

		if err != nil {
			log.Println("Error:", err)
			lresp.Status = "E"
			lresp.ErrMsg = "Error" + err.Error()
		} else {
			err = json.Unmarshal(body, &lfetchRecp)
			if err != nil {
				log.Println("Error:", err)
				lresp.Status = "E"
				lresp.ErrMsg = "error" + err.Error()
			} else {
				db, err := ftdb.LocalDbConnect(ftdb.MariaFTPRD)

				if err != nil {
					lresp.Status = "E"
					log.Println(err)
					lresp.ErrMsg = "Error: " + err.Error()
				} else {
					defer db.Close()

					// err := InsertDetails(db, lfetchRecp)

					err = repository.CreateNewRepo(lfetchRecp.ProgramCode, lfetchRecp.ProgramName, lfetchRecp.User, lfetchRecp.BusinessOwnwer, lfetchRecp.DeveloperName, lfetchRecp.DeveloperManager, lfetchRecp.AboutTheProgram)

					//log.Print("response", lresp)
					if err != nil {
						lresp.Status = "E"
						lresp.ErrMsg = "Error: " + err.Error()
						log.Println("error2", err)
					} else {
						lresp.FetchProgram, err = GetNumber(db)
						if err != nil {
							lresp.Status = "E"
							log.Println(err)
							lresp.ErrMsg = "Error: " + err.Error()
						}
					}
				}
			}
		}
		data, err := json.Marshal(lresp)
		if err != nil {
			fmt.Fprint(w, "Error taking data"+err.Error())
		} else {
			fmt.Fprint(w, string(data))
		}
		log.Println("InsertProgramNumber(-)")
	}
}

//===============================================================================================================================//

func InsertDetails(db *sql.DB, lfetchRecp programStruct) error {

	log.Println("InsertDetails1(+)")
	//log.Println("lfetchRecp", lfetchRecp)

	CoreString := `insert into programs(ProgramCode,ProgramName,ProgramDescription,BusinessOwner,DeveloperName,DevoloperManager,CreatedBy,CreatedDate,UpdatedBy,Updateddate)
						value(?,?,?,?,?,?,?,Now(),?,Now())`

	_, err := db.Exec(CoreString, lfetchRecp.ProgramCode, lfetchRecp.ProgramName, lfetchRecp.AboutTheProgram, lfetchRecp.BusinessOwnwer, lfetchRecp.DeveloperName, lfetchRecp.DeveloperManager, lfetchRecp.User, lfetchRecp.User)

	if err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Inserted successfully")
	}
	log.Println("InsertDetails(-)")
	return nil

}

// =============================================================== //

// func SelectDetails(db *sql.DB, lfetchRecp programStruct) (string, error) {

// 	// var lfetchRecp programStruct
// 	var check string
// 	log.Println("SelectDetails(+)")

// 	CoreString := `SELECT
// 	CASE WHEN COUNT(1) > 0 THEN 'Y' ELSE 'N' END
//      FROM
// 	 programs p
//      WHERE
// 	SUBSTRING_INDEX(p.ProgramName , '_', 2) =  SUBSTRING_INDEX(? , '_', 2) ; `

// 	rows, err := db.Query(CoreString, lfetchRecp.ProgramName)
// 	if err != nil {
// 		log.Println(err)
// 		return check, err
// 	} else {
// 		for rows.Next() {
// 			err := rows.Scan(&check)
// 			if err != nil {
// 				log.Println(err)
// 				return check, err
// 			}
// 		}
// 	}
// 	log.Println("check", check)
// 	log.Println("SelectDetails(-)")
// 	return check, err
// }

// ========================================================================================================== //

/*
Pupose: This API is used to get data for the database and it show the ProgramName, AboutTheProgram, BusinessOwnwer, DeveloperName, and DeveloperManager in the data table

MethodName:GetDetails

Response:

        On Sucess
        =========
	   {
	    "FetchProgram [{1  FCS tested kavya kavya kavya}]",
        "Status": "S",
        "ErrMsg": " ",
        }
        {
	   "Status": "E",
        "ErrMsg": ""Error",
        }


Author: KAVYA DHARSHANI
Date: 07JUNE2023
*/
func GetFetchProgramNumber(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "GET")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	log.Println("GetFetchProgramNumber(+)")

	// this variable is used for Error and tprogramCodehen marshals the value and comes back to the front end
	var lresp fetchProgramResp

	// this variable is used to find a error
	var err error

	lresp.Status = "S"
	db, err := ftdb.LocalDbConnect(ftdb.MariaFTPRD)
	if err != nil {
		log.Println(err.Error() + "(GD01)")
		lresp.ErrMsg = err.Error() + "(GD01)"
		lresp.Status = "E"
	} else {
		defer db.Close()
		lresp.FetchProgram, err = GetDetails(db)
		if err != nil {
			lresp.Status = "E"
			lresp.ErrMsg = err.Error() + "(GD02)"
			log.Println(err.Error() + "(GD02)")
		}

		FetchProgram, err := json.Marshal(lresp)
		if err != nil {
			fmt.Fprint(w, "Error taking FetchProgram"+err.Error())
		} else {
			fmt.Fprint(w, string(FetchProgram))
		}
		log.Println("GetFetchProgramNumber(-)")
	}
}

//===================================================================================================================//

func GetDetails(db *sql.DB) ([]programStruct, error) {
	log.Println("GetDetails+")

	// this variable is used to select the value of the data for the DataBase
	var lfetchRecp programStruct

	// this variable is used to append the value in array
	var FetchProgram []programStruct

	CoreString := `select Id,nvl(ProgramCode,''),nvl(ProgramName,''),nvl(ProgramDescription,''),nvl(BusinessOwner,''),nvl(DeveloperName,''),nvl(DevoloperManager,''),nvl(FilePath,''),nvl(Server,'') from programs where ProgramName!='unMappedPolicy'`
	rows, err := db.Query(CoreString)
	if err != nil {
		log.Println(err)
		return FetchProgram, err
	} else {
		for rows.Next() {
			err := rows.Scan(&lfetchRecp.Id, &lfetchRecp.ProgramCode, &lfetchRecp.ProgramName, &lfetchRecp.AboutTheProgram, &lfetchRecp.BusinessOwnwer, &lfetchRecp.DeveloperName, &lfetchRecp.DeveloperManager, &lfetchRecp.FilePath, &lfetchRecp.Server)
			if err != nil {
				log.Println("error1")
				return FetchProgram, err
			} else {
				FetchProgram = append(FetchProgram, lfetchRecp)
			}
		}
		log.Println("GetDetails(-)")
	}
	//log.Println("FetchProgram", FetchProgram)
	return FetchProgram, nil
}

/*
Pupose: This API is used to update the value in the data table

MethodName:UpdateMethods

Response:

        On Sucess
        =========
        {
	    FetchProgram [{1  FCS tested kavya kavya kavya} {2  BPO TESTED THE PROGRAM
          Raveeena Raveena Raveena} {3  IPO tested Karuniya Karuniya Karuniya} {4  EKYC tested priya priya priya}]
        "Status": "S",
        "ErrMsg": " ",
        }
        On Error
        ========
        {
	   "Status": "E",
        "ErrMsg": ""Error",
        }


Author: KAVYA DHARSHANI
Date: 07JUNE2023
*/

func UpdateFetchProgramNumber(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	log.Println("UpdateFetchProgramNumber(+)")
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()
	if r.Method == "PUT" {

		// this variable is used for unmarshal of the value and then pass to front end it store the value in DB
		var lfetchRecp programStruct

		// this variable is used for Error and then marshals the value and come back to front end
		var lresp fetchProgramResp

		lresp.Status = "S"
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			lresp.Status = "E"
			lresp.ErrMsg = "Error" + err.Error()
		} else {
			err = json.Unmarshal(body, &lfetchRecp)
			if err != nil {
				log.Println(err)
				lresp.Status = "E"
				lresp.ErrMsg = "error" + err.Error()
			} else {
				db, err := ftdb.LocalDbConnect(ftdb.MariaFTPRD)
				if err != nil {
					log.Println(err)
					lresp.Status = "E"
					lresp.ErrMsg = "error" + err.Error()
				} else {
					defer db.Close()

					lSess, lErr := base.CreateSession(lDebug)
					if lErr != nil {
						log.Println(err)
					}
					var lRepoName string
					if lfetchRecp.ProgramCode != "" {
						lRepoName = lfetchRecp.ProgramCode + "_" + lfetchRecp.ProgramName
					} else {
						lRepoName = lfetchRecp.ProgramName
					}
					if lfetchRecp.AboutTheProgram != "" {
						descinfo, lErr := reposity.UpdateDescription(lDebug, lSess, lRepoName, lfetchRecp.AboutTheProgram)
						if lErr != nil {
							lresp.ErrMsg = "error" + err.Error()
						}
						log.Println("descinfo", descinfo)
					}

					err := UpdateMethods(db, lfetchRecp)
					if err != nil {
						log.Println(err)
						lresp.Status = "E"
						lresp.ErrMsg = "error" + err.Error()
					}

				}
			}
		}

		data, err := json.Marshal(lresp)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Error taking data"+err.Error())
		} else {
			fmt.Fprint(w, string(data))
		}
	}
	log.Println("UpdateFetchProgramNumber(-)")
}

//============================================================================================================//

func UpdateMethods(db *sql.DB, lfetchRecp programStruct) error {
	log.Println("UpdateMethods(+)")
	// this variable is used to find a error
	var err error


	coreString := `update programs set ProgramCode=?, ProgramName=?, ProgramDescription=?, BusinessOwner=?,DeveloperName=?,DevoloperManager=?,FilePath=?,Server=?
	where Id=?`
	_, err = db.Exec(coreString, lfetchRecp.ProgramCode, lfetchRecp.ProgramName, lfetchRecp.AboutTheProgram, lfetchRecp.BusinessOwnwer, lfetchRecp.DeveloperName, lfetchRecp.DeveloperManager, lfetchRecp.FilePath, lfetchRecp.Server, lfetchRecp.Id)

	if err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Updated successfully")
	}
	log.Println("UpdateMethods(-)")
	return nil
}

// /*
// Pupose: This API is used to update the next Id number in the Vuetify page based on increment in the database

// MethodName:GetNumber

// Response:

//         On Sucess
//         =========
//         {
//          FetchProgram [{FCS_01}]
//         "Status": "S",
//         "ErrMsg": " ",
//         }
//         On Error
//         ========
//         {
// 	   "Status": "E",
//         "ErrMsg": ""Error",
//         }

// Author: KAVYA DHARSHANI
// Date: 07JUNE2023
// */

func GetNextNumber(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "GET")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	log.Println("GetNextNumber(+)")

	// this variable is used for Error and then marshals the value and comes back to the front end
	var lresp fetchProgramResp

	// this variable is used to find a error
	var err error

	lresp.Status = "S"
	db, err := ftdb.LocalDbConnect(ftdb.MariaFTPRD)
	if err != nil {
		log.Println(err.Error() + "(GD01)")
		lresp.ErrMsg = err.Error() + "(GD01)"
		lresp.Status = "E"
	} else {
		defer db.Close()
		lresp.FetchProgram, err = GetNumber(db)
		if err != nil {
			lresp.Status = "E"
			lresp.ErrMsg = err.Error() + "(GD02)"
			log.Println(err.Error() + "(GD02)")
		}

		FetchProgram, err := json.Marshal(lresp)
		if err != nil {
			fmt.Fprint(w, "Error taking FetchProgram"+err.Error())
		} else {
			fmt.Fprint(w, string(FetchProgram))
		}
		log.Println("GetNextNumber(-)")
	}
}

//==========================================================================================================//

func GetNumber(db *sql.DB) ([]programStruct, error) {
	log.Println("GetNumber+")

	// this variable is used to select the value of the data for the DataBase
	var lfetchRecp programStruct

	// this variable is used to append the value in array
	var FetchProgram []programStruct

	CoreString := `SELECT
	CASE
		WHEN MAX_number IS NULL THEN 'FCS_01'
		ELSE CONCAT('FCS_', LPAD(MAX_number + 1,CHAR_LENGTH(MAX_number + 1), '0'))
	END AS new_value
FROM (
	SELECT MAX(
		CAST(
			SUBSTRING_INDEX(ProgramCode, '_', -1)
		AS UNSIGNED)
	) AS MAX_number
	FROM Programs
) AS subquery`
	rows, err := db.Query(CoreString)
	if err != nil {
		log.Println(err)
		return FetchProgram, err
	} else {
		for rows.Next() {
			err := rows.Scan(&lfetchRecp.ProgramCode)
			if err != nil {
				log.Println("error1")
				return FetchProgram, err
			} else {
				FetchProgram = append(FetchProgram, lfetchRecp)
			}
		}

	}
	//log.Println("FetchProgram", FetchProgram)
	log.Println("GetNumber(-)")
	return FetchProgram, nil

}

// =============================================================================================================================== //
