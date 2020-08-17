package main

/*
Create a simple API endpoint which return true if number is prime

To use private repository or private git like github enterprise use ssh keys over username/password.
	The following command change all URL from https to git which use ssh keys

		git config --global url."git@github.com:".insteadOf "https://github.com/"

	Golang will also try to do a sum check, this will also need to be disable for private repo
		export GONOSUMDB=github.com/Tanmay-Teaches/golang
*/
import (
	"net/http"
	"github.com/Tanmay-Teaches/golang/chapter3/example3"
	"github.com/labstack/echo/v4"
	"strconv"
)


func main() {

	e := echo.New()
	e.GET("/:number", func(c echo.Context) error {
		nstr := c.Param("number")
		n, err := strconv.Atoi(nstr)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.String(http.StatusOK, strconv.FormatBool(example3.IsPrime(n)))
	})

	e.Logger.Fatal(e.Start(":1323"))
}
