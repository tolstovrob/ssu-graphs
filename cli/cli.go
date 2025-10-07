/*
 * WARNING! VIBE CODED! ALARM!
 */

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tolstovrob/graph-go/graph"
)

type CLIService struct {
	app         *tview.Application
	pages       *tview.Pages
	graph       *graph.Graph
	statusView  *tview.TextView
	currentView *tview.TextView
}

func NewCLIService() *CLIService {
	cli := &CLIService{
		app:   tview.NewApplication(),
		graph: graph.MakeGraph(),
	}

	cli.setupUI()
	return cli
}

func (cli *CLIService) Run() error {
	return cli.app.Run()
}

func (cli *CLIService) setupUI() {
	cli.pages = tview.NewPages()

	// Главное меню
	mainMenu := cli.createMainMenu()
	cli.pages.AddPage("main", mainMenu, true, true)

	// Статус бар
	cli.statusView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	cli.statusView.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	// Текущее представление графа с возможностью скролла
	cli.currentView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			cli.app.Draw()
		})
	cli.currentView.SetTitle(" Current Graph State ").SetBorder(true)

	// Основной layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(cli.statusView, 1, 0, false).
		AddItem(cli.pages, 0, 1, true).
		AddItem(cli.currentView, 10, 0, false)

	cli.app.SetRoot(flex, true)
	cli.updateStatus("Ready - Use arrow keys to navigate, Enter to select, Q to go back")
	cli.updateGraphView()
}

func (cli *CLIService) createMainMenu() tview.Primitive {
	list := tview.NewList().
		AddItem("Node Operations", "Add, remove, modify nodes", '1', cli.showNodeOperations).
		AddItem("Edge Operations", "Add, remove, modify edges", '2', cli.showEdgeOperations).
		AddItem("Graph Options", "Configure graph properties", '3', cli.showGraphOptions).
		AddItem("View Graph Info", "Display graph information", '4', cli.showGraphInfo).
		AddItem("JSON Operations", "Save/Load graph from JSON", '5', cli.showJSONOperations).
		AddItem("Quit", "Exit application", 'q', func() {
			cli.app.Stop()
		})

	list.SetBorder(true).SetTitle(" Graph CLI - Main Menu ")
	return list
}

func (cli *CLIService) showNodeOperations() {
	modal := tview.NewModal().
		SetText("Node Operations").
		AddButtons([]string{"Add Node", "Remove Node", "Modify Node", "List Nodes", "Back"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Add Node":
				cli.showAddNodeForm()
			case "Remove Node":
				cli.showRemoveNodeForm()
			case "Modify Node":
				cli.showModifyNodeForm()
			case "List Nodes":
				cli.showNodesList()
			case "Back":
				cli.pages.SwitchToPage("main")
			}
		})

	cli.pages.AddAndSwitchToPage("node_operations", modal, true)
}

func (cli *CLIService) showEdgeOperations() {
	modal := tview.NewModal().
		SetText("Edge Operations").
		AddButtons([]string{"Add Edge", "Remove Edge", "Modify Edge", "List Edges", "Back"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Add Edge":
				cli.showAddEdgeForm()
			case "Remove Edge":
				cli.showRemoveEdgeForm()
			case "Modify Edge":
				cli.showModifyEdgeForm()
			case "List Edges":
				cli.showEdgesList()
			case "Back":
				cli.pages.SwitchToPage("main")
			}
		})

	cli.pages.AddAndSwitchToPage("edge_operations", modal, true)
}

func (cli *CLIService) showGraphOptions() {
	form := tview.NewForm()

	// Добавляем чекбоксы с текущими значениями
	form.AddCheckbox("Directed Graph", cli.graph.Options.IsDirected, func(checked bool) {
		cli.graph.UpdateGraph(graph.WithGraphDirected(checked))
	})
	form.AddCheckbox("Multi Graph", cli.graph.Options.IsMulti, func(checked bool) {
		cli.graph.UpdateGraph(graph.WithGraphMulti(checked))
	})

	form.AddButton("Save", func() {
		cli.updateStatus("Graph options updated")
		cli.updateGraphView()
		cli.pages.SwitchToPage("main")
	})
	form.AddButton("Back", func() {
		cli.pages.SwitchToPage("main")
	})

	form.SetBorder(true).SetTitle(" Graph Options ")
	cli.pages.AddAndSwitchToPage("graph_options", form, true)
}

