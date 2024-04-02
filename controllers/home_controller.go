package controllers

import (
    "fmt"
    "net/http"
)

func Home_Controller (w http.ResponseWriter, r *http.Request){
        fmt.Println("HelloWorld")
    }
