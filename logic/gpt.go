package main

import (
    "context"
    "fmt"
    "log"
    "google.golang.org/genai"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
    ctx := context.Background()
    client, err := genai.NewClient(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash",
        genai.Text(""),
        nil,
    )

    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result.Text())
}




// Не имеет мотивов миграции
// Имеет четкие цели достижение которым поможет US
// Не имеет жалоб к совей стране и ее (власти, оброзовании итд)
