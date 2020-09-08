package nodegrid

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
		return "Offline"
	}
	return fmt.Sprintf("%s", metrics.NodeState)
}

func BuildImageOutput(target string, clusterOverview []NodeOverview, grid map[string]map[string]node.NodeInfo, outputTheme string) {
	baseXMargin := float64(4)

	ordOffsetX := baseXMargin
	nameOffsetX := ordOffsetX + 15
	addrOffsetX := nameOffsetX + 100
	versionOffsetX := addrOffsetX + 100
	snapshotOffsetX := versionOffsetX + 50
	latencyOffsetX := snapshotOffsetX + 50
	statusLbOffsetX := latencyOffsetX + 100
	statusSeparatorOffsetX := statusLbOffsetX + 5
	statusLocalOffsetX := statusSeparatorOffsetX + 5

	baseYMargin := float64(4)

	gridSize := float64(len(clusterOverview))

	iconWidth := 10
	iconHeight := 10
	iconMargin := 4

	textWidth := float64(8)
	textHeight := float64(12)
	// textHeight := iconHeight

	canvasWidth := math.Max(float64(iconWidth + iconMargin) * gridSize + 2*baseXMargin + 2*textWidth, 550)

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

	var backgroundColor = canvas.Transparent
    var textColor = canvas.Black

    ctx.SetFillColor(backgroundColor)

	if outputTheme == "light" {
		backgroundColor = canvas.White
		textColor = canvas.Black
	}else if outputTheme == "dark" {
		backgroundColor = canvas.Black
		textColor = canvas.White
	}

	ctx.SetFillColor(backgroundColor)
	backgroundRectangle := canvas.Rectangle(canvasWidth, canvasHeight)
	ctx.DrawPath(0, 0, backgroundRectangle)
	ctx.SetFillColor(canvas.Transparent)

	textFace := fontFamily.Face(20.0, textColor, canvas.FontRegular, canvas.FontNormal)

	colNumbersHeight := float64(iconHeight + iconMargin)
	textPosY := canvasHeight - 1

	ctx.DrawText(ordOffsetX, textPosY, canvas.NewTextBox(textFace, "##", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	ctx.DrawText(nameOffsetX, textPosY, canvas.NewTextBox(textFace, "Alias", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	ctx.DrawText(addrOffsetX, textPosY, canvas.NewTextBox(textFace,  "Address", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	ctx.DrawText(versionOffsetX, textPosY, canvas.NewTextBox(textFace, "Version", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	ctx.DrawText(snapshotOffsetX, textPosY, canvas.NewTextBox(textFace,  "Snapshot", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	ctx.DrawText(latencyOffsetX, textPosY, canvas.NewTextBox(textFace,  "Latency", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	ctx.DrawText(statusLbOffsetX, textPosY, canvas.NewTextBox(textFace,  "Status Lb", 0.0, 0.0, canvas.Right, canvas.Top, 0.0, 0.0))
	ctx.DrawText(statusSeparatorOffsetX, textPosY, canvas.NewTextBox(textFace,  "/", 0.0, 0.0, canvas.Center, canvas.Top, 0.0, 0.0))

	ctx.DrawText(statusLocalOffsetX, textPosY, canvas.NewTextBox(textFace,  "Status Node", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))

	var lastTextPosY = float64(0)

	for i, nodeOverview := range clusterOverview {
		var version = "?"
		var snap = "?"
		var latency = "?"
		statusTextFace1 :=  statusColorRGB(nodeOverview.info.Status)
		var statusTextFace2 = statusTextFace1
		if nodeOverview.metrics != nil {
			version = nodeOverview.metrics.Version
			snap = nodeOverview.metrics.LastSnapshotHeight
			latency = fmtLatency(nodeOverview.metricsResponseDuration)
			statusTextFace2 = statusColorRGB(nodeOverview.metrics.NodeState)
		}

		offsetY := float64(i+1) * textHeight
		textPosY := canvasHeight - 1 - offsetY
		lastTextPosY = textPosY

		ctx.DrawText(ordOffsetX, textPosY, canvas.NewTextBox(textFace,  fmt.Sprintf("%02d", i), 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
		ctx.DrawText(nameOffsetX, textPosY, canvas.NewTextBox(textFace, nodeOverview.info.Alias, 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
		ctx.DrawText(addrOffsetX, textPosY, canvas.NewTextBox(textFace,  nodeOverview.info.Ip.Host, 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
		ctx.DrawText(versionOffsetX, textPosY, canvas.NewTextBox(textFace, version, 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
		ctx.DrawText(snapshotOffsetX, textPosY, canvas.NewTextBox(textFace,  snap, 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
		ctx.DrawText(latencyOffsetX, textPosY, canvas.NewTextBox(textFace,  latency, 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))

		ctx.DrawText(statusLbOffsetX, textPosY, canvas.NewTextBox(fontFamily.Face(20.0,statusTextFace1, canvas.FontRegular, canvas.FontNormal),  fmt.Sprintf("%s", nodeOverview.info.Status), 0.0, 0.0, canvas.Right, canvas.Top, 0.0, 0.0))
		ctx.DrawText(statusSeparatorOffsetX, textPosY, canvas.NewTextBox(textFace,  "/", 0.0, 0.0, canvas.Center, canvas.Top, 0.0, 0.0))
		ctx.DrawText(statusLocalOffsetX, textPosY, canvas.NewTextBox(fontFamily.Face(20.0, statusTextFace2, canvas.FontRegular, canvas.FontNormal),  nodeStatusString(nodeOverview.metrics), 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
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

	c.WriteFile(target, rasterizer.PNGWriter(2))
}