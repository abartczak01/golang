package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"zad03/models"
)

func LoadPosts(filename string, count int) (map[int]models.Post, error) {
    jsonFile, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("error opening JSON file: %v", err)
    }
    defer jsonFile.Close()

    byteValue, err := io.ReadAll(jsonFile)
    if err != nil {
        return nil, fmt.Errorf("error reading JSON file: %v", err)
    }

    var posts []models.Post
    err = json.Unmarshal(byteValue, &posts)
    if err != nil {
        return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
    }

    rand.Shuffle(len(posts), func(i, j int) {
        posts[i], posts[j] = posts[j], posts[i]
    })
    postMap := make(map[int]models.Post)
    for i := 0; i < count && i < len(posts); i++ {
        posts[i].ID = i + 1
        postMap[i+1] = posts[i]
    }

    fmt.Println("Załadowano posty")

    return postMap, nil
}