<template>
  <div>
    <v-container class="my-5">
      <v-layout row wrap class="mb-1">
        <v-flex class="d-flex justify-left" lg2>
          <v-icon left color="blue darken-4" medium
            >mdi-file-cog</v-icon
          >
          <b>Toml File Config</b>
        </v-flex>
      </v-layout>
    </v-container>
    <Datepicker @getData="getData" />
    <v-container class="my-2">
      <v-card elevation="3" class="my-2">
        <v-layout class="mx-5 pt-5 pt-2">
          <div class="font-weight-bold mr-4">Program Name:</div>
          <v-flex>{{ programCode }}</v-flex>
          <v-flex class="text-right">
            <v-btn small class="mr-5 grey white--text" v-if="this.listedFileName.length>0" rounded @click="Filedialog=true">
              <v-icon class="white--text mr-3" size="15">mdi-plus</v-icon>Add File</v-btn>
          </v-flex>
        </v-layout>
        <v-text-field
          class="pa-5"
          v-model="search"
          append-icon="mdi-magnify"
          label="Search"
          single-line
          hide-details
          clearable
        ></v-text-field>
        <v-layout>
          <v-flex class="d-flex justify-end mr-2 mb-3">
            <v-btn
              v-if="selected.length > 0"
              color="blue darken-2"
              text
              small
              @click="downloadProgress"
            >
              <v-icon left color="blue darken-2" medium
                >mdi mdi-arrow-collapse-down</v-icon
              >Download
            </v-btn>
          </v-flex>
      
        </v-layout>
        <v-data-table
          :headers="headers"
          :items="listedFileName"
          :search="search"
          :loading="loading"
          class="elevation-0"
          height="400"
          fixed-header
          dense
          :items-per-page="-1"
          :single-select="true"
           @dblclick:row="(event, { item }) => openDialog(event, item)"
        >
          <template v-slot:item.Select="{ item }">
            <v-checkbox
              v-model="selected"
              :value="item"
              dense
              :key="item.fileName"
            ></v-checkbox>
          </template>
          <!-- <template v-slot:item.actions="{ item }">
    <v-btn small @click="performAction(item)">Action</v-btn>
  </template> -->
        </v-data-table>
      </v-card>
      <v-overlay :value="isLoading">
        <v-progress-circular
          indeterminate
          size="64"
          color="primary"
        ></v-progress-circular>
      </v-overlay>
      <Popup :dialog="dialog" :fileContent="this.fileContent"  :File="this.File" @save-file-data="SetFileData" @closetab="dialog = false"/>
      <AddFile  :Filedialog="Filedialog" @Addclosetab="Filedialog=false" @Add-file-data="AddFileData"/>
    </v-container>
  </div>
</template>

<script>
import Datepicker from "./logfileDate.vue";
import EventService from "../../../services/EventService";
import Popup from "./Popup.vue"
import AddFile from "./AddToml.vue"

