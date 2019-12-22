package helpers

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	nodesTable                 *widgets.Table
	resourcesGrid, syncGrid    *ui.Grid
	syncGauges                 []*widgets.Gauge
	cpuLoadGauge, memLoadGauge *widgets.Gauge

	netSentSparkLine, netRecvSparkLine *widgets.Sparkline
	netSparkLineGroup                  *widgets.SparklineGroup
	netSentHistory, netRecvHistory     []float64

	noncesGrid *ui.Grid
	plotData   [][][]float64
	plotWidget []*widgets.Plot
)

func InitDisplay() error {
	err := ui.Init()
	if err != nil {
		return err
	}
	go func() {
		defer ui.Close()
		displayInfo()
	}()
	return nil
}

func initNodesTable() {
	n := len(nodes)
	termWidth, _ := ui.TerminalDimensions()
	colWidth := (termWidth-1)/n - 1
	nodesTable = widgets.NewTable()
	nodesTable.ColumnWidths = make([]int, n)
	for i := 0; i < n; i++ {
		nodesTable.ColumnWidths[i] = colWidth
	}
	for i := 0; i < 6; i++ {
		nodesTable.Rows = append(nodesTable.Rows, make([]string, n))
	}
	nodesTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	nodesTable.RowStyles[0] = ui.NewStyle(ui.ColorCyan)
	nodesTable.RowSeparator = false
	nodesTable.SetRect(0, 0, n*(colWidth+1)+1, 8)
}

func initSyncGrid() {
	n := len(nodes)
	termWidth, _ := ui.TerminalDimensions()
	colWidth := (termWidth - 1) / n
	items := make([]interface{}, n)
	syncGauges = make([]*widgets.Gauge, 0)
	for i := 0; i < n; i++ {
		gauge := widgets.NewGauge()
		gauge.Title = "Sync"
		gauge.SetRect(i*colWidth, 8, (i+1)*colWidth-1, 11)
		items[i] = ui.NewCol(1.0/float64(n), gauge)
		syncGauges = append(syncGauges, gauge)
	}
	syncGrid = ui.NewGrid()
	syncGrid.SetRect(0, 8, colWidth*n, 11)
	syncGrid.Set(ui.NewRow(1.0, items...))
}

func initNonceGrids() {
	n := len(nodes)
	termWidth, _ := ui.TerminalDimensions()
	colWidth := (termWidth - 1) / n

	items := make([]interface{}, n)
	plotWidget = make([]*widgets.Plot, 0)
	plotData = make([][][]float64, n)
	for i := 0; i < n; i++ {
		plot := widgets.NewPlot()
		shard := "META"
		if nodes[i].shardID < 5 {
			shard = fmt.Sprintf("%v", nodes[i].shardID)
		}
		plot.Title = "Shard: " + shard
		plot.SetRect(i*colWidth, 29, (i+1)*colWidth-1, 41)
		items[i] = ui.NewCol(1.0/float64(n), plot)
		plotWidget = append(plotWidget, plot)
		plotData[i] = make([][]float64, 1)
		plotData[i][0] = make([]float64, 2)
		plotData[i][0][0] = float64(nodes[i].nonce)
		plotData[i][0][1] = float64(nodes[i].nonce)
		plot.Data = plotData[i]
	}

	noncesGrid = ui.NewGrid()
	noncesGrid.SetRect(0, 29, colWidth*n, 41)
	noncesGrid.Set(ui.NewRow(1.0, items...))
}

func initLoadGrids() {
	n := len(nodes)
	termWidth, _ := ui.TerminalDimensions()
	colWidth := (termWidth - 1) / n

	cpuLoadGauge = widgets.NewGauge()
	cpuLoadGauge.Title = "CPU Load"
	cpuLoadGauge.SetRect(0, 11, colWidth*n, 14)

	memLoadGauge = widgets.NewGauge()
	memLoadGauge.Title = "RAM Load"
	memLoadGauge.SetRect(0, 14, colWidth*n, 17)
}

func initNetworkSparkLines() {
	n := len(nodes)
	termWidth, _ := ui.TerminalDimensions()
	colWidth := (termWidth - 1) / n

	netSentHistory = make([]float64, 0)
	netRecvHistory = make([]float64, 0)

	netRecvSparkLine = widgets.NewSparkline()
	netRecvSparkLine.Title = "Recv"
	netRecvSparkLine.LineColor = ui.ColorMagenta

	netSentSparkLine = widgets.NewSparkline()
	netSentSparkLine.Title = "Sent"
	netSentSparkLine.LineColor = ui.ColorCyan

	netSparkLineGroup = widgets.NewSparklineGroup(netRecvSparkLine, netSentSparkLine)
	netSparkLineGroup.Title = "Network Load"
	netSparkLineGroup.SetRect(0, 17, colWidth*n, 29)
}

func displayInfo() {
	initNodesTable()
	initSyncGrid()
	initLoadGrids()
	initNetworkSparkLines()
	initNonceGrids()

	uiEvents := ui.PollEvents()
	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-time.After(time.Millisecond * 500):
			refreshWindow()
		case <-sigTerm:
			ui.Clear()
			return
		case e := <-uiEvents:
			processUiEvents(e)
		}
	}
}

func processUiEvents(e ui.Event) {
	switch e.ID {
	case "<Resize>":
		doResizeEvent(e)
	case "<C-c>":
		ui.Close()
		AppTerminated = true
		return
	}
}

