package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/sensepost/godoh/dnsclient"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var dnsDomain string
var dnsProviderName string
var dnsProvider dnsclient.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "godoh",
	Short: "A DNS (over-HTTPS) C2",
	Long:  `A DNS (over-HTTPS) C2`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(validateDNSProvider)
	cobra.OnInitialize(validateDNSDomain)
	cobra.OnInitialize(seedRand)

	rootCmd.PersistentFlags().StringVarP(&dnsDomain,
		"domain", "d", "", "DNS Domain to use. (ie: example.com)")
	rootCmd.PersistentFlags().StringVarP(&dnsProviderName,
		"provider", "p", "google", "Preferred DNS provider to use. [possible: google, cloudflare, raw]")
}

func seedRand() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func validateDNSDomain() {
	if dnsDomain == "" {
		log.Fatalf("A DNS domain to use is required.")
	}

	if strings.HasPrefix(dnsDomain, ".") {
		log.Fatalf("The DNS domain should be the base FQDN (without a leading dot).")
	}

	log.Infof("Using %s as DNS domain\n", dnsDomain)
}

func validateDNSProvider() {
	switch dnsProviderName {
	case "google":
		dnsProvider = dnsclient.NewGoogleDNS()
		break
	case "cloudflare":
		dnsProvider = dnsclient.NewCloudFlareDNS()
		break
	case "raw":
		dnsProvider = dnsclient.NewRawDNS()
		break
	default:
		log.Fatalf("DNS provider `%s` is not valid.\n", dnsProviderName)
	}

	log.Infof("Using `%s` as preferred provider\n", dnsProviderName)
}
