package main

import (
    "fmt"    
	"project/internal/config"
)
func main() {    
	cfg := config.MustLoad()
    fmt.Println(cfg)
    // TODO: init logger : slog
    // TODO: init db : gorm
    // TODO: init router : chi, chirender
    // TODO: init storage : sqlIte
}