export default {
  components: {
    Datepicker,
    Popup,
    AddFile
  },
  data() {
    return {
      showContainer: false,
      Filedialog:false,
      loading: false,
      selectedItems: [],
      listedFileName: [],
      search: "",
      headers: [
        {
          text: "Select the file",
          sortable: false,
          value: "Select",
          class: "blue lighten-5",
        },
        {
          text: "Toml File Names",
          sortable: true,
          value: 'fileName',
          class: "blue lighten-5",
          align: "left",
        },
  //       {
  //   text: "Actions",
  //   value: "actions",
  //   sortable: false,
  //   class: "blue lighten-5",
  // },
      ],
      isLoading: false,
      selected: [],
      checkBox: false,
      selectedItemsArray: [],
      Data: {
        fileName: [],
        path: "",
        server: "",
      },
      programCode: "",
      downloadServer:"",
      Status: "",
      dialog:false,
      fileContent:'',
      FileData:{
        fileContent: '',
        fileName:'',
        filePath:'',
        server:'',
      },
      Reload:{
        path:''
      },
      File:'',
    };
  },
  methods: {
    SetFileData(fileContent,fileName){
      this.FileData.fileContent=fileContent
      this.FileData.fileName=fileName
      EventService.TomlValues(this.FileData)
        .then((response) => {
          if (response.data.status == "S") {
            this.MessageBar("S", "Data Added Succesfully");  
            this.ReloadDataTable()
          } else {
            this.MessageBar("E", response.data.errMsg);       
          }
        })
        .catch((error) => {
          this.MessageBar("E", error);
        });
    },

    AddFileData(fileContent,fileName){
      this.SetFileData(fileContent,fileName)
      // console.log("AddFileData this.FileData",this.FileData)
    },

    ReloadDataTable(){
      console.log("ReloadDataTable (+)");
      if (this.FileData.server == "192.168.2.5") {
        /*--DAC Development Api Client--*/
        this.Status = "DAC";
      } else if (this.FileData.server == "192.168.150.12") {
        /*--PAC Production Api Client-- */
        this.Status = "PAC";
      } else if (this.FileData.server == "192.168.150.21") {
        /*--BAC Base Api Client-- */
        this.Status = "BAC";
      }else{
        this.MessageBar('E', 'Error in Server Address '+this.FileData.server)
        this.isLoading = false;
        this.loading=false
      }
      this.listedFileName=[]
      this.Reload.path=this.FileData.filePath
      if(this.Status!=''){     
      EventService.FetchingLogFile(this.Reload, this.Status)
        .then((response) => {
          if (response.data.status == "S") {
            //new condition
            if (response.data.fileNameArr != null) {
              this.loading=false
              this.listedFileName = response.data.fileNameArr.map(file => 
              ({ fileName: file.fileName,content: file.content }));
              console.log("this.listedFileName",this.listedFileName)
            }
            //here push path
            this.showContainer = true;
            this.isLoading = false;
            this.loading = false;
          } else {
            this.MessageBar("E", response.data.errMsg);
            this.isLoading = false;
            this.loading = false;
          }
        })
        .catch((error) => {
          this.MessageBar("E", error);
          console.log(error);
          this.isLoading = false;
          this.loading = false;
        });
      }
      this.Status=''
      console.log("ReloadDataTable (-)")
    },
    openDialog(event, item) {
      console.log("event",event)
      this.dialog = true;
      this.fileContent = this.listedFileName.find((file) => file.fileName === item.fileName).content;
      let contentString = '';
      for (const key in this.fileContent) {
        contentString += `${key} = "${this.fileContent[key]}"\n`;
      }
        this.fileContent = contentString;
        this.File=item.fileName
       console.log("this.fileContent",this.fileContent)
    },


    //method 1: Download btn processing
    downloadProgress() {
      this.selected.forEach((item) => {
        this.selectedItemsArray.push(item.fileName);
      });
      this.Data.fileName = this.selected;
      //calling csv file writing method
      this.downloadZip();
    },

    //method 2: Zip File downloading method
    downloadZip() {
      this.isLoading = true;
      this.showContainer = false;
      //Calling Zip file Download Event Serivce
      if(this.Data!='' && this.downloadServer!=''){
      EventService.LogDownload(this.Data, this.downloadServer)
        .then((response) => {
          if (response.statusText == "OK") {
            console.log("response.data", response.data);
            // this.MessageBar('S', 'File downloaded successfully')
            const blob = new Blob([response.data], { type: "application/zip" });
            const fileUrl = URL.createObjectURL(blob);
            // Create a link element and set its attributes
            const link = document.createElement("a");
            link.setAttribute("href", fileUrl);
            link.setAttribute("download", "TomlFiles.zip"); // Set the filename from response
            // Append the link to the document body and trigger click event
            link.style.visibility = "hidden";
            document.body.appendChild(link);
            link.click();
            // Cleanup
            document.body.removeChild(link);
            URL.revokeObjectURL(fileUrl);
            this.isLoading = false;
          } else {
            this.MessageBar("E", "Fail to downloaded");
            this.isLoading = false;
          }
        })
        .catch((error) => {
          this.isLoading = false;
          this.MessageBar("E", error);
          console.log(error);
        });
    }
      this.isLoading = false;
      this.selected = [];
      this.downloadServer=''
    },
    //method 3:Getting File Name for Data table
    getData(lDate) {
      this.loading = true;
      this.selectedItemsArray = [];
      this.isLoading = true;
      this.Data.path = lDate.path;
      this.programCode = lDate.programCode;
      this.Data.server = lDate.server;
      this.FileData.filePath=lDate.path;
      this.FileData.server=lDate.server;
      //Server Checking Condition
      if (this.Data.server == "192.168.2.5") {
        /*--DAC Development Api Client--*/
        this.Status = "DAC";
      } else if (this.Data.server == "192.168.150.12") {
        /*--PAC Production Api Client-- */
        this.Status = "PAC";
      } else if (this.Data.server == "192.168.150.21") {
        /*--BAC Base Api Client-- */
        this.Status = "BAC";
      }else{
        this.MessageBar('E', 'Error in Server Address '+this.Data.server)
        this.isLoading = false;
        this.loading=false
      }

      //Mt the array
      this.listedFileName = [];
      this.selected = [];
      //Calling Event service for user selected log file fetching
      if(this.Status!='' && lDate!=''){     
      EventService.FetchingLogFile(lDate, this.Status)
        .then((response) => {
          console.log("response.data",response.data);
   
          console.log("response.data.fileNameArr", response.data.fileNameArr); // Should print the array of objects
   

          if (response.data.status == "S") {
            //new condition
            if (response.data.fileNameArr != null) {
              // this.loading=false
              this.listedFileName = response.data.fileNameArr.map(file => ({ fileName: file.fileName,content: file.content }));
              console.log("this.listedFileName",this.listedFileName)
          
            }
            //here push path
            this.showContainer = true;
            this.isLoading = false;
            this.loading = false;
          } else {
            this.MessageBar("E", response.data.errMsg);
            this.isLoading = false;
            this.loading = false;
          }
        })
        .catch((error) => {
          this.MessageBar("E", error);
          console.log(error);
          this.isLoading = false;
          this.loading = false;
        });
      }
      this.downloadServer=this.Status
      this.Status=''
    },
  },
};
</script> 

