package helpers

var (
	nodes         []*nodeInfo
	AppTerminated bool = false
)

type nodeInfo struct {
	port           int
	nodeName       string
	nodeType       string
	shardID        uint64
	version        string
	blockKey       string
	txKey          string
	peers          uint64
	cpuLoadPercent uint64
	memLoadPercent uint64
	memUsedGo      uint64
	memTotal       uint64
	netRecvBps     uint64
	netRecvPeak    uint64
	netSendBps     uint64
	netSendPeak    uint64
	isSyncing      uint64
	syncedRound    uint64
	nonce          uint64
}

type statusType struct {
	Erd_app_version                                            string `json:"erd_app_version"`
	Erd_connected_nodes                                        uint64 `json:"erd_connected_nodes"`
	Erd_consensus_round_state                                  string `json:"erd_consensus_round_state"`
	Erd_consensus_state                                        string `json:"erd_consensus_state"`
	Erd_count_accepted_blocks                                  uint64 `json:"erd_count_accepted_blocks"`
	Erd_count_consensus                                        uint64 `json:"erd_count_consensus"`
	Erd_count_consensus_accepted_blocks                        uint64 `json:"erd_count_consensus_accepted_blocks"`
	Erd_count_leader                                           uint64 `json:"erd_count_leader"`
	Erd_cpu_load_percent                                       uint64 `json:"erd_cpu_load_percent"`
	Erd_current_block_hash                                     string `json:"erd_current_block_hash"`
	Erd_current_block_size                                     uint64 `json:"erd_current_block_size"`
	Erd_current_round                                          uint64 `json:"erd_current_round"`
	Erd_current_round_timestamp                                uint64 `json:"erd_current_round_timestamp"`
	Erd_fork_choice_count                                      uint64 `json:"erd_fork_choice_count"`
	Erd_highest_notarized_block_by_metachain_for_current_shard uint64 `json:"erd_highest_notarized_block_by_metachain_for_current_shard"`
	Erd_is_syncing                                             uint64 `json:"erd_is_syncing"`
	Erd_latest_tag_software_version                            string `json:"erd_latest_tag_software_version"`
	Erd_live_validator_nodes                                   uint64 `json:"erd_live_validator_nodes"`
	Erd_mem_load_percent                                       uint64 `json:"erd_mem_load_percent"`
	Erd_mem_total                                              uint64 `json:"erd_mem_total"`
	Erd_mem_used_golang                                        uint64 `json:"erd_mem_used_golang"`
	Erd_mem_used_sys                                           uint64 `json:"erd_mem_used_sys"`
	Erd_metric_community_percentage                            string `json:"erd_metric_community_percentage"`
	Erd_metric_consensus_group_size                            uint64 `json:"erd_metric_consensus_group_size"`
	Erd_metric_cross_check_block_height                        string `json:"erd_metric_cross_check_block_height"`
	Erd_metric_leader_percentage                               string `json:"erd_metric_leader_percentage"`
	Erd_metric_num_validators                                  uint64 `json:"erd_metric_num_validators"`
	Erd_mini_blocks_size                                       uint64 `json:"erd_mini_blocks_size"`
	Erd_network_recv_bps                                       uint64 `json:"erd_network_recv_bps"`
	Erd_network_recv_bps_peak                                  uint64 `json:"erd_network_recv_bps_peak"`
	Erd_network_recv_percent                                   uint64 `json:"erd_network_recv_percent"`
	Erd_network_sent_bps                                       uint64 `json:"erd_network_sent_bps"`
	Erd_network_sent_bps_peak                                  uint64 `json:"erd_network_sent_bps_peak"`
	Erd_network_sent_percent                                   uint64 `json:"erd_network_sent_percent"`
	Erd_node_display_name                                      string `json:"erd_node_display_name"`
	Erd_node_type                                              string `json:"erd_node_type"`
	Erd_nonce                                                  uint64 `json:"erd_nonce"`
	Erd_num_connected_peers                                    uint64 `json:"erd_num_connected_peers"`
	Erd_num_mini_blocks                                        uint64 `json:"erd_num_mini_blocks"`
	Erd_num_shard_headers_from_pool                            uint64 `json:"erd_num_shard_headers_from_pool"`
	Erd_num_shard_headers_processed                            uint64 `json:"erd_num_shard_headers_processed"`
	Erd_num_transactions_processed                             uint64 `json:"erd_num_transactions_processed"`
	Erd_num_tx_block                                           uint64 `json:"erd_num_tx_block"`
	Erd_probable_highest_nonce                                 uint64 `json:"erd_probable_highest_nonce"`
	Erd_public_key_block_sign                                  string `json:"erd_public_key_block_sign"`
	Erd_public_key_tx_sign                                     string `json:"erd_public_key_tx_sign"`
	Erd_rewards_value                                          string `json:"erd_rewards_value"`
	Erd_round_time                                             uint64 `json:"erd_round_time"`
	Erd_shard_id                                               uint64 `json:"erd_shard_id"`
	Erd_synchronized_round                                     uint64 `json:"erd_synchronized_round"`
	// Erd_tx_pool_load                                           uint64 `json:"erd_tx_pool_load"`
}

type statusMessage struct {
	Status statusType `json:"details"`
}