func doResizeEvent(e ui.Event) {
	payload := e.Payload.(ui.Resize)
	doResize(payload.Width, payload.Height)
}

func doResize(width int, height int) {
	initNodesTable()
	initSyncGrid()
	initLoadGrids()
	initNetworkSparkLines()
	initNonceGrids()
	refreshWindow()
}

func refreshWindow() {
	// mutex lock
	updateData()
	// mutex unlock
	ui.Clear()
	ui.Render(nodesTable, syncGrid, cpuLoadGauge, memLoadGauge, netSparkLineGroup, noncesGrid)
}

func updateData() {
	cpuLoadGauge.Percent = 0
	memLoadGauge.Percent = 0
	var recv, sent float64
	for i, node := range nodes {
		shard := "META"
		if node.shardID < 5 {
			shard = fmt.Sprintf("%v", node.shardID)
		}
		version := strings.Split(node.version, "-")[0]
		nodesTable.Rows[0][i] = node.nodeName
		nodesTable.Rows[1][i] = version
		nodesTable.Rows[2][i] = node.nodeType
		nodesTable.Rows[3][i] = "Shard: " + shard
		nodesTable.Rows[4][i] = node.txKey
		nodesTable.Rows[5][i] = node.blockKey

		if node.nonce > 0 {
			syncGauges[i].Percent = int(100 * node.syncedRound / node.nonce)
		} else {
			syncGauges[i].Percent = 0
		}
		if syncGauges[i].Percent > 100 {
			syncGauges[i].Percent = 100
		}
		if node.syncedRound < node.nonce {
			syncGauges[i].BarColor = ui.ColorYellow
			syncGauges[i].LabelStyle.Fg = ui.ColorRed
		} else {
			syncGauges[i].BarColor = ui.ColorGreen
			syncGauges[i].LabelStyle.Fg = ui.ColorWhite
		}

		cpuLoadGauge.Percent += int(node.cpuLoadPercent)
		memLoadGauge.Percent += int(node.memLoadPercent)
		recv += float64(node.netRecvBps)
		sent += float64(node.netSendBps)
	}
	if cpuLoadGauge.Percent > 100 {
		cpuLoadGauge.Percent = 100
	}
	if cpuLoadGauge.Percent < 51 {
		cpuLoadGauge.BarColor = ui.ColorGreen
		cpuLoadGauge.LabelStyle.Fg = ui.ColorWhite
	} else {
		if cpuLoadGauge.Percent < 76 {
			cpuLoadGauge.BarColor = ui.ColorYellow
			cpuLoadGauge.LabelStyle.Fg = ui.ColorBlack
		} else {
			cpuLoadGauge.BarColor = ui.ColorRed
			cpuLoadGauge.LabelStyle.Fg = ui.ColorWhite
		}
	}
	if memLoadGauge.Percent > 100 {
		memLoadGauge.Percent = 100
	}
	if memLoadGauge.Percent < 51 {
		memLoadGauge.BarColor = ui.ColorGreen
		memLoadGauge.LabelStyle.Fg = ui.ColorWhite
	} else {
		if memLoadGauge.Percent < 76 {
			memLoadGauge.BarColor = ui.ColorYellow
			memLoadGauge.LabelStyle.Fg = ui.ColorBlack
		} else {
			memLoadGauge.BarColor = ui.ColorRed
			memLoadGauge.LabelStyle.Fg = ui.ColorWhite
		}
	}

	var maxSent, maxRecv float64
	n := len(nodes)
	termWidth, _ := ui.TerminalDimensions()
	colWidth := (termWidth - 1) / n
	termWidth = colWidth*n - 2
	netRecvHistory = append(netRecvHistory, recv)
	if len(netRecvHistory) > termWidth {
		remove := len(netRecvHistory) - termWidth
		netRecvHistory = netRecvHistory[remove:]
	}
	netSentHistory = append(netSentHistory, sent)
	if len(netSentHistory) > termWidth {
		remove := len(netSentHistory) - termWidth
		netSentHistory = netSentHistory[remove:]
	}
	for _, v := range netSentHistory {
		if v > maxSent {
			maxSent = v
		}
	}
	for _, v := range netRecvHistory {
		if v > maxRecv {
			maxRecv = v
		}
	}
	netRecvSparkLine.Title = fmt.Sprintf("Recv %.3f Mbps (Peak %.3f Mpbs)", recv/1048576, maxRecv/1048576)
	netRecvSparkLine.Data = netRecvHistory
	netSentSparkLine.Title = fmt.Sprintf("Sent %.3f Mbps (Peak %.3f Mpbs)", sent/1048576, maxSent/1048576)
	netSentSparkLine.Data = netSentHistory

	if time.Now().Second()%6 == 0 {
		colWidth -= 5
		for i, node := range nodes {
			plotData[i][0] = append(plotData[i][0], float64(node.nonce))
			if len(plotData[i][0]) > colWidth {
				remove := len(plotData[i][0]) - colWidth
				plotData[i][0] = plotData[i][0][remove:]
			}
			data := make([][]float64, 1)
			data[0] = make([]float64, len(plotData[i][0]))
			for j := 0; j < len(plotData[i][0]); j++ {
				data[0][j] = plotData[i][0][j] - plotData[i][0][0]
			}
			plotWidget[i].Data = data
		}
	}
}
