<template>
  <div>
    <v-form ref="form" lazy-validation>
  <v-container >
  <v-card class="pt-8 elevation-0 blue lighten-5" >
  <v-row dense>
    <v-col class="pl-1">
      <v-autocomplete
                class="mx-5"
                label="Program Name"
                :menu-props="{ bottom: true, offsetY: true }"
                background-color="white lighten-3"
                :items="ProgramNo"
                item-text="productCode"
                item-value="descripition"
                v-model="ProgramCode"
                @change="getFilePathAndServer"
                outlined
                dense
                ></v-autocomplete>
    </v-col>
    <v-col class="pl-16 mt-2">
      <v-btn @click="fetchData()" color="primary" small elevation="2">
        Search
      </v-btn>
    </v-col>
  </v-row>
</v-card>
</v-container>
</v-form>
<v-overlay :value="overlay">
    <v-progress-circular indeterminate size="64"></v-progress-circular>
  </v-overlay>

</div>
</template>

<script>
import ProgramService from "../../../services/EventService"
export default {
  data() {
    return {
      overlay:false,
      menu1: false,
      menu2: false,
     lDate:{path:"",server:'',programCode:''},
    ProgramNo:[],
    ProgramCode:'',
    };
  
  },
  // props:{
  //   reset:Boolean,
  // },
  // watch:{
  //   reset(){
  //     if(this.reset){
  //       this.$refs.form.reset()
  //     }
  //   }
  // },
  methods: {
    //v-auto complete to select the values
    getFilePathAndServer() {
      const myString = this.ProgramCode.split(',');
        // this.lDate.path = myString[0];
        this.lDate.path="/home/user/Music/toml"
        console.log("this.lDate.path",this.lDate.path);
        this.lDate.server = myString[1];
        console.log("this.lDate.server",this.lDate.server);

        this.lDate.programCode=myString[2]
        console.log("this.lDate.programCode",this.lDate.programCode);

      },
  
    //emit the value of file path from date to date to the logfile download page
    fetchData(){
      this.overlay=true
      // if (this.lDate.fromDate!='' & this.lDate.toDate!='' & this.ProgramCode!=''){
      if (this.ProgramCode!=''){
            // this.MessageBar('S', 'Searching Data')
          this.$emit('getData', this.lDate)
          this.overlay=false
      }else{
          this.overlay=false
        this.MessageBar('E', 'Invalid Credentials')
      }
      // this.ProgramCode=''
      this.overlay=false
    },
  //This method is used to show the Program Code to the Select Box
  FetchProgramCode(){
      this.overlay=true
    ProgramService.ProgramCode()
        .then((response) =>{
          if (response.data.status == "S") {
            this.ProgramNo=response.data.productCodeArr
        // this.MessageBar('S', 'Program Code Fetched successfully')
        // });
        this.overlay=false
          }else{
              this.overlay=false
              this.MessageBar('E', response.data.errMsg)
            console.log("Fetching data failed");
          }}
        ).catch((error)=>{
          this.overlay=false
          this.MessageBar('E', error)
          // console.log("Error i :", error) 
      });
  }
},
mounted(){
//to shown page when load time
this.FetchProgramCode()
    } 
};
</script>


<!-- <template>
  <div>
    <v-form ref="form" lazy-validation>
  <v-container >
  <v-card class="pt-8 elevation-0 blue lighten-5" >
  <v-row dense>
    <v-col class="pl-2">
      <v-autocomplete
                class="mx-5"
                label="Program Code"
                :menu-props="{ bottom: true, offsetY: true }"
                background-color="white lighten-3"
                :items="ProgramNo"
                item-text="productCode"
                item-value="descripition"
                v-model="ProgramCode"
                @change="getFilePathAndServer"
                outlined
                dense
                ></v-autocomplete>
    </v-col>
  <v-col class="pl-2">
  <v-menu
      v-model="menu1"
      :close-on-content-click="false"
      :nudge-right="40"
      transition="scale-transition"
      offset-y
      min-width="auto"
    >
      <template v-slot:activator="{ on, attrs }">
        <v-text-field
          v-model="lDate.fromDate"
          label="From Date"
          prepend-icon="mdi-calendar"
          readonly
          outlined
          v-bind="attrs"
          v-on="on"
          dense
          background-color="white lighten-3"
        ></v-text-field>
      </template>
      <v-date-picker
        v-model="lDate.fromDate"
        @input="menu1 = false"
      ></v-date-picker>
    </v-menu>
    </v-col>
    <v-col class="pl-8">
      <v-menu
      v-model="menu2"
      :close-on-content-click="false"
      :nudge-right="40"
      transition="scale-transition"
      offset-y
      min-width="auto"
    >
      <template v-slot:activator="{ on, attrs }">
        <v-text-field
          v-model="lDate.toDate"
          label="To Date"
          prepend-icon="mdi-calendar"
          readonly
          outlined
          v-bind="attrs"
          v-on="on"
          dense
          background-color="white lighten-3"
          ></v-text-field>
      </template>
      <v-date-picker
        v-model="lDate.toDate"
         :max="new Date().toISOString().substr(0, 10)"
        @input="menu2 = false"
      ></v-date-picker>
    </v-menu>
    </v-col>
    <v-col class="pl-16 mt-2">
      <v-btn @click="fetchData()" color="primary" small elevation="2">
        Search
      </v-btn>
    </v-col>
  </v-row>
</v-card>
</v-container>
</v-form>
<v-overlay :value="overlay">
    <v-progress-circular indeterminate size="64"></v-progress-circular>
  </v-overlay>

</div>
</template>

<script>
import ProgramService from "../../../services/EventService"
export default {
  data() {
    return {
      overlay:false,
      menu1: false,
      menu2: false,
     lDate:{ fromDate: "", toDate:"",path:"",server:'',programCode:''},
    ProgramNo:[],
    ProgramCode:'',
    };
  
  },
  // props:{
  //   reset:Boolean,
  // },
  // watch:{
  //   reset(){
  //     if(this.reset){
  //       this.$refs.form.reset()
  //     }
  //   }
  // },
  methods: {
    //v-auto complete to select the values
    getFilePathAndServer() {
      const myString = this.ProgramCode.split(',');
        this.lDate.path = myString[0];
        this.lDate.server = myString[1];
        this.lDate.programCode=myString[2]
    },
  
    //emit the value of file path from date to date to the logfile download page
    fetchData(){
      this.overlay=true
      if (this.lDate.fromDate!='' & this.lDate.toDate!='' & this.ProgramCode!=''){
            // this.MessageBar('S', 'Searching Data')
          this.$emit('getData', this.lDate)
          this.overlay=false
      }else{
          this.overlay=false
        this.MessageBar('E', 'Invalid Credentials')
      }
      // this.ProgramCode=''
      this.overlay=false
    },
  //This method is used to show the Program Code to the Select Box
  FetchProgramCode(){
      this.overlay=true
    ProgramService.ProgramCode()
        .then((response) =>{
          if (response.data.status == "S") {
            this.ProgramNo=response.data.productCodeArr
        // this.MessageBar('S', 'Program Code Fetched successfully')
        // });
        this.overlay=false
          }else{
              this.overlay=false
              this.MessageBar('E', response.data.errMsg)
            console.log("Fetching data failed");
          }}
        ).catch((error)=>{
          this.overlay=false
          this.MessageBar('E', error)
          // console.log("Error i :", error) 
      });
  }
},
mounted(){
//to shown page when load time
this.FetchProgramCode()
    } 
};
</script>
 -->
