package sample

import (
	"archive/zip"
	"bytes"
	"ghostcloud/x/ghostcloud/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io"
	"math/rand"
	"strconv"
	"time"
)

const KB = 1024

// AccAddress returns a sample account address
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

func GetDeploymentMeta(i int) *types.DeploymentMeta {
	return &types.DeploymentMeta{
		Name:        strconv.Itoa(i),
		Description: strconv.Itoa(i),
		Domain:      strconv.Itoa(i),
	}
}

func GetDeploymentNameMeta(name string, i int) *types.DeploymentMeta {
	return &types.DeploymentMeta{
		Name:        name,
		Description: strconv.Itoa(i),
		Domain:      strconv.Itoa(i),
	}
}

func GetDeploymentFiles(i int) []*types.File {
	return []*types.File{
		{
			Meta:    &types.FileMeta{Name: strconv.Itoa(i)},
			Content: []byte{byte(i)},
		},
	}
}

func GetDeployment(i int) types.Deployment {
	return types.Deployment{
		Creator: AccAddress(),
		Meta:    GetDeploymentMeta(i),
		Files:   GetDeploymentFiles(i),
	}
}

func GetDeploymentList(i int) []types.Deployment {
	deployments := make([]types.Deployment, 0, i)
	for j := 0; j < i; j++ {
		deployments = append(deployments, GetDeployment(j))
	}
	return deployments
}

func GetDuplicatedDeploymentList() []types.Deployment {
	elem := GetDeployment(0)
	return []types.Deployment{elem, elem}
}

// generateRandomData generates a string with random data.
func generateRandomData(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// CreateZipWithHTML returns a []byte of a zip archive containing an index.html with random data.
func CreateZipWithHTML() []byte {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Add index.html to zip
	f, err := zipWriter.Create("index.html")
	if err != nil {
		panic(err)
	}

	// Write some random data to index.html
	randomData := generateRandomData(1 + rand.Intn(KB))
	_, err = io.WriteString(f, randomData)
	if err != nil {
		panic(err)
	}

	// Close the zip writer to finish the zip creation
	err = zipWriter.Close()
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
