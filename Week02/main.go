package main

import (
	"Week02/service"
	"log"
)

func main(){
	id := 100
	_ , err := service.HandleQuery(id)
	if err != nil {
		log.Printf("出问题了: %+v",err)
	}
}


