package node

type NodeId struct {
	Hex string `json:"hex"`
}

type NodeStatus string

const (
	PendingDownload NodeStatus = "PendingDownload"
	ReadyForDownload = "ReadyForDownload"
	DownloadInProgress = "DownloadInProgress"
	DownloadCompleteAwaitingFinalSync = "DownloadCompleteAwaitingFinalSync"
	SnapshotCreation  = "SnapshotCreation"
	Ready = "Ready"
	Leaving = "Leaving"
	Offline = "Offline"
)

var	ValidStatuses = [...]NodeStatus{PendingDownload, ReadyForDownload, DownloadInProgress, DownloadCompleteAwaitingFinalSync, SnapshotCreation, Ready, Leaving, Offline}

func IsRedownloading(status NodeStatus) bool {
	return status == PendingDownload || status == ReadyForDownload || status == DownloadInProgress || status == DownloadCompleteAwaitingFinalSync
}

func IsOffline(status NodeStatus) bool {
	return status == Leaving || status == Offline
}

// TODO: replace with hashmap?
func asNodeStatus(in string) *NodeStatus {
	for _, name := range ValidStatuses {
		if name == NodeStatus(in) {
			return &name
		}
	}
	return nil
}

type NodeAddr struct {
	Host string `json:"host"`
	Port int `json:"port"`
}

type NodeInfo struct {
	Alias string `json:"alias"`
	Id NodeId `json:"id"`
	Ip NodeAddr `json:"ip"`
	Status NodeStatus `json:"status"`
	Reputation int `json:"reputation"`
}

type ClusterInfo []NodeInfo