<!-- <template>
  <div>
    <v-container class="my-5">
      <v-layout row wrap class="mb-5">
        <v-flex class="d-flex justify-left" lg2>
          <v-icon left color="blue darken-4" medium
            >mdi mdi-arrow-collapse-down</v-icon
          >
          <b>Log File Download</b>
        </v-flex>
      </v-layout>
    </v-container>
    <Datepicker @getData="getData" />
    <v-container class="my-2">
      <v-card elevation="3" class="my-2">
        <v-layout class="mx-5 pt-5 pt-2">
          <div class="font-weight-bold mr-4">Program Code:</div>
          <v-flex>{{ programCode }}</v-flex>
          <div class="font-weight-bold mr-4">From Date:</div>
          <v-flex>{{ fromDate }}</v-flex>
          <div class="font-weight-bold mr-4">To Date:</div>
          <v-flex>{{ toDate }}</v-flex>
        </v-layout>
        <v-text-field
          class="pa-5"
          v-model="search"
          append-icon="mdi-magnify"
          label="Search"
          single-line
          hide-details
          clearable
        ></v-text-field>
        <v-layout>
          <v-flex class="d-flex justify-end mr-2 mb-3">
            <v-btn
              v-if="selected.length > 0"
              color="blue darken-2"
              text
              small
              @click="downloadProgress"
            >
              <v-icon left color="blue darken-2" medium
                >mdi mdi-arrow-collapse-down</v-icon
              >Download
            </v-btn>
          </v-flex>
        </v-layout>
        <v-data-table
          :headers="headers"
          :items="listedFileName"
          :search="search"
          :loading="loading"
          class="elevation-0"
          height="400"
          fixed-header
          :items-per-page="-1"
          :single-select="true"
        >
          <template v-slot:item.Select="{ item }">
            <v-checkbox
              v-model="selected"
              :value="item"
              :key="item.fileName"
            ></v-checkbox>
          </template>
        </v-data-table>
      </v-card>
      <v-overlay :value="isLoading">
        <v-progress-circular
          indeterminate
          size="64"
          color="primary"
        ></v-progress-circular>
      </v-overlay>
    </v-container>
  </div>
</template>

<script>
import Datepicker from "./logfileDate.vue";
import EventService from "../../../services/EventService";

