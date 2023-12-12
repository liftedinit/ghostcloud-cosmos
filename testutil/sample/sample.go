package sample

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"ghostcloud/x/ghostcloud/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const HelloWorldHTMLBody = "<html><body>Hello World</body></html>"

// AccAddress returns a sample account address
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

func CreateNDatasetPayloads(n int, datasetSize int) ([]*types.Meta, []*types.Payload) {
	metas := make([]*types.Meta, n)
	payloads := make([]*types.Payload, n)
	for i := 0; i < n; i++ {
		metas[i], payloads[i] = CreateDatasetPayload(i, datasetSize)
	}
	return metas, payloads
}

func CreateNDatasetPayloadsWithIndexHtml(n int, datasetSize int) ([]*types.Meta, []*types.Payload) {
	metas := make([]*types.Meta, n)
	payloads := make([]*types.Payload, n)
	for i := 0; i < n; i++ {
		metas[i], payloads[i] = CreateDatasetPayloadWithIndexHtml(i, datasetSize)
	}
	return metas, payloads
}

func CreateNArchivePayloads(n int) ([]*types.Meta, []*types.Payload) {
	metas := make([]*types.Meta, n)
	payloads := make([]*types.Payload, n)
	for i := 0; i < n; i++ {
		metas[i], payloads[i] = CreateArchivePayload(i)
	}
	return metas, payloads
}

func CreateNDeployments(n int, datasetSize int) []*types.Deployment {
	return CreateNDeploymentsWithAddr(AccAddress(), n, datasetSize)
}

func CreateNDeploymentsWithAddr(addr string, n int, datasetSize int) []*types.Deployment {
	deployments := make([]*types.Deployment, n)
	for i := 0; i < n; i++ {
		deployments[i] = CreateDeploymentWithAddr(addr, i, datasetSize)
	}
	return deployments
}

func CreateNMetaDataset(n int, datasetSize int) ([]*types.Meta, []*types.Dataset) {
	metas := make([]*types.Meta, n)
	datasets := make([]*types.Dataset, n)
	for i := 0; i < n; i++ {
		metas[i], datasets[i] = CreateMetaDataset(i, datasetSize)
	}
	return metas, datasets
}

func createDatasetPayload(addr string, i int, datasetSize int) (*types.Meta, *types.Payload) {
	return CreateMetaWithAddr(addr, i), &types.Payload{
		PayloadOption: &types.Payload_Dataset{Dataset: CreateDataset(datasetSize)},
	}
}

func CreateDatasetPayload(i int, datasetSize int) (*types.Meta, *types.Payload) {
	return createDatasetPayload(AccAddress(), i, datasetSize)
}

func CreateDatasetPayloadWithIndexHtml(i int, datasetSize int) (*types.Meta, *types.Payload) {
	return CreateMeta(i), &types.Payload{
		PayloadOption: &types.Payload_Dataset{Dataset: CreateDatasetWithIndexHtml(datasetSize)},
	}
}

func CreateDatasetPayloadWithAddr(addr string, i int, datasetSize int) (*types.Meta, *types.Payload) {
	return createDatasetPayload(addr, i, datasetSize)
}

func CreateDatasetPayloadWithAddrAndIndexHtml(addr string, i int, datasetSize int) (*types.Meta, *types.Payload) {
	return CreateMetaWithAddr(addr, i), &types.Payload{
		PayloadOption: &types.Payload_Dataset{Dataset: CreateDatasetWithIndexHtml(datasetSize)},
	}
}

func CreateArchivePayloadWithAddrAndIndexHtml(addr string, i int) (*types.Meta, *types.Payload) {
	return CreateMetaWithAddr(addr, i), &types.Payload{
		PayloadOption: &types.Payload_Archive{Archive: CreateArchive()},
	}
}

func CreateArchivePayload(i int) (*types.Meta, *types.Payload) {
	return CreateMeta(i), &types.Payload{
		PayloadOption: &types.Payload_Archive{Archive: CreateArchive()},
	}
}

func CreateRandomArchivePayload(i int, size int64, name string) (*types.Meta, *types.Payload) {
	body, err := generateRandomBytes(int(size))
	if err != nil {
		panic(err)
	}
	return CreateMeta(i), &types.Payload{
		PayloadOption: &types.Payload_Archive{Archive: &types.Archive{
			Type:    types.ArchiveType_Zip,
			Content: CreateZip(name, string(body)),
		}},
	}
}

func CreateBombArchivePayload(i int, size uint64, name string) (*types.Meta, *types.Payload) {
	return CreateMeta(i), &types.Payload{
		PayloadOption: &types.Payload_Archive{Archive: &types.Archive{
			Type:    types.ArchiveType_Zip,
			Content: CreateZip(name, strings.Repeat("a", int(size))),
		}},
	}
}

