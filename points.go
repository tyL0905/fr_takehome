package main

import (
    "errors"
    "math"
    "regexp"
    "strconv"
    "strings"
    "time"
)

func calculatePoints(receipt Receipt) (int, error) {
    // try parsing first, throw if there's error
    total, err := strconv.ParseFloat(receipt.Total, 64)
    if err != nil {
        return 0, errors.New("invalid total amount format")
    }

    date, err := time.Parse("2006-01-02", receipt.PurchaseDate)
    if err != nil {
        return 0, errors.New("invalid purchase date format")
    }

    purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
    if err != nil {
        return 0, errors.New("invalid purchase time format")
    }

    for _, item := range receipt.Items {
        if _, err := strconv.ParseFloat(item.Price, 64); err != nil {
            return 0, errors.New("invalid item price format")
        }
    }

    // all the input is legal, start calc points
    points := 0

    // Rule 1
    alphanumeric := regexp.MustCompile(`[a-zA-Z0-9]`)
    points += len(alphanumeric.FindAllString(receipt.Retailer, -1))

    // Rule 2
    if total == float64(int(total)) {
        points += 50
    }

    // Rule 3
    if math.Mod(total, 0.25) == 0 {
        points += 25
    }

    // Rule 4
    points += (len(receipt.Items) / 2) * 5

    // Rule 5
    for _, item := range receipt.Items {
        trimmedDesc := strings.TrimSpace(item.ShortDescription)
        if len(trimmedDesc)%3 == 0 {
            price, _ := strconv.ParseFloat(item.Price, 64) 
            points += int(math.Ceil(price * 0.2))
        }
    }

    // Rule 6
    if date.Day()%2 != 0 {
        points += 6
    }

    // Rule 7
    if (purchaseTime.Hour() == 14 && purchaseTime.Minute() > 0) || purchaseTime.Hour() == 15 {
        points += 10
    }

    return points, nil
}
