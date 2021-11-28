package vault

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/testcontainers/testcontainers-go"
// 	"github.com/testcontainers/testcontainers-go/wait"
// )

// type nginxContainer struct {
// 	testcontainers.Container
// 	URI string
// }

// func setupNginx(ctx context.Context) (*nginxContainer, error) {
// 	req := testcontainers.ContainerRequest{
// 		Image:        "vault",
// 		ExposedPorts: []string{"8200/tcp"},
// 		WaitingFor:   wait.ForHTTP("/"),
// 		Cmd:          []string{"server", "-dev", "-dev-root-token-id", "root"},
// 	}
// 	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
// 		ContainerRequest: req,
// 		Started:          true,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	ip, err := container.Host(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mappedPort, err := container.MappedPort(ctx, "8200")
// 	if err != nil {
// 		return nil, err
// 	}

// 	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

// 	return &nginxContainer{Container: container, URI: uri}, nil
// }

// func TestPrint(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	ctx := context.Background()

// 	nginxC, err := setupNginx(ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Clean up the container after the test is complete
// 	defer nginxC.Terminate(ctx)

// 	v, err := NewClient()
// 	assert.NoError(t, err)

// 	m, err := v.ReadSecrets("secret", "")

// 	assert.NoError(t, err)

// 	fmt.Println(m)
// }
