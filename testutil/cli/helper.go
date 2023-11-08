package cli

import (
	"os"
	"testing"

	"ghostcloud/testutil/network"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/client/cli"
	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
)

const (
	NewDescription = "new description"
	NewDomain      = "newdomain"
	NewContent     = "<h1>new content</h1>"
	IndexHTML      = "index.html"
)

func CreateDeployment(t *testing.T, nc *network.Context, id int, commonFlags []string) *types.Deployment {
	archive, err := sample.CreateTempArchive(IndexHTML, sample.HelloWorldHTMLBody)
	require.NoError(t, err)
	defer os.RemoveAll(archive.Name())

	deployment := &types.Deployment{
		Meta: sample.CreateMetaWithAddr(nc.Val.Address.String(), id),
		Dataset: &types.Dataset{
			Items: []*types.Item{
				{
					Meta: &types.ItemMeta{Path: IndexHTML},
					Content: &types.ItemContent{
						Content: []byte(sample.HelloWorldHTMLBody),
					},
				},
			},
		},
	}

	args := append([]string{deployment.GetMeta().GetName(), archive.Name()}, commonFlags...)
	_, err = clitestutil.ExecTestCLICmd(nc.Ctx, cli.CmdCreateDeployment(), args)
	require.NoError(t, err)
	return deployment
}