func (cli *CLIService) showGraphInfo() {
	info := cli.getDetailedGraphInfo()
	cli.showScrollableModal("Graph Information", info, "main")
}

func (cli *CLIService) showJSONOperations() {
	modal := tview.NewModal().
		SetText("JSON Operations").
		AddButtons([]string{"Save to JSON", "Load from JSON", "Show JSON", "Back"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Save to JSON":
				cli.showSaveJSONForm()
			case "Load from JSON":
				cli.showLoadJSONForm()
			case "Show JSON":
				cli.showJSONView()
			case "Back":
				cli.pages.SwitchToPage("main")
			}
		})

	cli.pages.AddAndSwitchToPage("json_operations", modal, true)
}

// Node Operations
func (cli *CLIService) showAddNodeForm() {
	form := tview.NewForm()
	var key, label string

	form.AddInputField("Key", "", 10, nil, func(text string) {
		key = text
	})
	form.AddInputField("Label", "", 20, nil, func(text string) {
		label = text
	})
	form.AddButton("Add", func() {
		keyVal, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format")
			return
		}

		node := graph.MakeNode(graph.TKey(keyVal))
		if label != "" {
			node.UpdateNode(graph.WithNodeLabel(label))
		}

		if err := cli.graph.AddNode(node); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err))
		} else {
			cli.updateStatus(fmt.Sprintf("Node %d added successfully", keyVal))
			cli.updateGraphView()
			cli.pages.SwitchToPage("main")
		}
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("node_operations")
	})

	form.SetBorder(true).SetTitle(" Add Node ")
	cli.pages.AddAndSwitchToPage("add_node", form, true)
}

func (cli *CLIService) showRemoveNodeForm() {
	form := tview.NewForm()
	var key string

	form.AddInputField("Node Key", "", 10, nil, func(text string) {
		key = text
	})
	form.AddButton("Remove", func() {
		keyVal, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format")
			return
		}

		if err := cli.graph.RemoveNodeByKey(graph.TKey(keyVal)); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err))
		} else {
			cli.updateStatus(fmt.Sprintf("Node %d removed successfully", keyVal))
			cli.updateGraphView()
			cli.pages.SwitchToPage("main")
		}
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("node_operations")
	})

	form.SetBorder(true).SetTitle(" Remove Node ")
	cli.pages.AddAndSwitchToPage("remove_node", form, true)
}

func (cli *CLIService) showModifyNodeForm() {
	form := tview.NewForm()
	var key, newLabel string

	form.AddInputField("Node Key", "", 10, nil, func(text string) {
		key = text
	})
	form.AddInputField("New Label", "", 20, nil, func(text string) {
		newLabel = text
	})
	form.AddButton("Modify", func() {
		keyVal, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format")
			return
		}

		node, err := cli.graph.GetNodeByKey(graph.TKey(keyVal))
		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err))
			return
		}

		node.UpdateNode(graph.WithNodeLabel(newLabel))
		cli.updateStatus(fmt.Sprintf("Node %d modified successfully", keyVal))
		cli.updateGraphView()
		cli.pages.SwitchToPage("main")
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("node_operations")
	})

	form.SetBorder(true).SetTitle(" Modify Node ")
	cli.pages.AddAndSwitchToPage("modify_node", form, true)
}

func (cli *CLIService) showNodesList() {
	nodesInfo := "Nodes:\n\n"
	for key, node := range cli.graph.Nodes {
		nodesInfo += fmt.Sprintf("Key: %d, Label: %s\n", key, node.Label)
	}

	cli.showScrollableModal("Nodes List", nodesInfo, "node_operations")
}

