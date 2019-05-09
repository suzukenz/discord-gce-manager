package model

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

// GameServer is definition of the game server
type GameServer struct {
	IDName          string
	ShowName        string
	InstanceName    string
	Zone            string
	Port            int
	LastStatus      string
	EnablePortCheck bool
	Disable         bool
}

// SaveToDatastore puts game server definition to datastore.
func (gs *GameServer) SaveToDatastore(ctx context.Context, client *datastore.Client) error {
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

// CheckServerIsRunning returns game server is running or not
func (gs *GameServer) CheckServerIsRunning(service *compute.InstancesService, projectID string) (bool, error) {
	ins, err := service.Get(projectID, gs.Zone, gs.InstanceName).Do()
	if err != nil {
		return false, fmt.Errorf("error getting instance, err: %s", err)
	}

	return (ins.Status == "RUNNING"), nil
}

// CheckServerIsReady returns game server is ready or not
func (gs *GameServer) CheckServerIsReady(service *compute.InstancesService, projectID string) (externalIP string, status string, err error) {
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

// RunServer runs game server on compute engine.
func (gs *GameServer) RunServer(service *compute.InstancesService, projectID string) error {
	_, err := service.Start(projectID, gs.Zone, gs.InstanceName).Do()
	if err != nil {
		return err
	}
	return nil
}

// StopServer stops game server on compute engine.
func (gs *GameServer) StopServer(service *compute.InstancesService, projectID string) error {
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

// GetGameServers returns all enabled game server definitions from datastore.
func GetGameServers(ctx context.Context, client *datastore.Client) ([]*GameServer, error) {
	q := datastore.NewQuery(EntityName).Filter("Disable =", false)
	servers := make([]*GameServer, 0)
	_, err := client.GetAll(ctx, q, &servers)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

// GetGameServer returns single game server definition.
func GetGameServer(ctx context.Context, client *datastore.Client, idName string) (*GameServer, error) {
	k := datastore.NameKey(EntityName, idName, nil)

	server := new(GameServer)
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
