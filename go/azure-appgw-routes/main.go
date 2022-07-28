package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func run(subscriptionID string, appgwName string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}
	client, err := armnetwork.NewApplicationGatewaysClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}
	if err := dumpRoutes(client, appgwName); err != nil {
		return err
	}

	return nil
}

func dumpRoutes(client *armnetwork.ApplicationGatewaysClient, appgwName string) error {
	appgw, err := findAppGw(client, appgwName)
	if err != nil {
		return err
	}
	rs := rules(appgw)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&rs); err != nil {
		return fmt.Errorf("failed to encode roules: %w", err)
	}

	return nil
}

func findAppGw(client *armnetwork.ApplicationGatewaysClient, appgwName string) (*armnetwork.ApplicationGateway, error) {
	ctx := context.Background()
	pager := client.NewListAllPager(&armnetwork.ApplicationGatewaysClientListAllOptions{})
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, appgw := range resp.ApplicationGatewayListResult.Value {
			if *appgw.Name == appgwName {
				return appgw, nil
			}
		}
	}
	return nil, fmt.Errorf("appgw %s not found", appgwName)
}

func subResourceID(res *armnetwork.SubResource) string {
	if res != nil && res.ID != nil {
		return *res.ID
	}
	return ""
}

type Rule struct {
	Name        string
	Type        string
	Listener    Listener
	BackendAddr []string
	Routes      []Route
}
type Route struct {
	Name        string
	Paths       []string
	BackendAddr []string
	RedirectURL string
}

func rules(appgw *armnetwork.ApplicationGateway) []Rule {
	rules := []Rule{}
	for _, rr := range appgw.Properties.RequestRoutingRules {
		rules = append(rules, Rule{
			Name:        *rr.Name,
			Type:        string(*rr.Properties.RuleType),
			Listener:    listener(appgw, subResourceID(rr.Properties.HTTPListener)),
			BackendAddr: backendPoolAddrs(appgw, subResourceID(rr.Properties.BackendAddressPool)),
			Routes:      routes(appgw, subResourceID(rr.Properties.URLPathMap)),
		})
	}
	return rules
}

func routes(appgw *armnetwork.ApplicationGateway, id string) []Route {
	for _, urlPathMap := range appgw.Properties.URLPathMaps {
		if *urlPathMap.ID != id {
			continue
		}

		routes := []Route{}
		for _, pathRule := range urlPathMap.Properties.PathRules {
			paths := []string{}
			for _, p := range pathRule.Properties.Paths {
				paths = append(paths, *p)
			}
			routes = append(
				routes,
				Route{
					Name:        *pathRule.Name,
					Paths:       paths,
					BackendAddr: backendPoolAddrs(appgw, subResourceID(pathRule.Properties.BackendAddressPool)),
					RedirectURL: redirectURL(appgw, subResourceID(pathRule.Properties.RedirectConfiguration)),
				})
		}
		return routes
	}
	return nil
}

func backendPoolAddrs(appgw *armnetwork.ApplicationGateway, id string) []string {
	for _, bp := range appgw.Properties.BackendAddressPools {
		if *bp.ID != id {
			continue
		}
		addrs := []string{}
		for _, a := range bp.Properties.BackendAddresses {
			if a.Fqdn != nil {
				addrs = append(addrs, *a.Fqdn)
			} else {
				addrs = append(addrs, *a.IPAddress)
			}
		}
		return addrs
	}
	return nil
}

func redirectURL(appgw *armnetwork.ApplicationGateway, id string) string {
	for _, rc := range appgw.Properties.RedirectConfigurations {
		if *rc.ID != id {
			continue
		}
		return *rc.Properties.TargetURL
	}
	return ""
}

type Listener struct {
	Host     string
	Protocol string
	Port     int32
}

func listener(appgw *armnetwork.ApplicationGateway, id string) Listener {
	for _, ln := range appgw.Properties.HTTPListeners {
		if *ln.ID != id {
			continue
		}

		return Listener{
			Host:     *ln.Properties.HostName,
			Protocol: string(*ln.Properties.Protocol),
			Port:     port(appgw, subResourceID(ln.Properties.FrontendPort)),
		}
	}
	return Listener{}
}

func port(appgw *armnetwork.ApplicationGateway, id string) int32 {
	for _, fp := range appgw.Properties.FrontendPorts {
		if *fp.ID != id {
			continue
		}
		return *fp.Properties.Port
	}
	return -1
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s <subscription-id> <appgw-name>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(flag.Arg(0), flag.Arg(1)); err != nil {
		log.Fatal(err)
	}
}