// Edge Operations
func (cli *CLIService) showAddEdgeForm() {
	form := tview.NewForm()
	var edgeKey, srcKey, dstKey, weightStr, label string

	form.AddInputField("Edge Key", "", 10, nil, func(text string) {
		edgeKey = text
	})
	form.AddInputField("Source Node Key", "", 10, nil, func(text string) {
		srcKey = text
	})
	form.AddInputField("Destination Node Key", "", 10, nil, func(text string) {
		dstKey = text
	})
	form.AddInputField("Weight", "0", 10, nil, func(text string) {
		weightStr = text
	})
	form.AddInputField("Label", "", 20, nil, func(text string) {
		label = text
	})
	form.AddButton("Add", func() {
		key, err := strconv.ParseUint(edgeKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid edge key format")
			return
		}

		src, err := strconv.ParseUint(srcKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid source key format")
			return
		}

		dst, err := strconv.ParseUint(dstKey, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid destination key format")
			return
		}

		weight, err := strconv.ParseUint(weightStr, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid weight format")
			return
		}

		edge := graph.MakeEdge(graph.TKey(key), graph.TKey(src), graph.TKey(dst))
		if weight > 0 {
			edge.UpdateEdge(graph.WithEdgeWeight(graph.TWeight(weight)))
		}
		if label != "" {
			edge.UpdateEdge(graph.WithEdgeLabel(label))
		}

		if err := cli.graph.AddEdge(edge); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err))
		} else {
			cli.updateStatus(fmt.Sprintf("Edge %d added successfully", key))
			cli.updateGraphView()
			cli.pages.SwitchToPage("main")
		}
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("edge_operations")
	})

	form.SetBorder(true).SetTitle(" Add Edge ")
	cli.pages.AddAndSwitchToPage("add_edge", form, true)
}

func (cli *CLIService) showRemoveEdgeForm() {
	form := tview.NewForm()
	var key string

	form.AddInputField("Edge Key", "", 10, nil, func(text string) {
		key = text
	})
	form.AddButton("Remove", func() {
		keyVal, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format")
			return
		}

		if err := cli.graph.RemoveEdgeByKey(graph.TKey(keyVal)); err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err))
		} else {
			cli.updateStatus(fmt.Sprintf("Edge %d removed successfully", keyVal))
			cli.updateGraphView()
			cli.pages.SwitchToPage("main")
		}
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("edge_operations")
	})

	form.SetBorder(true).SetTitle(" Remove Edge ")
	cli.pages.AddAndSwitchToPage("remove_edge", form, true)
}

func (cli *CLIService) showModifyEdgeForm() {
	form := tview.NewForm()
	var key, weightStr, label string

	form.AddInputField("Edge Key", "", 10, nil, func(text string) {
		key = text
	})
	form.AddInputField("New Weight", "", 10, nil, func(text string) {
		weightStr = text
	})
	form.AddInputField("New Label", "", 20, nil, func(text string) {
		label = text
	})
	form.AddButton("Modify", func() {
		keyVal, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			cli.updateStatus("Error: Invalid key format")
			return
		}

		edge, err := cli.graph.GetEdgeByKey(graph.TKey(keyVal))
		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error: %v", err))
			return
		}

		if weightStr != "" {
			weight, err := strconv.ParseUint(weightStr, 10, 64)
			if err != nil {
				cli.updateStatus("Error: Invalid weight format")
				return
			}
			edge.UpdateEdge(graph.WithEdgeWeight(graph.TWeight(weight)))
		}

		if label != "" {
			edge.UpdateEdge(graph.WithEdgeLabel(label))
		}

		cli.updateStatus(fmt.Sprintf("Edge %d modified successfully", keyVal))
		cli.updateGraphView()
		cli.pages.SwitchToPage("main")
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("edge_operations")
	})

	form.SetBorder(true).SetTitle(" Modify Edge ")
	cli.pages.AddAndSwitchToPage("modify_edge", form, true)
}

func (cli *CLIService) showEdgesList() {
	edgesInfo := "Edges:\n\n"
	for key, edge := range cli.graph.Edges {
		edgesInfo += fmt.Sprintf("Key: %d, Source: %d -> Destination: %d, Weight: %d, Label: %s\n",
			key, edge.Source, edge.Destination, edge.Weight, edge.Label)
	}

	cli.showScrollableModal("Edges List", edgesInfo, "edge_operations")
}

// JSON Operations
func (cli *CLIService) showSaveJSONForm() {
	form := tview.NewForm()
	var filename string

	form.AddInputField("Filename", "graph.json", 30, nil, func(text string) {
		filename = text
	})
	form.AddButton("Save", func() {
		if filename == "" {
			cli.updateStatus("Error: Filename cannot be empty")
			return
		}

		jsonData, err := cli.graph.ToJSON()
		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error generating JSON: %v", err))
			return
		}

		err = os.WriteFile(filename, []byte(jsonData), 0644)
		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error writing file: %v", err))
			return
		}

		cli.updateStatus(fmt.Sprintf("Graph saved to %s successfully", filename))
		cli.pages.SwitchToPage("main")
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("json_operations")
	})

	form.SetBorder(true).SetTitle(" Save Graph to JSON ")
	cli.pages.AddAndSwitchToPage("save_json", form, true)
}

