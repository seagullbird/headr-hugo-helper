package main

func main() {
	newsiteChannel := make(chan NewSiteEvent)
	dequeueEvents(newsiteChannel)
	consumeEvents(newsiteChannel)
}
