package main

import (
    "math/rand"
)

func reportResponds(fromUser string) map[string]interface{} {
        praise := []string{"There you go", "Keep up the good work", "Keep it up", "Good job", "Awesome", "Top Banana", "Just Wow", "Brilliant", "Amazing work", "Top notch"}
        var msTeamsRespondMessage = make(map[string]interface{})
        msTeamsRespondMessage["text"] = "Acknowledged!" + praise[rand.Intn(len(praise))] + ", " + fromUser + "! Thank you!"
        msTeamsRespondMessage["type"] = "message"
        return msTeamsRespondMessage
}