type Metrics struct {
	Version string `json:"version"`
	NodeState NodeStatus `json:"nodeState"`

	TransactionServiceUnknownSize string `json:"transactionService_unknown_size"`
	GenesisAccepted string `json:"genesisAccepted"`
	ClusterOwnJoinedHeight string `json:"cluster_ownJoinedHeight"`
	NodeStateSuccess string `json:"nodeState_success"`
	TotalNumCBsInShapshots string `json:"totalNumCBsInShapshots"`
	LastSnapshotHash string `json:"lastSnapshotHash"`

	CheckpointTipsRemoved string `json:"checkpointTipsRemoved"`
	RedownloadMaxCreatedSnapshotHeight string `json:"redownload_maxCreatedSnapshotHeight"`
	PeerApiRXFinishedCheckpoint string `json:"peerApiRXFinishedCheckpoint"`
	ConsensusParticipateInRound string `json:"consensus_participateInRound"`
	NodeStartTimeMS string `json:"nodeStartTimeMS"`
	TransactionAccepted string `json:"transactionAccepted"`
	ExteralHost string `json:"externalHost"`
	SnapshotAttemptHeightIntervalNotMet string `json:"snapshotAttempt_heightIntervalNotMet"`

	RedownloadLastMajorityStateHeight string `json:"redownload_lastMajorityStateHeight"`
	AcceptMajorityCheckpointBlockUniquesCount1 string `json:"acceptMajorityCheckpointBlockUniquesCount_1"`
	ResolveMajorityCheckpointBlockUniquesCoun1 string `json:"resolveMajorityCheckpointBlockUniquesCount_1"`
	TpsAll string `json:"TPS_all"`
	ResolveMajorityCheckpointBlockProposalCount3 string `json:"resolveMajorityCheckpointBlockProposalCount_3"`
	ParentSOEServiceQueryFailed string `json:"parentSOEServiceQueryFailed"`
	PeerAddedFromRegistrationFlow string `json:"peerAddedFromRegistrationFlow"`
	RewardsSnapshotReward string `json:"rewards_snapshotReward"`
	RewardsSnapshotRewardWithoutStardust string `json:"rewards_snapshotRewardWithoutStardust"`
	TrustDataPollingRound string `json:"blacklistedAddressesSize"`
	BlacklistedAddressesSize string `json:"trustDataPollingRound"`
	RewardsStardustBalanceAfterReward string `json:"rewards_stardustBalanceAfterReward"`
	ConsensusParticipateInRoundInvalidNodeStateError string `json:"consensus_participateInRound_invalidNodeStateError"`
	DeadPeer string `json:"deadPeer"`
	AcceptedCBCacheMatchesAcceptedSize string `json:"acceptedCBCacheMatchesAcceptedSize"`
	AwaitingForAcceptance string `json:"awaitingForAcceptance"`
	SnapshotHeightIntervalConditionMet string `json:"snapshotHeightIntervalConditionMet"`
	NodeStartDate string `json:"nodeStartDate"`
	WriteSnapshotSuccess string `json:"writeSnapshot_success"`
	SnapshotWriteToDiskSuccess string `json:"snapshotWriteToDisk_success"`
	ChannelCount string `json:"channelCount"`
	ObservationServiceUnknownSize string `json:"observationService_unknown_size"`
	Alias string `json:"alias"`

	AddressCount string `json:"addressCount"`
	SnapshotAttempt_success string `json:"snapshotAttempt_success"`
	Id string `json:"id"`
	TransactionServiceAcceptedSize string `json:"transactionService_accepted_size"`
	CheckpointsAcceptedWithDummyTxs string `json:"checkpointsAcceptedWithDummyTxs"`
	DownloadFinishedTotal string `json:"downloadFinished_total"`
	NextSnapshotHeight string `json:"nextSnapshotHeight"`
	AllowedForAcceptance string `json:"allowedForAcceptance"`
	NodeCurrentDate string `json:"nodeCurrentDate"`
	RewardsSelfBalanceAfterReward string `json:"rewards_selfBalanceAfterReward"`

	SnapshotCount string `json:"snapshotCount"`
	ObservationServiceInConsensusSize string `json:"observationService_inConsensus_size"`
	TPSLast10seconds string `json:"TPS_last_10_seconds"`
	AddPeerWithRegistrationSymmetricSuccess string `json:"addPeerWithRegistrationSymmetric_success"`
	BalancesBySnapshot string `json:"balancesBySnapshot"`
	RewardsLastRewardedHeight string `json:"rewards_lastRewardedHeight"`
	MissingParents string `json:"missingParents"`

	RedownloadLastSentHeight string `json:"redownload_lastSentHeight"`
	MinTipHeight string `json:"minTipHeight"`
	ConsensusStartOwnRoundConsensusError string `json:"consensus_startOwnRound_consensusError"`
	CheckpointTipsIncremented string `json:"checkpointTipsIncremented"`
	HeightCalculationParentMissing string `json"heightCalculationParentMissing`
	ConsensusStartOwnRound string `json:"consensus_startOwnRound`
	Address string `json:"address"`
	CheckpointValidationSuccess string `json:"checkpointValidationSuccess"`
	RedownloadMaxAcceptedSnapshotHeight string `json:"redownload_maxAcceptedSnapshotHeight"`
	ConsensusStartOwnRoundUnknownError string `json:"consensus_startOwnRound_unknownError"`
	AcceptedCBSinceSnapshot string `json:"acceptedCBSinceSnapshot"`
	CheckpointAccepted string `json:"checkpointAccepted"`
	AcceptMajorityCheckpointBlockSelectedCount3 string `json:"acceptMajorityCheckpointBlockSelectedCount_3"`
	TransactionAcceptedFromRedownload string `json:"transactionAcceptedFromRedownload"`
	BadPeerAdditionAttempt string `json:"badPeerAdditionAttempt"`
	NodeCurrentTimeMS string `json:"nodeCurrentTimeMS"`
	BatchTransactionsEndpoint string `json:"batchTransactionsEndpoint"`
	ConsensusStartOwnRoundConsensusStartError string `json:"consensus_startOwnRound_consensusStartError"`
	RewardsSelfSnapshotReward string `json:"rewards_selfSnapshotReward"`
	MetricsRound string `json:"metricsRound"`
	GenesisHash string `json:"genesisHash"`
	Peers string `json:"peers"`
	RewardsSnapshotCount string `json:"rewards_snapshotCount"`
	ReDownloadFinishedTotal string `json:"reDownloadFinished_total"`
	SyncBufferSize string `json:"syncBufferSize"`
	LastSnapshotHeight string `json:"lastSnapshotHeight"`
	ObservationServiceAcceptedSize string `json:"ObservationService_accepted_size"`
	ObservationServicePendingSize string `json:"observationService_pending_size"`
	HeightBelow string `json:"heightBelow"`
	TransactionServicePendingSize string `json:"transactionService_pending_size"`
	RewardsStardustSnapshotReward string `json:"rewards_stardustSnapshotReward"`
	SnapshotHeightIntervalConditionNotMet string `json:"snapshotHeightIntervalConditionNotMet"`
	CheckpointAcceptBlockAlreadyStored string `json:"checkpointAcceptBlockAlreadyStored"`
	ActiveTips string `json:"activeTips"`
	TransactionServiceInConsensusSize string `json:"transactionService_inConsensus_size"`
}

type MetricsEnvelope struct {
	Metrics Metrics `json:"metrics"`
}