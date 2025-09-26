package main

import (
	"fmt"
	"peach/utils"
)

const create_sql = `
create table if not exists test(
	name    text,
	age     int
)
`

func main() {
	defer utils.Recover()
	fmt.Println("hello world.")

}
