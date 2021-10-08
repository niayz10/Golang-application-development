# Rating

This program rates the work of the Dosa Center.

### How it works
The program takes data from the user, such as name and rating. Then the program checks the rating for validity. Determines the date of the assessment and displays the answer from the data center.

### Functions of the program
    1. ValidityOfInput(mynumRating float64) (float64, error)
        Validate the input from the user.
    2. Response(mynumRating float64)
        This function gives an answer based on the rating

### Packets that were used
    ~
    "fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
    "testing"
    ~

### How to run a program
To run a project on your computer or laptop, you need to download the project from the repository. Then open it up in your work environment.Open terminal and write `go run Rating.go"`

### How to test a program
Open terminal and write `go test`