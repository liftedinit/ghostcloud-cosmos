package network

import (
	"fmt"
	"testing"
	"time"

	"ghostcloud/testutil/keeper"
	"ghostcloud/testutil/sample"
	"ghostcloud/x/ghostcloud/types"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	tmdb "github.com/cometbft/cometbft-db"
	tmrand "github.com/cometbft/cometbft/libs/rand"

	"ghostcloud/app"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	cosmosnetwork "github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type (
	Network = network.Network
	Config  = network.Config
	Context struct {
		Net *network.Network
		Val *cosmosnetwork.Validator
		Ctx client.Context
	}

	TxTestCase struct {
		Name string
		Args []string
		Err  error
		Code uint32
	}
)

const FlagPattern = "--%s=%s"

func setupCommon(t *testing.T, cfg Config) *Context {
	t.Helper()
	net := New(t, cfg)
	val := net.Validators[0]
	ctx := val.ClientCtx
	return &Context{
		Net: net,
		Val: val,
		Ctx: ctx,
	}
}

func setupTxCommonFlags(t *testing.T, nc *Context, addr string) []string {
	t.Helper()
	return []string{
		fmt.Sprintf(FlagPattern, flags.FlagFrom, addr),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf(FlagPattern, flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(FlagPattern, flags.FlagFees, sdk.NewCoins(sdk.NewCoin(nc.Net.Config.BondDenom, sdkmath.NewInt(10))).String()),
	}
}

func SetupTxCommonFlags(t *testing.T, nc *Context) []string {
	t.Helper()
	return setupTxCommonFlags(t, nc, nc.Val.Address.String())
}

func SetupQueryCommonFlags(t *testing.T) []string {
	t.Helper()
	return []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
}

func Setup(t *testing.T) *Context {
	return setupCommon(t, DefaultConfig())
}

func SetupWithDeployments(t *testing.T, n int) (*Context, []*types.Deployment) {
	t.Helper()
	return SetupWithDeploymentsAndAddr(t, n, sample.AccAddress())
}

func SetupWithDeploymentsAndAddr(t *testing.T, n int, addr string) (*Context, []*types.Deployment) {
	t.Helper()
	cfg := DefaultConfig()
	state := types.GenesisState{
		Params:      types.DefaultParams(),
		Deployments: sample.CreateNDeploymentsWithAddr(addr, n, keeper.DATASET_SIZE),
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	ctx := setupCommon(t, cfg)
	return ctx, state.Deployments
}

// New creates instance with fully configured cosmos network.
// Accepts optional config, that will be used in place of the DefaultConfig() if provided.
func New(t *testing.T, configs ...Config) *Network {
	if len(configs) > 1 {
		panic("at most one config should be provided")
	}
	var cfg network.Config
	if len(configs) == 0 {
		cfg = DefaultConfig()
	} else {
		cfg = configs[0]
	}
	net, err := network.New(t, t.TempDir(), cfg)
	require.NoError(t, err)
	_, err = net.WaitForHeight(1)
	require.NoError(t, err)
	t.Cleanup(net.Cleanup)
	return net
}

// DefaultConfig will initialize config for the network with custom application,
// genesis and single validator. All other parameters are inherited from cosmos-sdk/testutil/network.DefaultConfig
func DefaultConfig() network.Config {
	var (
		encoding = app.MakeEncodingConfig()
		chainID  = "chain-" + tmrand.NewRand().Str(6)
	)
	return network.Config{
		Codec:             encoding.Marshaler,
		TxConfig:          encoding.TxConfig,
		LegacyAmino:       encoding.Amino,
		InterfaceRegistry: encoding.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor: func(val network.ValidatorI) servertypes.Application {
			return app.New(
				val.GetCtx().Logger,
				tmdb.NewMemDB(),
				nil,
				true,
				map[int64]bool{},
				val.GetCtx().Config.RootDir,
				0,
				encoding,
				simtestutil.EmptyAppOptions{},
				baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
				baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
				baseapp.SetChainID(chainID),
			)
		},
		GenesisState:    app.ModuleBasics.DefaultGenesis(encoding.Marshaler),
		TimeoutCommit:   2 * time.Second,
		ChainID:         chainID,
		NumValidators:   1,
		BondDenom:       sdk.DefaultBondDenom,
		MinGasPrices:    fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		AccountTokens:   sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
		StakingTokens:   sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction),
		BondedTokens:    sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
		PruningStrategy: pruningtypes.PruningOptionNothing,
		CleanupDir:      true,
		SigningAlgo:     string(hd.Secp256k1Type),
		KeyringOptions:  []keyring.Option{},
	}
}