func (cli *CLIService) showLoadJSONForm() {
	form := tview.NewForm()
	var filename string

	form.AddInputField("Filename", "graph.json", 30, nil, func(text string) {
		filename = text
	})
	form.AddButton("Load", func() {
		if filename == "" {
			cli.updateStatus("Error: Filename cannot be empty")
			return
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			cli.updateStatus(fmt.Sprintf("Error reading file: %v", err))
			return
		}

		newGraph := graph.MakeGraph()
		if err := newGraph.FromJSON(string(data)); err != nil {
			cli.updateStatus(fmt.Sprintf("Error parsing JSON: %v", err))
			return
		}

		cli.graph = newGraph
		cli.updateStatus(fmt.Sprintf("Graph loaded from %s successfully", filename))
		cli.updateGraphView()
		cli.pages.SwitchToPage("main")
	})
	form.AddButton("Cancel", func() {
		cli.pages.SwitchToPage("json_operations")
	})

	form.SetBorder(true).SetTitle(" Load Graph from JSON ")
	cli.pages.AddAndSwitchToPage("load_json", form, true)
}

func (cli *CLIService) showJSONView() {
	jsonData, err := cli.graph.ToJSON()
	if err != nil {
		cli.updateStatus(fmt.Sprintf("Error generating JSON: %v", err))
		return
	}

	// Форматируем JSON для лучшего отображения
	var formattedJSON map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &formattedJSON); err != nil {
		cli.updateStatus(fmt.Sprintf("Error formatting JSON: %v", err))
		return
	}

	prettyJSON, err := json.MarshalIndent(formattedJSON, "", "  ")
	if err != nil {
		cli.updateStatus(fmt.Sprintf("Error formatting JSON: %v", err))
		return
	}

	// Создаем текстовое view с возможностью скролла
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			cli.app.Draw()
		}).
		SetText(string(prettyJSON))

	textView.SetBorder(true).SetTitle(" Graph JSON - Use arrow keys to scroll, Q to go back ")

	// Добавляем кнопки навигации
	buttons := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewButton("Scroll Up").SetSelectedFunc(func() {
			row, col := textView.GetScrollOffset()
			textView.ScrollTo(row-5, col)
		}), 0, 1, false).
		AddItem(tview.NewButton("Scroll Down").SetSelectedFunc(func() {
			row, col := textView.GetScrollOffset()
			textView.ScrollTo(row+5, col)
		}), 0, 1, false).
		AddItem(tview.NewButton("Top").SetSelectedFunc(func() {
			textView.ScrollToBeginning()
		}), 0, 1, false).
		AddItem(tview.NewButton("Bottom").SetSelectedFunc(func() {
			textView.ScrollToEnd()
		}), 0, 1, false).
		AddItem(tview.NewButton("Back (Q)").SetSelectedFunc(func() {
			cli.pages.SwitchToPage("json_operations")
		}), 0, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true).
		AddItem(buttons, 1, 0, false)

	cli.pages.AddAndSwitchToPage("json_view", flex, true)
}

// Универсальный метод для создания скроллируемых модальных окон с поддержкой клавиши Q
func (cli *CLIService) showScrollableModal(title, content, returnPage string) {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			cli.app.Draw()
		}).
		SetText(content)

	textView.SetBorder(true).SetTitle(fmt.Sprintf(" %s - Use arrow keys to scroll, Q to go back ", title))

	// Добавляем кнопки навигации
	buttons := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewButton("Scroll Up").SetSelectedFunc(func() {
			row, col := textView.GetScrollOffset()
			textView.ScrollTo(row-5, col)
		}), 0, 1, false).
		AddItem(tview.NewButton("Scroll Down").SetSelectedFunc(func() {
			row, col := textView.GetScrollOffset()
			textView.ScrollTo(row+5, col)
		}), 0, 1, false).
		AddItem(tview.NewButton("Top").SetSelectedFunc(func() {
			textView.ScrollToBeginning()
		}), 0, 1, false).
		AddItem(tview.NewButton("Bottom").SetSelectedFunc(func() {
			textView.ScrollToEnd()
		}), 0, 1, false).
		AddItem(tview.NewButton("Back (Q)").SetSelectedFunc(func() {
			cli.pages.SwitchToPage(returnPage)
		}), 0, 1, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true).
		AddItem(buttons, 1, 0, false)

	// Обработка клавиши Q для всего flex контейнера
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'Q' {
			cli.pages.SwitchToPage(returnPage)
			return nil
		}
		return event
	})

	cli.pages.AddAndSwitchToPage(strings.ToLower(strings.ReplaceAll(title, " ", "_"))+"_view", flex, true)
}