func createDeployment(addr string, i int, datasetSize int) *types.Deployment {
	return &types.Deployment{
		Meta:    CreateMetaWithAddr(addr, i),
		Dataset: CreateDataset(datasetSize),
	}
}
func CreateDeployment(i int, datasetSize int) *types.Deployment {
	return createDeployment(AccAddress(), i, datasetSize)
}

func CreateDeploymentWithAddr(addr string, i int, datasetSize int) *types.Deployment {
	return createDeployment(addr, i, datasetSize)
}

func CreateDeploymentWithAddrAndIndexHtml(addr string, i int, datasetSize int) *types.Deployment {
	return &types.Deployment{
		Meta:    CreateMetaWithAddr(addr, i),
		Dataset: CreateDatasetWithIndexHtml(datasetSize),
	}
}

func CreateMetaDataset(i int, datasetSize int) (*types.Meta, *types.Dataset) {
	return CreateMeta(i), CreateDataset(datasetSize)
}

func createMeta(addr string, i int) *types.Meta {
	return &types.Meta{
		Creator:     addr,
		Name:        strconv.Itoa(i),
		Description: strconv.Itoa(i),
		Domain:      strconv.Itoa(i),
	}
}
func CreateMeta(i int) *types.Meta {
	return createMeta(AccAddress(), i)
}

func CreateMetaWithAddr(addr string, i int) *types.Meta {
	return createMeta(addr, i)
}

func CreateMetaInvalidAddress() *types.Meta {
	return &types.Meta{
		Creator:     "invalid_address",
		Name:        "name",
		Description: "description",
		Domain:      "domain",
	}
}

func CreateDataset(n int) *types.Dataset {
	return &types.Dataset{Items: CreateNItems(n)}
}

func CreateDatasetWithIndexHtml(n int) *types.Dataset {
	return &types.Dataset{Items: CreateNItemsWithIndexHtml(n)}
}

func CreateArchive() *types.Archive {
	return &types.Archive{
		Type:    types.ArchiveType_Zip,
		Content: CreateZip("index.html", HelloWorldHTMLBody),
	}
}

func CreateZip(fileName string, body string) []byte {
	return createInMemoryZip(fileName, body)
}

func CreateItem(i int) *types.Item {
	return &types.Item{
		Meta:    &types.ItemMeta{Path: strconv.Itoa(i)},
		Content: &types.ItemContent{Content: []byte{byte(i)}},
	}
}

func CreateItemWithIndexHtml() *types.Item {
	return &types.Item{
		Meta:    &types.ItemMeta{Path: "index.html"},
		Content: &types.ItemContent{Content: []byte{0x00}},
	}
}

func CreateNItems(n int) []*types.Item {
	items := make([]*types.Item, n)
	for i := 0; i < n; i++ {
		items[i] = CreateItem(i)
	}
	return items
}

func CreateNItemsWithIndexHtml(n int) []*types.Item {
	items := make([]*types.Item, n)
	items[0] = CreateItemWithIndexHtml()
	for i := 1; i < n; i++ {
		items[i] = CreateItem(i)
	}
	return items
}

func createInMemoryZip(fileName string, body string) []byte {
	// Step 1: Create a buffer to hold the zip archive's data in memory
	var buffer bytes.Buffer

	// Step 2: Create a new zip archive writing to the buffer
	zipWriter := zip.NewWriter(&buffer)

	// Dummy data to put into the zip file
	files := []struct {
		Name, Body string
	}{
		{fileName, body},
	}

	// Step 3: Add files to the archive
	for _, file := range files {
		f, err := zipWriter.Create(file.Name)
		if err != nil {
			panic(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			panic(err)
		}
	}

	// Step 4: Close the zip archive
	err := zipWriter.Close()
	if err != nil {
		panic(err)
	}

	// Step 5: Convert buffer's contents to bytes
	return buffer.Bytes()
}

func CreateTempDataset() (dir string, err error) {
	dir, err = os.MkdirTemp("", "example")
	if err != nil {
		return dir, fmt.Errorf("error creating temporary directory: %v", err)
	}

	indexFilePath := filepath.Join(dir, "index.html")
	_, err = os.Create(indexFilePath)
	if err != nil {
		return dir, fmt.Errorf("error creating index.html file: %v", err)
	}

	return dir, nil
}

func CreateTempArchive(fileName string, body string) (file *os.File, err error) {
	file, err = os.CreateTemp("", "test-archive-*.zip")
	if err != nil {
		return file, fmt.Errorf("error creating temporary file: %v", err)
	}

	data := CreateZip(fileName, body)
	_, err = file.Write(data)
	if err != nil {
		return file, fmt.Errorf("error writing to temporary file: %v", err)
	}

	return file, nil
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
