package main

func main() {
	newsiteChannel := make(chan NewSiteEvent)
	forever := make(chan bool)
	dequeueEvents(newsiteChannel)
	consumeEvents(newsiteChannel)
	<-forever
}
