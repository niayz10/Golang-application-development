package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

func main(){

	var name string
	var userRating string

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter your full name")
	name, _ = reader.ReadString('\n')

	reader = bufio.NewReader(os.Stdin)
	fmt.Println("Please rate our Dosa center between 1 and 5: ")
	userRating, _ = reader.ReadString('\n')
	mynumRating, _ := strconv.ParseFloat(strings.TrimSpace(userRating), 64)

	Rating, ok := ValidityOfInput(mynumRating)
	if ok != nil{
		panic("not correct input")
	}

	fmt.Printf("Hello %v,  \nThanks for rating our dosa center with %v star rating. \n\n Your rating was recorded in our system at %v\n\n", name, Rating, time.Now().Format(time.Stamp))
	Response(Rating)
	
}

func ValidityOfInput(mynumRating float64) (float64, error){

	if mynumRating<1 || mynumRating>5 {
		return 0, fmt.Errorf("Not corect input")
	}
	return mynumRating, nil;
}

func Response(mynumRating float64){
	if mynumRating == 5 {
		fmt.Println("Bonus for team for 5 star service")
	} 
	 if mynumRating <= 4 && mynumRating >= 3 {
		fmt.Println("We are always improving")
	} 
	 if mynumRating < 3 {
		fmt.Println("Need Serious work on our side")
	}
}