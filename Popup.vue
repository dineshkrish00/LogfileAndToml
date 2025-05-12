<template>
    <div class="text-center">
      <v-dialog v-model="dialog" persistent width="1000">
        <v-card height="750">
          <v-toolbar color="primary" height="50">
            <v-btn text small class="primary ml-n2" fab @click="closetab">
              <v-icon class="white--text" size="25">mdi-close</v-icon>
            </v-btn>
            <v-flex class="white--text font-weight-bold text-subtitle-1  d-flex justify-start mt-1">
              Toml File Values</v-flex>
          </v-toolbar>
                <v-card-title>{{File}} Values</v-card-title>
                <v-card-text>
                <v-textarea  class="font-weight-bold d-flex justify-start mt-1" v-model="fileData" :readonly="ReadCheck"   rows="19" row-height="20"
                :label="File" outlined ></v-textarea>
                <div class="d-flex justify-end">
                    <v-btn v-if="Check" class="mr-10 blue white--text"  rounded @click="EditEnable"><v-icon class="white--text mr-3" size="15">mdi-pencil</v-icon>Edit</v-btn>
                    <v-btn v-if="Check!=true" class="mr-10 blue white--text" rounded @click="Preview"><v-icon class="white--text mr-3" size="15">mdi-eye</v-icon>Read-Only</v-btn>
                    <v-btn class="mr-10 blue white--text" :disabled="ValueIsNotNull" rounded @click="SetFileData"><v-icon class="white--text mr-3" size="15">mdi-content-save</v-icon>Save</v-btn>
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
        fileData:'',
      }
    },
    props: {
        dialog:Boolean,
        fileContent:String,
        File:String
    },
    methods: {
        SetFileData(){
            this.$emit('closetab')
            console.log(" this.fileContent", this.fileData)
            this.$emit('save-file-data', this.fileData, this.File);

        },
        Preview(){
            this.Check=true
            this.ReadCheck=true
        },
        EditEnable(){
             this.Check=false
            this.ReadCheck=false
            console.log("this.ReadCheck")
        },
      closetab() {
        this.$emit('closetab')
      },
    },
    watch:{
        fileContent(){
            this.fileData=this.fileContent
        }
    },
    computed:{
      ValueIsNotNull(){
        return (this.fileData!="")?false:true
      }
    }

  }
  </script>