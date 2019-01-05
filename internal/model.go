package internal

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	compute "google.golang.org/api/compute/v1"
)

// EntityName is a name of the entity on Google Cloud Datastore.
const EntityName = "GameServer"

// StatusReady... are status names.
const (
	StatusReady    = "READY"
	StatusNotReady = "NOT READY"
)

type gameServer struct {
	IDName          string
	ShowName        string
	InstanceName    string
	Zone            string
	Port            int
	LastStatus      string
	EnablePortCheck bool
	Disable         bool
}

func (gs *gameServer) saveToDatastore(ctx context.Context, client *datastore.Client) error {
	// (example)
	// gs := &gameServer{
	// 	IDName:          "7d2d",
	// 	ShowName:        "7 Days To Die",
	// 	InstanceName:    "seven-days-2d-v2",
	// 	Zone:            "asia-northeast1-b",
	// 	Port:            26900,
	// 	LastStatus:      "",
	// 	EnablePortCheck: true,
	// 	Disable:         false,
	// }
	// gs.saveToDatastore(ctx, datastoreClient)
	k := datastore.NameKey(EntityName, gs.IDName, nil)
	_, err := client.Put(ctx, k, gs)
	if err != nil {
		return err
	}
	return nil
}

func (gs *gameServer) checkServerIsRunning(service *compute.InstancesService) (bool, error) {
	ins, err := service.Get(projectID, gs.Zone, gs.InstanceName).Do()
	if err != nil {
		return false, fmt.Errorf("error getting instance, err: %s", err)
	}

	return (ins.Status == "RUNNING"), nil
}

func (gs *gameServer) checkServerIsReady(service *compute.InstancesService) (externalIP string, status string, err error) {
	ins, err := service.Get(projectID, gs.Zone, gs.InstanceName).Do()
	if err != nil {
		return "", StatusNotReady, fmt.Errorf("error getting instance, err: %s", err)
	}

	if ins.Status != "RUNNING" {
		return "", StatusNotReady, fmt.Errorf("instance is not running :%s", ins.Name)
	}

	if len(ins.NetworkInterfaces) == 0 {
		return "", StatusNotReady, fmt.Errorf("error, network interface is not attached :%s", ins.Name)
	}

	nwIf := ins.NetworkInterfaces[0]
	if len(nwIf.AccessConfigs) == 0 {
		return "", StatusNotReady, fmt.Errorf("error, nothing access configs on network interface, instance:%s, nwif: %s", ins.Name, nwIf.Name)
	}

	exIP := nwIf.AccessConfigs[0].NatIP
	if gs.EnablePortCheck {
		err = scanPort(exIP, gs.Port, 5*time.Second)
		if err != nil {
			return "", StatusNotReady, err
		}
	}

	return exIP, StatusReady, nil
}

func (gs *gameServer) runServer(service *compute.InstancesService) error {
	_, err := service.Start(projectID, gs.Zone, gs.InstanceName).Do()
	if err != nil {
		return err
	}
	return nil
}

func (gs *gameServer) stopServer(service *compute.InstancesService) error {
	_, err := service.Stop(projectID, gs.Zone, gs.InstanceName).Do()
	if err != nil {
		return err
	}
	return nil
}

func scanPort(ip string, port int, timeout time.Duration) error {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			return scanPort(ip, port, timeout)
		}
		return fmt.Errorf("%s is closed, err: %s", target, err)
	}
	defer conn.Close()
	return nil
}

func getGameServers(ctx context.Context, client *datastore.Client) ([]*gameServer, error) {
	q := datastore.NewQuery(EntityName).Filter("Disable =", false)
	servers := make([]*gameServer, 0)
	_, err := client.GetAll(ctx, q, &servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func getGameServer(ctx context.Context, client *datastore.Client, idName string) (*gameServer, error) {
	k := datastore.NameKey(EntityName, idName, nil)

	server := new(gameServer)
	err := client.Get(ctx, k, server)
	if err != nil {
		if strings.Contains(err.Error(), "no such entity") {
			// no such entity is not error
			return nil, nil
		}
		return nil, err
	}

	return server, nil
}
