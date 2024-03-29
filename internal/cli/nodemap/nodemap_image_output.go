package nodemap

import (
	"constellation/pkg/node"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/rasterizer"
	"image/color"
	"math"
)

func statusColorRGB(status node.NodeState) color.RGBA {

	switch status {
	case node.StartingSession:
		return color.RGBA{R: 247, G: 247, B: 52, A: 255}
	case node.ReadyToJoin:
		return color.RGBA{R: 247, G: 247, B: 52, A: 255}
	case node.WaitingForDownload:
		return color.RGBA{R: 247, G: 247, B: 52, A: 255}
	case node.LoadingGenesis:
		return color.RGBA{R: 247, G: 247, B: 52, A: 255}
	case node.Initial:
		return color.RGBA{R: 247, G: 247, B: 52, A: 255}
	case node.Leaving:
		return color.RGBA{R: 230, G: 78, B: 18, A: 255}
	case node.SessionStarted:
		return color.RGBA{R: 98, G: 191, B: 67, A: 255}
	case node.Observing:
		return color.RGBA{R: 96, G: 255, B: 253, A: 255}
	case node.GenesisReady:
		return color.RGBA{R: 98, G: 191, B: 67, A: 255}
	case node.Ready:
		return color.RGBA{R: 98, G: 191, B: 67, A: 255}
	case node.Offline:
		return color.RGBA{R: 230, G: 78, B: 18, A: 255}
	case node.NotSupported:
		return color.RGBA{R: 98, G: 186, B: 221, A: 255}
	case node.Undefined:
		return color.RGBA{R: 153, G: 102, B: 102, A: 255}
	}
	return color.RGBA{R: 153, G: 102, B: 102, A: 255}
}

func BuildImageOutput(target string, clusterOverview []ClusterNode, grid map[string]map[string]node.PeerInfo, outputTheme string) {
	log.Info("Drawing network according to the discovered map")

	baseXMargin := float64(4)

	ordOffsetX := baseXMargin
	nameOffsetX := ordOffsetX + 15
	addrOffsetX := nameOffsetX + 100
	versionOffsetX := addrOffsetX + 100
	snapshotOffsetX := versionOffsetX + 50
	latencyOffsetX := snapshotOffsetX + 50
	statusLocalOffsetX := latencyOffsetX + 120

	baseYMargin := float64(4)

	gridSize := float64(len(clusterOverview))

	iconWidth := 10
	iconHeight := 10
	iconMargin := 4

	textWidth := float64(8)
	textHeight := float64(12)

	canvasWidth := math.Max(float64(iconWidth+iconMargin)*gridSize+2*baseXMargin+2*textWidth, 550)

	canvasHeight := float64(iconHeight+iconMargin)*gridSize + 2*baseYMargin + textHeight + (gridSize+2)*textHeight

	c := canvas.New(canvasWidth, canvasHeight)

	ctx := canvas.NewContext(c)

	nodeIcon := canvas.RoundedRectangle(float64(iconWidth), float64(iconHeight), 1)

	var fontFamily *canvas.FontFamily
	fontFamily = canvas.NewFontFamily("Monospace")
	if err := fontFamily.LoadLocalFont("Monospace", canvas.FontRegular); err != nil {
		panic(err)
	}

	var backgroundColor = canvas.Transparent
	var textColor = canvas.Black

	ctx.SetFillColor(backgroundColor)

	if outputTheme == "light" {
		backgroundColor = canvas.White
		textColor = canvas.Black
	} else if outputTheme == "dark" {
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
	ctx.DrawText(nameOffsetX, textPosY, canvas.NewTextBox(textFace, "Id", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	ctx.DrawText(addrOffsetX, textPosY, canvas.NewTextBox(textFace, "Address", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	ctx.DrawText(versionOffsetX, textPosY, canvas.NewTextBox(textFace, "Version", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	ctx.DrawText(snapshotOffsetX, textPosY, canvas.NewTextBox(textFace, "Snapshot", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))

	ctx.DrawText(latencyOffsetX, textPosY, canvas.NewTextBox(textFace, "Latency", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	ctx.DrawText(statusLocalOffsetX, textPosY, canvas.NewTextBox(textFace, "Status Node", 0.0, 0.0, canvas.Right, canvas.Top, 0.0, 0.0))

	var lastTextPosY = float64(0)

	for i, nodeOverview := range clusterOverview {

		// segfault
		selfInfoState := node.Undefined
		if nodeOverview.SelfInfo != nil {
			selfInfoState = nodeOverview.SelfInfo.CardinalState()
		}

		statusTextFace1 := statusColorRGB(selfInfoState)

		offsetY := float64(i+1) * textHeight
		textPosY := canvasHeight - 1 - offsetY
		lastTextPosY = textPosY

		ctx.DrawText(ordOffsetX, textPosY, canvas.NewTextBox(textFace, fmt.Sprintf("%02d", i), 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))

		ctx.DrawText(nameOffsetX, textPosY, canvas.NewTextBox(textFace, nodeOverview.ShortId(), 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)) // TODO: Alias
		ctx.DrawText(addrOffsetX, textPosY, canvas.NewTextBox(textFace, nodeOverview.Addr.Ip, 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))

		ctx.DrawText(statusLocalOffsetX, textPosY, canvas.NewTextBox(fontFamily.Face(20.0, statusTextFace1, canvas.FontRegular, canvas.FontNormal), fmt.Sprintf("%s", selfInfoState), 0.0, 0.0, canvas.Right, canvas.Top, 0.0, 0.0))
	}

	lastTextPosY = lastTextPosY - textHeight*2

	for col, _ := range clusterOverview {
		tb := canvas.NewTextBox(textFace, fmt.Sprintf("%02d", col), 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)

		textPosY := lastTextPosY - 2
		textPosX := textWidth + float64((iconWidth+iconMargin)*col) + baseXMargin + textWidth

		ctx.DrawText(textPosX, textPosY, tb)
	}

	for row, rowNode := range clusterOverview {

		offsetY := colNumbersHeight + float64((iconHeight+iconMargin)*row)

		rowMap := grid[rowNode.Id]

		tb := canvas.NewTextBox(textFace, fmt.Sprintf("%02d", row), 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)

		textPosY := lastTextPosY - offsetY

		ctx.DrawText(baseXMargin, textPosY, tb)

		for col, colNode := range clusterOverview {

			offsetX := textWidth + float64((iconWidth+iconMargin)*col)

			cell, nocell := rowMap[colNode.Id]

			if nocell {
				ctx.SetFillColor(statusColorRGB(cell.CardinalState()))
				ctx.SetStrokeColor(statusColorRGB(cell.CardinalState()))
			} else {
				ctx.SetFillColor(statusColorRGB(node.Undefined))
				ctx.SetStrokeColor(statusColorRGB(node.Undefined))
			}

			iconPosY := lastTextPosY - (offsetY + float64(iconHeight))

			ctx.DrawPath(baseXMargin+offsetX+textWidth, iconPosY, nodeIcon.Copy())
		}
	}

	c.WriteFile(target, rasterizer.PNGWriter(2))
}
