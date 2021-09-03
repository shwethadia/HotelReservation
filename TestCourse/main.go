package main

import (
	"errors"
	"fmt"
)

func main(){


	res , err := divide(100,0)
	if err!= nil {
		fmt.Println(err)
	}

	fmt.Println(res)

}


func divide(x,y float32) (float32,error){

	var result float32

	if y == 0 {

		return result,errors.New("can't devide by zero")
	}
	result= x /y
	return result,nil
}