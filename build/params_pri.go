// +build pri

package build

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/actors/policy"
)

var DrandSchedule = map[abi.ChainEpoch]DrandEnum{
	0: DrandMainnet,
}

const BootstrappersFile = "prinet.pi"
const GenesisFile = "https://prinet-car.xchain.icu:21970/prinet.car"

const UpgradeBreezeHeight = -1
const BreezeGasTampingDuration = 0

const UpgradeSmokeHeight = -1
const UpgradeIgnitionHeight = -2
const UpgradeRefuelHeight = -3
const UpgradeTapeHeight = -4

const UpgradeActorsV2Height = 3
const UpgradeLiftoffHeight = -5

const UpgradeKumquatHeight = 5
const UpgradeCalicoHeight = 6
const UpgradePersianHeight = 7
const UpgradeOrangeHeight = 8
const UpgradeClausHeight = 9

const UpgradeActorsV3Height = 10
const UpgradeNorwegianHeight = 11

var UpgradeActorsV4Height = abi.ChainEpoch(12)

func init() {
	MessageConfidence = 2
	policy.SetSupportedProofTypes(
		abi.RegisteredSealProof_StackedDrg2KiBV1,
		abi.RegisteredSealProof_StackedDrg512MiBV1,
		abi.RegisteredSealProof_StackedDrg32GiBV1,
		abi.RegisteredSealProof_StackedDrg64GiBV1,
	)
	policy.SetConsensusMinerMinPower(abi.NewStoragePower(2048))
	policy.SetMinVerifiedDealSize(abi.NewStoragePower(256))
	policy.SetPreCommitChallengeDelay(15)
	policy.SetChainFinality(90)
	policy.SetWPoStChallengeWindow(20)
	// InteractivePoRepConfidence = 3
	BuildType |= BuildPri
}

// Seconds
const BlockDelaySecs = uint64(30)

const PropagationDelaySecs = 6

// SlashablePowerDelay is the number of epochs after ElectionPeriodStart, after
// which the miner is slashed
//
// Epochs
const SlashablePowerDelay = 20

const BootstrapPeerThreshold = 1

// we skip checks on message validity in this block to sidestep the zero-bls signature
var WhitelistedBlock = MustParseCid("bafy2bzaceapyg2uyzk7vueh3xccxkuwbz3nxewjyguoxvhx77malc2lzn2ybi")