export default {
  components: {
    Datepicker,
  },
  data() {
    return {
      showContainer: false,
      loading: false,
      selectedItems: [],
      listedFileName: [],
      search: "",
      headers: [
        {
          text: "Select the file",
          sortable: false,
          value: "Select",
          class: "blue lighten-5",
        },
        {
          text: "Available Log Files",
          sortable: true,
          value: "fileName",
          class: "blue lighten-5",
          align: "left",
        },
      ],
      isLoading: false,
      selected: [],
      checkBox: false,
      selectedItemsArray: [],
      Data: {
        fileName: [],
        path: "",
        server: "",
      },
      programCode: "",
      downloadServer:"",
      fromDate: "",
      toDate: "",
      Status: "",
      // ResetForm:false,
    };
  },
  methods: {
    //method 1: Download btn processing
    downloadProgress() {
      this.selected.forEach((item) => {
        this.selectedItemsArray.push(item.fileName);
      });
      this.Data.fileName = this.selected;
      //calling csv file writing method
      this.downloadZip();
    },

    //method 2: Zip File downloading method
    downloadZip() {
      this.isLoading = true;
      this.showContainer = false;
      //Calling Zip file Download Event Serivce
      if(this.Data!='' && this.downloadServer!=''){
      EventService.LogDownload(this.Data, this.downloadServer)
        .then((response) => {
          if (response.statusText == "OK") {
            console.log("response.data", response.data);
            // this.MessageBar('S', 'File downloaded successfully')
            const blob = new Blob([response.data], { type: "application/zip" });
            const fileUrl = URL.createObjectURL(blob);
            // Create a link element and set its attributes
            const link = document.createElement("a");
            link.setAttribute("href", fileUrl);
            link.setAttribute("download", "logfiles.zip"); // Set the filename from response
            // Append the link to the document body and trigger click event
            link.style.visibility = "hidden";
            document.body.appendChild(link);
            link.click();
            // Cleanup
            document.body.removeChild(link);
            URL.revokeObjectURL(fileUrl);
            this.isLoading = false;
          } else {
            this.MessageBar("E", "Fail to downloaded");
            this.isLoading = false;
          }
        })
        .catch((error) => {
          this.isLoading = false;
          this.MessageBar("E", error);
          console.log(error);
        });
    }
      // Field MT
      // this.ResetForm=true
      // this.programCode = ''
      // this.toDate = ''
      // this.fromDate = ''
      // this.listedFileName = []
      this.isLoading = false;
      this.selected = [];
      this.downloadServer=''
    },
    //method 3:Getting File Name for Data table
    getData(lDate) {
      this.loading = true;
      this.selectedItemsArray = [];
      this.isLoading = true;
      //From Date Conversion
      this.fromDate = lDate.fromDate;
      const FromDateParts = lDate.fromDate.split("-");
      lDate.fromDate = FromDateParts.reverse().join("");
      //To Date Conversion
      this.toDate = lDate.toDate;
      const ToDateParts = lDate.toDate.split("-");
      lDate.toDate = ToDateParts.reverse().join("");

      this.Data.path = lDate.path;
      this.programCode = lDate.programCode;
      this.Data.server = lDate.server;

      //Server Checking Condition
      if (this.Data.server == "192.168.2.5") {
        /*--DAC Development Api Client--*/
        this.Status = "DAC";
      } else if (this.Data.server == "192.168.150.12") {
        /*--PAC Production Api Client-- */
        this.Status = "PAC";
      } else if (this.Data.server == "192.168.150.21") {
        /*--BAC Base Api Client-- */
        this.Status = "BAC";
      }else{
        this.MessageBar('E', 'Error in Server Address '+this.Data.server)
        this.isLoading = false;
        this.loading=false
         lDate.fromDate =this.fromDate
       lDate.toDate =this.toDate
      }

      //Mt the array
      this.listedFileName = [];
      this.selected = [];
      //Calling Event service for user selected log file fetching
      if(this.Status!='' && lDate!=''){     
      EventService.FetchingLogFile(lDate, this.Status)
        .then((response) => {
          if (response.data.status == "S") {
            //new condition
            if (response.data.fileNameArr != null) {
              // this.loading=false
              this.listedFileName = response.data.fileNameArr;
              // this.MessageBar('S', 'File Data Fetched Succesully')
            }
            //here push path
            this.showContainer = true;
            this.isLoading = false;
            this.loading = false;
          } else {
            this.MessageBar("E", response.data.errMsg);
            this.isLoading = false;
            this.loading = false;
          }
        })
        .catch((error) => {
          this.MessageBar("E", error);
          console.log(error);
          this.isLoading = false;
          this.loading = false;
        });
      }
      this.downloadServer=this.Status
      this.Status=''
      lDate.fromDate = this.fromDate;
      lDate.toDate = this.toDate;
    },
  },
};
</script> -->
