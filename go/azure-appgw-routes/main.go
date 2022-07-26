package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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
	Name            string
	Type            string
	Listener        Listener
	BackendAddr     []string
	BackendSettings *BackendSettings
	Routes          []Route
}
type Route struct {
	Name            string
	Paths           []string
	BackendAddr     []string
	BackendSettings *BackendSettings
	RedirectURL     string
}

func rules(appgw *armnetwork.ApplicationGateway) []Rule {
	rules := []Rule{}
	for _, rr := range appgw.Properties.RequestRoutingRules {
		rules = append(rules, Rule{
			Name:            *rr.Name,
			Type:            string(*rr.Properties.RuleType),
			Listener:        listener(appgw, subResourceID(rr.Properties.HTTPListener)),
			BackendAddr:     backendPoolAddrs(appgw, subResourceID(rr.Properties.BackendAddressPool)),
			BackendSettings: backendSettings(appgw, subResourceID(rr.Properties.BackendHTTPSettings)),
			Routes:          routes(appgw, subResourceID(rr.Properties.URLPathMap)),
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
					Name:            *pathRule.Name,
					Paths:           paths,
					BackendAddr:     backendPoolAddrs(appgw, subResourceID(pathRule.Properties.BackendAddressPool)),
					BackendSettings: backendSettings(appgw, subResourceID(pathRule.Properties.BackendHTTPSettings)),
					RedirectURL:     redirectURL(appgw, subResourceID(pathRule.Properties.RedirectConfiguration)),
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

type BackendSettings struct {
	Name     string
	Port     int32
	Protocol string
	Probe    Probe
}

func backendSettings(appgw *armnetwork.ApplicationGateway, id string) *BackendSettings {
	for _, bs := range appgw.Properties.BackendHTTPSettingsCollection {
		if *bs.ID != id {
			continue
		}

		return &BackendSettings{
			Name:     *bs.Name,
			Port:     *bs.Properties.Port,
			Protocol: string(*bs.Properties.Protocol),
			Probe:    probe(appgw, subResourceID(bs.Properties.Probe)),
		}

	}
	return nil
}

type Probe struct {
	Name string
	Host string
	Port int32
	Path string
}

func probe(appgw *armnetwork.ApplicationGateway, id string) Probe {
	for _, pb := range appgw.Properties.Probes {
		if *pb.ID != id {
			continue
		}

		port := int32(-1)
		if pb.Properties.Port != nil {
			port = *pb.Properties.Port
		}
		return Probe{
			Name: *pb.Name,
			Host: *pb.Properties.Host,
			Port: port,
			Path: *pb.Properties.Path,
		}
	}
	return Probe{}
}

func redirectURL(appgw *armnetwork.ApplicationGateway, id string) string {
	for _, rc := range appgw.Properties.RedirectConfigurations {
		if *rc.ID != id {
			continue
		}
		url := ""
		if rc.Properties.TargetURL != nil {
			url = *rc.Properties.TargetURL
		}
		return url
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

		host := ""
		if ln.Properties.HostName != nil {
			host = *ln.Properties.HostName
		} else if ln.Properties.HostNames != nil {
			xs := []string{}
			for _, x := range ln.Properties.HostNames {
				xs = append(xs, *x)
			}
			host = strings.Join(xs, ",")
		}

		return Listener{
			Host:     host,
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
