package reporter

import (
	"constellation_cli/pkg/node"
	"fmt"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/rasterizer"
	"image/color"
	"math"
)


func statusColorRGB(status node.NodeStatus) color.RGBA {
	switch status {
	case node.DownloadCompleteAwaitingFinalSync:
		return color.RGBA{R: 247, G: 247, B: 52, A: 255}
	case node.ReadyForDownload:
		return color.RGBA{R: 247, G: 247, B: 52, A: 255}
	case node.DownloadInProgress:
		return color.RGBA{R: 247, G: 247, B: 52, A: 255}
	case node.PendingDownload:
		return color.RGBA{R: 247, G: 247, B: 52, A: 255}
	case node.Leaving:
		return color.RGBA{R: 230, G: 78, B: 18, A: 255}
	case node.SnapshotCreation:
		return color.RGBA{R: 78, G: 191, B: 189, A: 255}
	case node.Ready:
		return color.RGBA{R: 98, G: 191, B: 67, A: 255}
	default:
		return color.RGBA{R: 230, G: 78, B: 18, A: 255}
	}
}

func nodeStatusString(metrics *node.Metrics) string {
	if metrics == nil {
		return "/Offline"
	}
	return fmt.Sprintf("/%s", metrics.NodeState)
}

func BuildImageOutput(target string, clusterOverview []NodeOverview, grid map[string]map[string]node.NodeInfo) {
	baseXMargin := float64(4)
	baseYMargin := float64(4)

	gridSize := float64(len(clusterOverview))

	iconWidth := 10
	iconHeight := 10
	iconMargin := 4

	textWidth := float64(8)
	textHeight := float64(12)
	// textHeight := iconHeight

	canvasWidth := math.Max(float64(iconWidth + iconMargin) * gridSize + 2*baseXMargin + 2*textWidth, 500)

	canvasHeight := float64(iconHeight + iconMargin) * gridSize + 2*baseYMargin + textHeight + (gridSize+2) * textHeight

	c := canvas.New(canvasWidth, canvasHeight)

	ctx := canvas.NewContext(c)

	nodeIcon := canvas.RoundedRectangle(float64(iconWidth), float64(iconHeight), 1)
	// fontFamily := canvas.NewFontFamily("Courier Regular")

	var fontFamily *canvas.FontFamily
	fontFamily = canvas.NewFontFamily("Monospace")
	//fontFamily.Use(canvas.CommonLigatures)
	if err := fontFamily.LoadLocalFont("Monospace", canvas.FontRegular); err != nil {
		panic(err)
	}

	textFace := fontFamily.Face(20.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	colNumbersHeight := float64(iconHeight + iconMargin)


	headerText := fmt.Sprintf("#  %-30s %-21s %-10s %-10s %-10s %s", "Alias", "Address", "Version", "Snapshot", "Latency", "Status Lb/Node")

	tb := canvas.NewTextBox(textFace, headerText, 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)

	textPosY := canvasHeight - 1

	ctx.DrawText(baseXMargin, textPosY, tb)

	var lastTextPosY = float64(0)

	for i, nodeOverview := range clusterOverview {
		var version = "?"
		var snap = "?"
		var latency = "?"

		if nodeOverview.metrics != nil {
			version = nodeOverview.metrics.Version
			snap = nodeOverview.metrics.LastSnapshotHeight
			latency = fmtLatency(nodeOverview.metricsResponseDuration)
		}

		rowText := fmt.Sprintf("%02d %-30s %-21s %-10s %-10s %-10s %s%s", i, nodeOverview.info.Alias,
			fmt.Sprintf("%s:%d", nodeOverview.info.Ip.Host, nodeOverview.info.Ip.Port),
			version,
			snap,
			latency,
			nodeOverview.info.Status,
			nodeStatusString(nodeOverview.metrics))

		tb := canvas.NewTextBox(textFace,  rowText, 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)

		offsetY := float64(i+1) * textHeight

		textPosY := canvasHeight - 1 - offsetY

		lastTextPosY = textPosY

		ctx.DrawText(baseXMargin, textPosY, tb)
	}

	lastTextPosY = lastTextPosY - textHeight *2

	for col, _ := range clusterOverview {
		tb := canvas.NewTextBox(textFace, fmt.Sprintf("%02d", col), 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)

		textPosY := lastTextPosY - 2
		textPosX := textWidth + float64((iconWidth + iconMargin) * col) + baseXMargin + textWidth

		ctx.DrawText(textPosX, textPosY, tb)
	}

	for row, rowNode := range clusterOverview {

		offsetY := colNumbersHeight + float64((iconHeight + iconMargin) * row)

		rowMap := grid[rowNode.info.Ip.Host]

		tb := canvas.NewTextBox(textFace, fmt.Sprintf("%02d", row), 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)

		textPosY := lastTextPosY - offsetY

		ctx.DrawText(baseXMargin, textPosY, tb)

		for col, colNode := range clusterOverview {

			offsetX := textWidth + float64((iconWidth + iconMargin) * col)

			cell := rowMap[colNode.info.Ip.Host]

			ctx.SetFillColor(statusColorRGB(cell.Status))
			ctx.SetStrokeColor(statusColorRGB(cell.Status))

			iconPosY :=  lastTextPosY - (offsetY + float64(iconHeight))

			ctx.DrawPath(baseXMargin + offsetX + textWidth, iconPosY, nodeIcon.Copy())
		}
	}

	c.WriteFile(target, rasterizer.PNGWriter(5))
}