// Детальное представление графа с полной информацией
func (cli *CLIService) getDetailedGraphInfo() string {
	var info strings.Builder

	// 1. ОПЦИИ ГРАФА (в начале)
	info.WriteString("GRAPH OPTIONS\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	info.WriteString(fmt.Sprintf("Directed: %v\n", cli.graph.Options.IsDirected))
	info.WriteString(fmt.Sprintf("Multi-graph: %v\n", cli.graph.Options.IsMulti))
	info.WriteString(fmt.Sprintf("Graph Type: %s%s\n\n",
		map[bool]string{true: "Directed", false: "Undirected"}[cli.graph.Options.IsDirected],
		map[bool]string{true: " Multi", false: ""}[cli.graph.Options.IsMulti]))

	// 2. ОСНОВНАЯ СТАТИСТИКА
	info.WriteString("BASIC STATISTICS\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	info.WriteString(fmt.Sprintf("Total Nodes: %d\n", len(cli.graph.Nodes)))
	info.WriteString(fmt.Sprintf("Total Edges: %d\n", len(cli.graph.Edges)))
	info.WriteString(fmt.Sprintf("Density: %.4f\n\n", cli.calculateGraphDensity()))

	// 3. СТАТИСТИКА ПО СТЕПЕНЯМ ВЕРШИН
	degreeDistribution := make(map[int]int)
	maxDegree := 0
	minDegree := -1
	totalDegree := 0
	isolatedNodes := 0

	for _, neighbors := range cli.graph.AdjacencyMap {
		degree := len(neighbors)
		degreeDistribution[degree]++
		totalDegree += degree
		if degree > maxDegree {
			maxDegree = degree
		}
		if minDegree == -1 || degree < minDegree {
			minDegree = degree
		}
		if degree == 0 {
			isolatedNodes++
		}
	}

	avgDegree := 0.0
	if len(cli.graph.AdjacencyMap) > 0 {
		avgDegree = float64(totalDegree) / float64(len(cli.graph.AdjacencyMap))
	}

	info.WriteString("DEGREE STATISTICS\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	info.WriteString(fmt.Sprintf("Minimum degree: %d\n", minDegree))
	info.WriteString(fmt.Sprintf("Maximum degree: %d\n", maxDegree))
	info.WriteString(fmt.Sprintf("Average degree: %.2f\n", avgDegree))
	info.WriteString(fmt.Sprintf("Isolated nodes: %d\n", isolatedNodes))

	// Распределение степеней
	info.WriteString("\nDegree Distribution:\n")
	degrees := make([]int, 0, len(degreeDistribution))
	for degree := range degreeDistribution {
		degrees = append(degrees, degree)
	}
	sort.Ints(degrees)

	for _, degree := range degrees {
		count := degreeDistribution[degree]
		percentage := float64(count) / float64(len(cli.graph.Nodes)) * 100
		info.WriteString(fmt.Sprintf("  Degree %d: %d nodes (%.1f%%)\n", degree, count, percentage))
	}
	info.WriteString("\n")

	// 4. ПОЛНЫЙ СПИСОК УЗЛОВ
	info.WriteString("COMPLETE NODES LIST\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	if len(cli.graph.Nodes) == 0 {
		info.WriteString("No nodes in graph\n")
	} else {
		// Сортируем узлы по ключу
		keys := make([]graph.TKey, 0, len(cli.graph.Nodes))
		for key := range cli.graph.Nodes {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		for _, key := range keys {
			node := cli.graph.Nodes[key]
			degree := len(cli.graph.AdjacencyMap[key])
			info.WriteString(fmt.Sprintf("Key: %4d | Label: %-20s | Degree: %d\n",
				key, node.Label, degree))
		}
	}
	info.WriteString("\n")

	// 5. ПОЛНЫЙ СПИСОК РЕБЕР
	info.WriteString("COMPLETE EDGES LIST\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	if len(cli.graph.Edges) == 0 {
		info.WriteString("No edges in graph\n")
	} else {
		// Сортируем ребра по ключу
		keys := make([]graph.TKey, 0, len(cli.graph.Edges))
		for key := range cli.graph.Edges {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		for _, key := range keys {
			edge := cli.graph.Edges[key]
			info.WriteString(fmt.Sprintf("Key: %4d | %4d → %4d | Weight: %4d | Label: %s\n",
				key, edge.Source, edge.Destination, edge.Weight, edge.Label))
		}
	}
	info.WriteString("\n")

	// 6. ПОЛНЫЙ ADJACENCY LIST
	info.WriteString("COMPLETE ADJACENCY LIST\n")
	info.WriteString(strings.Repeat("─", 50) + "\n")
	if len(cli.graph.AdjacencyMap) == 0 {
		info.WriteString("Empty adjacency list\n")
	} else {
		// Сортируем узлы adjacency list по ключу
		keys := make([]graph.TKey, 0, len(cli.graph.AdjacencyMap))
		for key := range cli.graph.AdjacencyMap {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		for _, key := range keys {
			neighbors := cli.graph.AdjacencyMap[key]
			info.WriteString(fmt.Sprintf("%4d → [", key))

			// Сортируем соседей
			sort.Slice(neighbors, func(i, j int) bool { return neighbors[i] < neighbors[j] })

			for i, neighbor := range neighbors {
				if i > 0 {
					info.WriteString(", ")
				}
				info.WriteString(fmt.Sprintf("%d", neighbor))
			}
			info.WriteString(fmt.Sprintf("] (degree: %d)\n", len(neighbors)))
		}
	}

	return info.String()
}

// Вспомогательные методы для статистики
func (cli *CLIService) calculateGraphDensity() float64 {
	n := len(cli.graph.Nodes)
	if n <= 1 {
		return 0.0
	}

	maxEdges := n * (n - 1)
	if !cli.graph.Options.IsDirected {
		maxEdges = maxEdges / 2
	}

	if maxEdges == 0 {
		return 0.0
	}

	return float64(len(cli.graph.Edges)) / float64(maxEdges)
}

// Вспомогательные методы
func (cli *CLIService) updateStatus(message string) {
	cli.statusView.SetText(fmt.Sprintf("[green]%s", message))
}

func (cli *CLIService) updateGraphView() {
	info := cli.getCompactGraphInfo()
	cli.currentView.SetText(info)
	cli.currentView.ScrollToBeginning()
}

func (cli *CLIService) getCompactGraphInfo() string {
	var info strings.Builder

	info.WriteString(fmt.Sprintf("Graph: %s%s | ",
		map[bool]string{true: "Directed", false: "Undirected"}[cli.graph.Options.IsDirected],
		map[bool]string{true: " Multi", false: ""}[cli.graph.Options.IsMulti]))

	info.WriteString(fmt.Sprintf("Nodes: %d | Edges: %d\n\n", len(cli.graph.Nodes), len(cli.graph.Edges)))

	if len(cli.graph.AdjacencyMap) > 0 {
		info.WriteString("Adjacency Preview:\n")
		count := 0
		for node, neighbors := range cli.graph.AdjacencyMap {
			if count >= 8 {
				info.WriteString(fmt.Sprintf("  ... and %d more nodes\n", len(cli.graph.AdjacencyMap)-8))
				break
			}
			info.WriteString(fmt.Sprintf("  %d -> ", node))
			for i, neighbor := range neighbors {
				if i >= 3 {
					info.WriteString(fmt.Sprintf("... (%d more)", len(neighbors)-3))
					break
				}
				if i > 0 {
					info.WriteString(", ")
				}
				info.WriteString(fmt.Sprintf("%d", neighbor))
			}
			info.WriteString("\n")
			count++
		}
	} else {
		info.WriteString("Graph is empty\n")
	}

	return info.String()
}
