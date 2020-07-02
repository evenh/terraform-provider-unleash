package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func RunWithUnleash(testFunc func(portNumber int) int) {
	commonIdentifier := randomIdentifier()

	// Containers definitely run in the background
	ctx := context.Background()

	// Define network
	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			CheckDuplicate: true,
			Name:           commonIdentifier,
			Driver:         "bridge",
			SkipReaper:     false,
		},
	})

	if err != nil {
		log.Fatal("Could not create Docker network: ", err)
	}

	log.Println("Container network created: " + commonIdentifier)

	// Define containers
	pgReq, pgPort := postgresRequest(commonIdentifier)
	unleashReq, unleashPort := unleashRequest(commonIdentifier, pgPort)

	// Start containers in order
	startedContainers := startContainers(ctx, pgReq, unleashReq)

	// Expose mapped unleash port
	actualUnleashPort, err := startedContainers[1].MappedPort(ctx, unleashPort)
	if err != nil {
		log.Fatal("Could not get mapped port of Unleash container: ", err)
	}

	// Run the actual tests
	log.Printf("Unleash container running at port %d", actualUnleashPort.Int())
	returnCode := testFunc(actualUnleashPort.Int())

	// Stop containers after test run
	log.Println("Tearing down containers")
	for _, c := range startedContainers {
		_ = c.Terminate(ctx)
	}
	_ = network.Remove(ctx)

	os.Exit(returnCode)
}

func startContainers(context context.Context, reqs ...testcontainers.ContainerRequest) []testcontainers.Container {
	var containers []testcontainers.Container
	for _, req := range reqs {
		c, err := testcontainers.GenericContainer(context, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

		if err != nil {
			log.Fatal("Could not start container: "+req.Name, err)
		}

		containers = append(containers, c)
	}

	return containers
}

func postgresRequest(identifier string) (testcontainers.ContainerRequest, nat.Port) {
	port, _ := nat.NewPort("tcp", "5432")
	dbName := fmt.Sprintf("%s-db", identifier)
	req := testcontainers.ContainerRequest{
		Name:         dbName,
		Image:        "postgres:12.3-alpine",
		ExposedPorts: []string{port.Port()},
		Env: map[string]string{
			"POSTGRES_DB":               "db",
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		},
		Networks:   []string{identifier},
		WaitingFor: wait.ForListeningPort(port),
	}

	return req, port
}

func unleashRequest(identifier string, pgPort nat.Port) (testcontainers.ContainerRequest, nat.Port) {
	port, _ := nat.NewPort("tcp", "4242")
	req := testcontainers.ContainerRequest{
		Name:         fmt.Sprintf("%s-app", identifier),
		Image:        "unleashorg/unleash-server:3.3",
		ExposedPorts: []string{port.Port()},
		Env: map[string]string{
			"NODE_ENV":     "acceptance-test",
			"DATABASE_URL": fmt.Sprintf("postgres://postgres:unleash@%s:%s/db", identifier+"-db", pgPort.Port()),
		},
		Networks:   []string{identifier},
		WaitingFor: wait.ForLog("Unleash has started").WithStartupTimeout(30 * time.Second),
	}

	return req, port
}

func randomIdentifier() string {
	n := acctest.RandIntRange(1, 9999)
	return fmt.Sprintf("tf-unleash-%d", n)
}
