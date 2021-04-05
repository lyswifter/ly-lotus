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
const UpgradeCalicoHeight = 1000
const UpgradePersianHeight = 1005
const UpgradeOrangeHeight = 1010
const UpgradeClausHeight = 1015

const UpgradeActorsV3Height = 1020
const UpgradeNorwegianHeight = 1025

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
