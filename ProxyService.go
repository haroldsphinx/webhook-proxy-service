package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/namsral/flag"
	"github.com/webhook-proxy-service/proxy"
)

var (
	flagSet       = flag.NewFlagSetWithEnvPrefix(os.Args[0], "GWP", 0)
	listenAddress = flagSet.String("listen", ":8080", "Listening port for the webhook service")
	secret        = flagSet.String("secret", "", "Secret of the webhook API")
	upstreamURL   = flagSet.String("upstreamURL", "", "URL to which the proxy request will be forwarded (required)")
	provider      = flagSet.String("provider", "gitlab", "Since this is meant to be a general webhook-service, its just proper we set different type of usage here, but for now our concern is gitlab")
	allowedPaths  = flagSet.String("allowedPaths", "", "Comma separated list of the paths we want to accept from the client")
	ignoredUsers  = flagSet.String("ignoredUsers", "", "Comma separated list of silent or banned users")
	allowedUsers  = flagSet.String("allowedUsers", "", "Comma separated list of allowed users")
)

func validateRequiredFlags() {
	isValid := true
	if len(strings.TrimSpace(*upstreamURL)) == 0 {
		log.Println("Required flag 'upstreamURL' not specified")
		isValid = false
	}

	if !isValid {
		fmt.Println("")
		flagSet.Usage()
		fmt.Println("")
		panic("See Flag Usage")
	}
}

func main() {
	flagSet.Parse(os.Args[1:])
	validateRequiredFlags()
	lowerProvider := strings.ToLower(*provider)

	//split , into an array
	allowedPathsArray := []string{}
	if len(*allowedPaths) > 0 {
		allowedPathsArray = strings.Split(*ignoredUsers, ",")
	}

	ignoredUsersArray := []string{}
	if len(*ignoredUsers) > 0 {
		ignoredUsersArray = strings.Split(*ignoredUsers, ",")
	}

	log.Printf("Quidax Webhook Service Started")
	p, err := proxy.NewProxy(*upstreamURL, allowedPathsArray, lowerProvider, *secret, ignoredUsersArray)
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Run(*listenAddress); err != nil {
		log.Fatal(err)
	}

}
