<template>
    <div class="text-center">
      <v-dialog v-model="Filedialog" persistent width="1000">
        <v-card height="690">
          <v-toolbar color="primary" height="50">
            <v-btn text small class="primary ml-n2" fab @click="closetab">
              <v-icon class="white--text" size="25">mdi-close</v-icon>
            </v-btn>
            <v-flex class="white--text font-weight-bold text-subtitle-1  d-flex justify-start mt-1">
              Add Toml File</v-flex>
          </v-toolbar>
                <v-card-text class="mt-5">
                <v-textarea  class="font-weight-bold d-flex justify-start mt-1" v-model="fileName" 
                rows="1" row-height="1"
                :rules="[v => v.endsWith('.toml') || 'File name must end with .toml']"
                ref="fileNameInput" 
                label="Toml File Name" outlined ></v-textarea>
                <v-textarea  class="font-weight-bold d-flex justify-start mt-1" v-model="fileData"   rows="15" row-height="10"
                label="Toml Values" outlined ></v-textarea>
                <div class="d-flex justify-end">
                    <v-btn class="mr-10 blue white--text" :disabled="FieldValidation" rounded @click="AddFileData">
                    <v-icon class="white--text mr-3" size="15">mdi-content-save</v-icon>Save</v-btn>
                </div>
        </v-card-text>
        </v-card>
      </v-dialog>
    </div>
  </template>

  <script>
  export default {
    data() {
      return {
        ReadCheck:true,
        Check:true,
        fileName:'',
        fileData:'',
      }
    },
    props: {
        Filedialog:Boolean,
    },
    methods: {
        AddFileData(){
            this.$emit('Addclosetab')
            console.log(" this.fileContent", this.fileData)
            console.log("this.fileName",this.fileName)
            this.$emit('Add-file-data',  this.fileData,this.fileName);
            this.fileData=''
            this.fileName=''
        },
      closetab() {
        this.$emit('Addclosetab')
        this.$refs.fileNameInput.resetValidation();
            this.fileData=''
            this.fileName=''
      },
    },

    computed:{
      ValueIsNotNull(){
        return (this.fileData!="")?false:true
      },
      FieldValidation(){
        return (this.fileName!="" && this.fileData !="")?false:true
      }
    }

  }
  </script>