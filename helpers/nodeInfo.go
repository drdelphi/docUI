package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func copyNodeInfo(status *statusMessage, info *nodeInfo) {
	info.nodeName = status.Status.Erd_node_display_name
	info.nodeType = status.Status.Erd_node_type
	info.shardID = status.Status.Erd_shard_id
	info.version = status.Status.Erd_app_version
	info.blockKey = status.Status.Erd_public_key_block_sign
	info.txKey = status.Status.Erd_public_key_tx_sign
	info.peers = status.Status.Erd_num_connected_peers
	info.cpuLoadPercent = status.Status.Erd_cpu_load_percent
	info.memLoadPercent = status.Status.Erd_mem_load_percent
	info.memUsedGo = status.Status.Erd_mem_used_golang
	info.memTotal = status.Status.Erd_mem_total
	info.netRecvBps = status.Status.Erd_network_recv_bps
	info.netRecvPeak = status.Status.Erd_network_recv_bps_peak
	info.netSendBps = status.Status.Erd_network_sent_bps
	info.netSendPeak = status.Status.Erd_network_sent_bps_peak
	info.isSyncing = status.Status.Erd_is_syncing
	info.syncedRound = status.Status.Erd_nonce
	info.nonce = status.Status.Erd_probable_highest_nonce
}

func InitializeNodes() int {
	nodes = make([]*nodeInfo, 0)
	for port := 8080; port < 8086; port++ {
		status, err := getNodeStatus(port)
		if err != nil {
			continue
		}
		node := nodeInfo{port: port}
		copyNodeInfo(status, &node)
		nodes = append(nodes, &node)
	}
	return len(nodes)
}

func getNodeStatus(port int) (*statusMessage, error) {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("http://localhost:%v/node/status", port), nil)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status statusMessage
	json.Unmarshal(body, &status)
	resp.Body.Close()
	return &status, nil
}

func GetNodesInfo() {
	for _, node := range nodes {
		status, err := getNodeStatus(node.port)
		if err != nil {
			continue
		}
		copyNodeInfo(status, node)
	}
}
