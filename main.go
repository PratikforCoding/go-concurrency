package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

const apiKey string = "7108ee70131e48b29d0164114232709"

type WeatherData struct {
    City   string `json:"-"`    // Not present in response, manually set
    Current struct {
        TempC float64 `json:"temp_c"`
    } `json:"current"`
}

func fetchWeather(city string, ch chan<- WeatherData, wg *sync.WaitGroup) error {
    url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city)

    resp, err := http.Get(url)
    if err != nil {
        fmt.Printf("Error fetching weather for %s: %s\n", city, err)
        wg.Done()
        return err
    }

    defer resp.Body.Close()

    var data WeatherData
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        fmt.Printf("Error decoding weather for %s: %s\n", city, err)
        wg.Done()
        return err
    }

    data.City = city  // Set city name
    ch <- data
    wg.Done()
    return nil
}

func main() {
    startNow := time.Now()
    cities := []string{"Kolkata", "London", "Paris", "Tokyo"}

    ch := make(chan WeatherData)
    var wg sync.WaitGroup
    for _, city := range cities {
        wg.Add(1)
        go fetchWeather(city, ch, &wg)
    }

    go func() {
        wg.Wait()
        close(ch)
    }()

    for weatherData := range ch {
        fmt.Printf("The temperature in %s is %f°C\n", weatherData.City, weatherData.Current.TempC)
    }
    fmt.Println("This operation took:", time.Since(startNow))
}

func fetchWeather1(city string) (WeatherData, error) {
    url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city)

    resp, err := http.Get(url)
    if err != nil {
        fmt.Printf("Error fetching weather for %s: %s\n", city, err)
        return WeatherData{}, err
    }

    defer resp.Body.Close()

    var data WeatherData
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        fmt.Printf("Error decoding weather for %s: %s\n", city, err)
        return WeatherData{}, err
    }
	data.City = city
    return data, nil
}

// func main() {
//     startNow := time.Now()
//     cities := []string{"Kolkata", "London", "Paris", "Tokyo"}

//     for _, city := range cities {
//         data, err := fetchWeather1(city)
//         if err != nil {
//             fmt.Println("Error fetching weather:", err)
//             continue
//         }
//         fmt.Printf("The temperature in %s is %f°C\n", data.City, data.Current.TempC)
//     }

//     fmt.Println("This operation took:", time.Since(startNow))
// }
