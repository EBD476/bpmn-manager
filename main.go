package main

import (
	"fmt"
	"os"
	"time"

	"bpmn-manager/api"
	"bpmn-manager/models"

	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// -----------------------------------------------------------------------
type BPMNManager struct {
	app         *tview.Application
	pages       *tview.Pages
	apiClient   *api.APIClient
	infoPanel   *tview.TextView
	mainContent *tview.Flex
	nav         *tview.List
	currentPage string
	// contentView *tview.TextView // Add the contentView field here

	baseURL string
}

// -----------------------------------------------------------------------
func NewBPMNManager(baseURL string) *BPMNManager {
	manager := &BPMNManager{
		app:       tview.NewApplication(),
		pages:     tview.NewPages(),
		apiClient: api.NewAPIClient(baseURL),
		baseURL:   baseURL,
	}
	// Set up proper encoding for Persian/Arabic text
	manager.setupEncoding()

	return manager

}

// -----------------------------------------------------------------------
func (m *BPMNManager) setupEncoding() {
	// Set environment variables for proper Unicode support
	os.Setenv("LANG", "fa_IR.utf8")
	os.Setenv("LC_ALL", "fa_IR.utf8")
	os.Setenv("LC_MESSAGES", "fa_IR.utf8")
}

// -----------------------------------------------------------------------
// Helper function to check if text contains Persian/Arabic characters
func containsPersian(text string) bool {
	for _, r := range text {
		// Persian/Arabic character ranges
		if (r >= 0x0600 && r <= 0x06FF) || // Arabic
			(r >= 0x0750 && r <= 0x077F) || // Arabic Supplement
			(r >= 0x08A0 && r <= 0x08FF) || // Arabic Extended-A
			(r >= 0xFB50 && r <= 0xFDFF) || // Arabic Presentation Forms-A
			(r >= 0xFE70 && r <= 0xFEFF) { // Arabic Presentation Forms-B
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------------
func (m *BPMNManager) setupUI() {

	m.app.EnableMouse(true)

	// Create main menu
	menu := m.createMainMenu()
	m.pages.AddPage("main", menu, true, true)
	m.app.SetRoot(m.pages, true)
}

// -----------------------------------------------------------------------
func (m *BPMNManager) createMainMenu() tview.Primitive {
	// Create main flex layout
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Header
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("🏭 BPMN Activity Manager - %s", m.baseURL)).
		// SetTextColor(tcell.ColorYellow)
		SetTextColor(tcell.ColorBeige)
	header.SetBorder(true).SetTitle(" BPMN Manager ")
	// header.SetBackgroundColor(tcell.ColoDarkSlateBlue)
	header.SetBackgroundColor(tcell.Color142)

	// Content area
	m.mainContent = tview.NewFlex().SetDirection(tview.FlexColumn)

	// Navigation panel
	m.nav = m.createNavigationPanel()

	// 	welcomeText := `🎯 به سیستم مدیریت BPMN خوش آمدید!

	// از منوی ناوبری برای موارد زیر استفاده کنید:
	// • مشاهده وظایف محول شده
	// • نظارت بر فرآیندهای در حال اجرا
	// • مشاهده جزئیات فرآیند

	// برای بروزرسانی اطلاعات کلید F5 را فشار دهید`

	// Main content
	// content := tview.NewFlex().SetDirection(tview.FlexColumn)
	// leftPanel := m.createUserTasksPanel()

	// Main content view
	m.infoPanel = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)
	// contentView.SetText(m.formatRTLText(welcomeText))
	m.infoPanel.SetBorder(true).SetTitle(" Dashboard ").SetBorderColor(tcell.Color102)

	m.updateDashboardPanel()

	// m.mainContent.AddItem(m.nav, 35, 1, true)
	// m.mainContent.AddItem(m.infoPanel, 0, 3, false)

	// Footer
	footer := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Press F5 to refresh • Tab to navigate • Ctrl+C to exit")
	footer.SetBorder(false)

	flex.AddItem(header, 3, 1, false)
	flex.AddItem(m.mainContent, 0, 8, true)
	flex.AddItem(footer, 1, 1, false)

	return flex
}

// -----------------------------------------------------------------------

func (m *BPMNManager) createNavigationPanel() *tview.List {
	nav := tview.NewList().
		AddItem("👤 My Tasks", "View assigned tasks", 't', func() {
			// m.showUserTasks()
			m.showDashboard()
		}).
		AddItem("🔄 Running Processes", "View active processes", 'r', func() {
			// m.showRunningProcesses()
			m.createRunningProcesses()
		}).
		AddItem("📊 Completed tasks", "View completed tasks", 'c', func() {
			m.showCompletedTaskDetails()
		}).
		AddItem("🔎 Find Process ", "Find process by id", 'f', func() {
			m.showProcessSearch()
		}).
		AddItem("🚀 Start Process Instance ", "Launch process instance", 'l', func() {
			m.apiClient.StartProcess()
			m.showMessage("New process instance started!")
		}).
		AddItem("📊 Process Details", "View process information", 'd', func() {
			m.showProcessSelection()
		}).
		AddItem("🔄 Refresh Data", "Reload all data", 'f', func() {
			m.updateDashboardPanel()
		}).
		AddItem("⚙️ Settings", "Configure connection", 's', func() {
			m.showSettings()
		}).
		AddItem("❌ Quit", "Exit application", 'q', func() {
			m.app.Stop()
		})
	nav.SetBorder(true).SetTitle(" Navigation ").SetBorderColor(tcell.Color102)
	selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)
	nav.SetSelectedStyle(selectedStyle)

	return nav
}

// -----------------------------------------------------------------------

func (m *BPMNManager) updateDashboardPanel() {

	tasks, _ := m.apiClient.GetUserTasks()
	completedTasks, _ := m.apiClient.GetCompletedTasks()
	processes, _ := m.apiClient.GetRunningProcesses()
	completedProcesses, _ := m.apiClient.GetCompletedProcesses()

	detailsText := fmt.Sprintf(`

 🎯 Welcome to BPMN Manager!
  [yellow]---------------------------------	
 [yellow] 📊  Total Running Tasks: %d   
  📊  Total Completed Tasks: %d 
  ---------------------------------
[yellow]  📊  Total Running Processes: %d   
  📊  Total Completed Processes: %d 
 ----------------------------------[white] 
 
  Use the navigation menu to:
  • View your assigned tasks
  • View completed tasks
  • Monitor running processes
  • Check process details
  
  Press F5 to refresh data`, len(tasks), len(completedTasks), len(processes), len(completedProcesses))

	m.infoPanel.SetText(detailsText)
	m.mainContent.Clear()
	m.mainContent.AddItem(m.nav, 35, 1, true)
	m.mainContent.AddItem(m.infoPanel, 0, 3, false)
	m.app.SetFocus(m.nav)
}

// -----------------------------------------------------------------------
func (m *BPMNManager) showDashboard() {
	// Create a simple dashboard view
	modal := tview.NewModal().
		SetText("🔄 Loading dashboard data...").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			m.pages.SwitchToPage("main")
		})
	m.pages.AddPage("loading_dashboard", modal, true, true)
	m.currentPage = "task_page"
	// Fetch data in background
	go func() {
		// time.Sleep(500 * time.Microsecond)
		// Fetch user tasks
		tasks, err := m.apiClient.GetUserTasks()
		if err != nil {
			m.app.QueueUpdateDraw(func() {
				m.showError("Failed to load user tasks: " + err.Error())
			})
			return
		}

		m.app.QueueUpdateDraw(func() {
			// time.Sleep(1 * time.Second)
			m.pages.SwitchToPage("main")
			m.updateDashboard(tasks)
		})
	}()

	// Simulate API call with goroutine
	// go func() {
	// 	time.Sleep(1 * time.Second) // Simulate loading
	// 	m.app.QueueUpdateDraw(func() {
	// 		m.pages.SwitchToPage("main")
	// 	})
	// }()
}

// -----------------------------------------------------------------------
func (m *BPMNManager) showError(message string) {
	modal := tview.NewModal().
		SetText("❌ " + message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			m.pages.SwitchToPage("main")
		})
	m.pages.AddAndSwitchToPage("error", modal, true)
}

// -----------------------------------------------------------------------
func (m *BPMNManager) createUserTasksPanel(tasks []models.UserTask) *tview.List {

	list := tview.NewList()
	list.SetBorder(true).SetTitle(" Process Groups ")

	for _, task := range tasks {

		fmt.Fprintf(m.infoPanel, "  %s -  %s  -  %s \n",
			task.ID, reverseString(task.Name), task.Assignee)

		list.AddItem(task.Name, task.ID, 0, func() {

		})
	}
	return list
}

// -----------------------------------------------------------------------
func (m *BPMNManager) createTaskDetails(tasks []models.UserTask) *tview.TextView {

	taskDetails := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	completedTasks, _ := m.apiClient.GetCompletedTasks()

	detailsText := fmt.Sprintf(`Task Summary:
[yellow:darkgreen]    📊  Total Completed Tasks: %d    
	📊  Total Available Tasks: %d    `, len(completedTasks), len(tasks))
	// taskDetails.SetText(fmt.Sprintf("\n[yellow:darkgreen]📊 Total completed tasks : %d", len(completedTasks)))

	// 	detailsText := fmt.Sprintf(`Process: Order Processing
	// ID: %s
	// Description: Handles customer order fulfillment
	// Version: 2.1
	// Status: Active

	// Activities:
	// 1. Receive Order
	// 2. Validate Payment
	// 3. Prepare Shipment
	// 4. Deliver Order

	// Statistics:
	// • Total Instances: 156
	// • Running: 12
	// • Completed: 144
	// • Success Rate: 92%%`, "123")

	// taskDetails.SetText(fmt.Sprintf("\n[yellow:darkgreen]📊 Total tasks : %d - %s[white]\n\n", len(tasks), time.Now().Format("2006-01-02 15:04:05")))

	taskDetails.SetBorder(true).SetBorderColor(tcell.Color102)
	taskDetails.SetTitle("Task Details")
	taskDetails.SetText(detailsText)

	return taskDetails
}

// -----------------------------------------------------------------------
func (m *BPMNManager) createProcessDetailsPanel(processes []models.RunningProcess) *tview.TextView {

	taskDetails := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	completedTasks, _ := m.apiClient.GetCompletedProcesses()

	detailsText := fmt.Sprintf(`Task Summary:
[yellow]  📊  Total Running Processes: %d   
  📊  Total Completed Processes: %d [white] `, len(processes), len(completedTasks))

	taskDetails.SetBorder(true).SetBorderColor(tcell.Color102)
	taskDetails.SetText(detailsText)
	taskDetails.SetTitle("Process Details")

	return taskDetails
}

// -----------------------------------------------------------------------
// CreateModalForTaskCompletion function
func (m *BPMNManager) CreateModalForTaskCompletion(taskId string) tview.Primitive {
	// Create a new modal with a title and content
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Complete Task: %s", taskId)).
		AddButtons([]string{"Complete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {

			switch buttonLabel {
			case "Complete":

				data := models.FormData{
					ReDefineDecision: false,
					DbDecision:       false,
					Comment:          "comment",
				}

				err := m.apiClient.CompleteTask(taskId, data)
				if err != nil {
					m.infoPanel.SetText(err.Error())
				} else {
					m.infoPanel.SetText("Task successfully done." + taskId)
					m.showDashboard()
				}
				m.pages.SwitchToPage("main")
				return
			case "Cancel":
				// fmt.Println("Task completion canceled.")
				m.pages.SwitchToPage("main")
				return
			}
		})

	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If the Escape key is pressed, stop the app (close the modal)
		if event.Key() == tcell.KeyEsc {
			m.pages.SwitchToPage("main")
		}
		return event
	})

	// Simulate API call
	go func() {

		m.app.QueueUpdateDraw(func() {
			m.pages.AddPage("user_task_complete", modal, true, true)
		})
	}()

	// // Use a flex layout to arrange the modal and form
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(modal, 0, 1, true)
	// 	AddItem(form, 0, 1, false)

	m.pages.AddPage("user_task_complete", modal, true, true)

	return flex
}

// -----------------------------------------------------------------------
func (m *BPMNManager) createProcessDetails(processId string) string {

	process, err := m.apiClient.GetProcessDetails(processId)

	if err != nil {
		m.infoPanel.SetText(err.Error())
		return err.Error()
	}

	if process == nil {
		return ""
	}

	detailsText := fmt.Sprintf(`Process:[white]	
	ID: %s
	[DarkCyan:green]Total completed activities:  %d [white:darkgreen]
	ProcessDefinitionId: %s
	ProcessDefinitionKey: %s
	StartTime: %s
	EndtTime: %s
	Duration: %s
`,
		process.ID,
		len(process.Activities),
		process.ProcessDefinitionId,
		process.ProcessDefinitionKey,
		process.StartTime,
		process.EndTime,
		formatDuration(process.Duration/1000),
	)

	detailsText += "[violet]CurrentVariables:\n"
	for key, variable := range process.CurrentVariables {
		switch v := variable.(type) {
		case string:
			detailsText += fmt.Sprintf("	[violet]%s:[violet] %s\n", key, v)
		case bool:
			detailsText += fmt.Sprintf("	[violet]%s:[violet] %t\n", key, v)

		}
	}

	detailsText += "[orange]Activities:"

	for _, activity := range process.Activities {

		detailsText += fmt.Sprintf(`[orange]
	[orange:gray]TaskId: %s	[orange:darkgreen] 
	ActivityId: %s
	ActivityName: %s
	[orange:green]ActivityType: %s [orange:darkgreen] 
	Assignee: %s
	StartTime: %s
	EndTime: %s
	Duration: %s
	--------------------------------------`,
			activity.TaskId,
			activity.ID,
			reverseString(activity.Name),
			activity.Type,
			activity.Assignee,
			activity.StartTime.Format("2006-01-02 15:04:05"),
			formatEndTime(activity.EndTime), //.Format("2025-01-02 15:04:05"),
			formatDuration(activity.Duration/1000),
		)
	}

	return detailsText

}

// -----------------------------------------------------------------------
// FormatDuration takes a duration in seconds and returns it in "XXh XXm XXs" format
func formatDuration(durationInSeconds int) string {
	// Convert seconds to time.Duration
	duration := time.Duration(durationInSeconds) * time.Second

	// Get the hours, minutes, and seconds from the duration
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	// Return the formatted duration string
	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}

// -----------------------------------------------------------------------
// formatEndTime formats the EndTime if it's not the zero value
func formatEndTime(t time.Time) string {
	if t.IsZero() {
		// Handle the zero value case (process is still ongoing or no end time is set)
		return "Not Finished" // You can replace this with an empty string or another placeholder
	}
	// Return the formatted EndTime if valid
	return t.Format("2006-01-02 15:04:05")
}

// -----------------------------------------------------------------------
func (m *BPMNManager) updateDashboard(tasks []models.UserTask) {

	if len(tasks) < 1 {
		m.infoPanel.Clear()
		fmt.Fprintf(m.infoPanel, "⚠️ No data available ...")
		return
	}

	for i, task := range tasks {
		if i >= 4 { // Show only first 5 tasks
			fmt.Fprintf(m.infoPanel, "  ... and %d more tasks\n", len(tasks)-5)
			break
		}
		fmt.Fprintf(m.infoPanel, "  %s -  %s  -  %s \n",
			task.ID, reverseString(task.Name), task.Assignee)
	}

	m.mainContent.Clear()
	m.mainContent.AddItem(m.nav, 35, 1, true)

	leftPanel := m.createUserTasksTable()
	m.mainContent.AddItem(leftPanel, 70, 1, true)
	m.infoPanel = m.createTaskDetails(tasks)
	m.mainContent.AddItem(m.infoPanel, 0, 1, true)
	m.app.SetFocus(leftPanel)
}

// -----------------------------------------------------------------------
func reverseString(s string) string {
	// Convert string to slice of runes to handle multi-byte characters correctly
	runes := []rune(s)
	// Reverse the slice of runes
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	// Convert the reversed rune slice back to a string and return it
	return string(runes)
}

// -----------------------------------------------------------------------
func (m *BPMNManager) showCompletedTaskDetails() {

	// Create tasks table
	table := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(true, false).SetBordersColor(tcell.Color100)

	selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)
	table.SetSelectedStyle(selectedStyle)

	// Headers
	headers := []string{"TaskID|", "TaskName", "|ProcessID", "|Assignee"}
	for i, header := range headers {
		table.SetCell(0, i,
			tview.NewTableCell(header).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}

	completedTasks, _ := m.apiClient.GetCompletedTasks()
	// Add data to table
	for row, task := range completedTasks {
		table.SetCell(row+1, 0, tview.NewTableCell(task.ID+" |"))
		table.SetCell(row+1, 1, tview.NewTableCell(reverseString(task.Name)).SetAlign(tview.AlignRight))
		table.SetCell(row+1, 2, tview.NewTableCell("|"+task.ProcessID))
		statusCell := tview.NewTableCell("|" + task.Assignee)
		statusCell.SetTextColor(tcell.ColorYellow)
		table.SetCell(row+1, 3, statusCell)
	}

	table.SetBorder(true).SetBorderColor(tcell.Color102)
	table.SetTitle(" Completed Tasks List")

	m.mainContent.Clear()
	m.mainContent.AddItem(m.nav, 35, 1, true)
	m.mainContent.AddItem(table, 0, 3, true)
	m.app.SetFocus(table)

	// leftPanel := m.createUserTasksTable()
	// m.mainContent.AddItem(leftPanel, 70, 1, true)
	// m.infoPanel = m.createTaskDetails(completedTasks)
	// m.mainContent.AddItem(m.infoPanel, 0, 1, true)
	// m.app.SetFocus(leftPanel)
}

// -----------------------------------------------------------------------
func (m *BPMNManager) createUserTasksTable() *tview.Flex {

	// Create tasks table
	table := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(true, false).SetBordersColor(tcell.Color100)

	selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)
	table.SetSelectedStyle(selectedStyle)

	// Headers
	headers := []string{"TaskID|", "TaskName", "|TaskDefinitionKey", "|ProcessID", "|Assignee"}
	for i, header := range headers {
		table.SetCell(0, i,
			tview.NewTableCell(header).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}

	tasks, _ := m.apiClient.GetUserTasks()

	// Add data to table
	for row, task := range tasks {
		table.SetCell(row+1, 0, tview.NewTableCell(task.ID+" |"))
		table.SetCell(row+1, 1, tview.NewTableCell(reverseString(task.Name)).SetAlign(tview.AlignRight))
		table.SetCell(row+1, 2, tview.NewTableCell("|"+task.TaskDefinitionKey))
		table.SetCell(row+1, 3, tview.NewTableCell("|"+task.ProcessID))
		statusCell := tview.NewTableCell("|" + task.Assignee)
		statusCell.SetTextColor(tcell.ColorYellow)
		table.SetCell(row+1, 4, statusCell)
	}

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF2:
			selectedRow, _ := table.GetSelection()
			taskId := table.GetCell(selectedRow, 0).Text
			m.CreateModalForTaskCompletion(taskId)
			// m.mainContent.AddItem(m.CreateModalForTaskCompletion(taskId), 1, 0, false)
			return nil
		case tcell.KeyEnter:

			selectedRow, _ := table.GetSelection()
			taskId := strings.ReplaceAll(table.GetCell(selectedRow, 0).Text, "|", "")
			taskId = strings.TrimSpace(taskId)

			taskKey := strings.ReplaceAll(table.GetCell(selectedRow, 2).Text, "|", "")
			taskKey = strings.TrimSpace(taskKey)

			processId := strings.ReplaceAll(table.GetCell(selectedRow, 3).Text, "|", "")
			processId = strings.TrimSpace(processId)

			reDefineDecision := tview.NewCheckbox().
				SetLabel("ReDefineDecision: ").SetChecked(false)
			// reDefineDecision.SetBorderPadding(1, 1, 1, 1)

			dbDecision := tview.NewCheckbox().
				SetLabel("DbDecision: ").SetChecked(false)
			// dbDecision.SetBorderPadding(1, 1, 1, 1)

			businessApproved := tview.NewCheckbox().
				SetLabel("BusinessApproved: ").SetChecked(false)

			technicalApproved := tview.NewCheckbox().
				SetLabel("TechnicalApproved: ").SetChecked(false)

			operationApproved := tview.NewCheckbox().
				SetLabel("OperationApproved: ").SetChecked(false)

			form := tview.NewForm().
				AddTextView("Task Id:", taskId, 10, 1, true, false).
				AddTextView("Task Definition Key:", taskKey, 30, 1, true, false).
				AddTextView("Process Id:", processId, 10, 1, true, false).

				// AddInputField("Name", "", 20, nil, nil).
				// AddInputField("ReDefineDecision", "", 10, nil, nil).
				// AddInputField("DbDecision", "", 10, nil, nil).
				// AddCheckbox("ReDefineDecision", false, nil).
				// AddCheckbox("DbDecision", false, nil).
				// AddTextArea("Task Description", "", 20, 8, 30, nil).
				AddButton("Complete Task", func() {

					data := models.FormData{
						ReDefineDecision:  reDefineDecision.IsChecked(),
						DbDecision:        dbDecision.IsChecked(),
						BusinessApproved:  businessApproved.IsChecked(),
						TechnicalApproved: technicalApproved.IsChecked(),
						OperationApproved: operationApproved.IsChecked(),
						Comment:           "test comment",
						Message:           "Comeleted by BPMN-MANAGER",
					}
					err := m.apiClient.CompleteTask(taskId, data)
					if err != nil {
						m.infoPanel.SetText(err.Error())
					} else {
						lastItem := m.mainContent.GetItem(2)
						m.mainContent.RemoveItem(lastItem)
						m.app.SetFocus(m.mainContent.GetItem(0))
						// m.showMessage("Task Successfully Completed !")
						m.mainContent.AddItem(m.infoPanel, 0, 1, true)
						m.showDashboard()
					}

					// m.pages.SwitchToPage("main")

					// Simulate completing the task by stopping the app
					// Normally, this would trigger a BPMN event like moving to the next task or workflow step
				})
			// SetButtonsAlign(tview.AlignCenter).

			switch taskKey {
			case "Activity_0ol9pgw":
				form.AddFormItem(reDefineDecision).
					AddFormItem(dbDecision)
			case "Activity_0bowttv":
				form.AddFormItem(businessApproved)
			case "Activity_06k5ayj":
				form.AddFormItem(technicalApproved)
			case "Activity_018w7i0":
				form.AddFormItem(operationApproved)
			}

			form.SetBorder(true).SetBorderColor(tcell.Color102)
			form.SetTitle(" User Task Form ")
			// Create a status message to show task completion (could simulate a BPMN task state)
			// statusMessage := tview.NewTextView().
			// 	SetText("Task not completed yet. Fill the form to complete it.").
			// 	SetTextAlign(tview.AlignCenter).
			// 	SetDynamicColors(true).
			// 	SetBorder(true).
			// 	SetTitle("Status")

			// Create a flex layout to arrange the form and status message
			layout := tview.NewFlex().
				SetDirection(tview.FlexRow).
				// AddItem(statusMessage, 0, 1, false).
				AddItem(form, 0, 3, true)

			// m.mainContent.Clear()
			m.mainContent.RemoveItem(m.infoPanel)
			// m.mainContent.AddItem(m.nav, 35, 1, true)
			m.mainContent.AddItem(layout, 0, 1, true)
			m.app.SetFocus(form)
		}
		return event
	})

	// Set up the layout: add the box containing the table to the app
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(table, 0, 2, true)

	flex.SetTitle("Task List")
	flex.SetBorder(true).SetBorderColor(tcell.Color102)

	return flex

}

// -----------------------------------------------------------------------
func (m *BPMNManager) showUserTasks() {
	// Create loading modal
	modal := tview.NewModal().
		SetText("🔄 Loading user tasks...").
		AddButtons([]string{"Cancel"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			m.pages.SwitchToPage("main")
		}
	})
	m.pages.AddPage("user_tasks_loading", modal, true, true)

	// Simulate API call
	go func() {
		// Simulate API delay
		// time.Sleep(2 * time.Second)

		m.app.QueueUpdateDraw(func() {
			// Create tasks table
			table := tview.NewTable().
				SetBorders(true).
				SetFixed(1, 0).
				SetSelectable(true, false)

			// Headers
			headers := []string{"Task", "Process", "Status", "Priority"}
			for i, header := range headers {
				table.SetCell(0, i,
					tview.NewTableCell(header).
						SetTextColor(tcell.ColorYellow).
						SetAlign(tview.AlignCenter).
						SetSelectable(false))
			}

			tasks, _ := m.apiClient.GetUserTasks()

			// Add data to table
			for row, task := range tasks {
				table.SetCell(row+1, 0, tview.NewTableCell(task.ID))
				table.SetCell(row+1, 1, tview.NewTableCell(reverseString(task.Name)).SetAlign(tview.AlignRight))
				table.SetCell(row+1, 2, tview.NewTableCell("In Progress"))
				// table.SetCell(row+1, 3, tview.NewTableCell("High"))
				// switch task.Status {
				// case "completed":
				// 	statusCell.SetTextColor(tcell.ColorGreen)
				// case "in-progress":
				// 	statusCell.SetTextColor(tcell.ColorYellow)
				// default:
				// 	statusCell.SetTextColor(tcell.ColorRed)
				// }
				// statusCell.SetTextColor(tcell.ColorYellow)
				// table.SetCell(row+1, 2, statusCell)
				// table.SetCell(row+1, 3, tview.NewTableCell(task.DueDate.Format("2006-01-02")))
			}

			// Create layout with buttons
			buttons := tview.NewFlex().
				AddItem(tview.NewButton("Complete Task").SetSelectedFunc(func() {
					m.showMessage("Task completion feature would be implemented here")
				}), 0, 1, false).
				AddItem(tview.NewButton("Refresh").SetSelectedFunc(func() {
					m.showUserTasks()
				}), 0, 1, false).
				AddItem(tview.NewButton("Back").SetSelectedFunc(func() {
					m.pages.SwitchToPage("main")
				}), 0, 1, false)

			layout := tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(table, 0, 1, true).
				AddItem(buttons, 1, 1, false)

			layout.SetBorder(true).SetTitle(" My Tasks ")
			m.pages.AddPage("user_tasks", layout, true, true)
			m.pages.SwitchToPage("user_tasks")
		})
	}()
}

// -----------------------------------------------------------------------
func (m *BPMNManager) createRunningProcesses() {

	// Create tasks table
	table := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(true, false).SetBordersColor(tcell.Color100)

	selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)
	table.SetSelectedStyle(selectedStyle)

	// Headers
	headers := []string{"ProcessID", "| ProcessStatus", "| ProcessDefKey", "| StartTime"}
	for i, header := range headers {
		table.SetCell(0, i,
			tview.NewTableCell(header).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}

	processes, _ := m.apiClient.GetRunningProcesses()

	// Add data to table
	for row, process := range processes {
		table.SetCell(row+1, 0, tview.NewTableCell(process.ProcessID))
		table.SetCell(row+1, 1, tview.NewTableCell("| "+process.Status))
		table.SetCell(row+1, 2, tview.NewTableCell("| "+process.ProcessDefinitionKey))
		table.SetCell(row+1, 3, tview.NewTableCell("| "+process.StartTime.Format("2006-01-02 15:04:05")))
	}

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:

			selectedRow, _ := table.GetSelection()
			// m.infoPanel.SetText(strconv.Itoa(selectedRow))
			m.infoPanel.SetText(table.GetCell(selectedRow, 0).Text)
			selectedId := table.GetCell(selectedRow, 0).Text
			details := m.createProcessDetails(selectedId)
			// m.infoPanel.SetBackgroundColor(0x005F87)
			m.infoPanel.SetBackgroundColor(tcell.ColorDarkGreen)
			m.infoPanel.SetDynamicColors(true)
			m.infoPanel.SetTitle("Process Details")
			m.infoPanel.SetText(details)

		}
		return event
	})

	totalText := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	totalText.SetText(fmt.Sprintf("[orange]Total Proceses: %d", len(processes)))
	// Set up the layout: add the box containing the table to the app
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(table, 0, 22, true).
		AddItem(totalText, 0, 1, true)

	flex.SetTitle(" Process List ")
	flex.SetBorder(true).SetBorderColor(tcell.Color102)

	m.mainContent.Clear()
	m.mainContent.AddItem(m.nav, 35, 1, true)
	m.mainContent.AddItem(flex, 0, 3, true)
	m.app.SetFocus(flex)

	m.infoPanel = m.createProcessDetailsPanel(processes)
	m.mainContent.AddItem(m.infoPanel, 0, 3, true)

}

// -----------------------------------------------------------------------
func (m *BPMNManager) showRunningProcesses() {
	modal := tview.NewModal().
		SetText("🔄 Loading running processes...").
		AddButtons([]string{"Cancel"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			m.pages.SwitchToPage("main")
		}
	})
	m.pages.AddPage("processes_loading", modal, true, true)

	go func() {
		time.Sleep(2 * time.Second)

		m.app.QueueUpdateDraw(func() {
			table := tview.NewTable().
				SetBorders(true).
				SetFixed(1, 0).
				SetSelectable(true, false)

			headers := []string{"Process", "Current Activity", "Status", "Started"}
			for i, header := range headers {
				table.SetCell(0, i,
					tview.NewTableCell(header).
						SetTextColor(tcell.ColorYellow).
						SetAlign(tview.AlignCenter).
						SetSelectable(false))
			}

			buttons := tview.NewFlex().
				AddItem(tview.NewButton("View Details").SetSelectedFunc(func() {
					m.showMessage("Process details would be shown here")
				}), 0, 1, false).
				AddItem(tview.NewButton("Refresh").SetSelectedFunc(func() {
					m.showRunningProcesses()
				}), 0, 1, false).
				AddItem(tview.NewButton("Back").SetSelectedFunc(func() {
					m.pages.SwitchToPage("main")
				}), 0, 1, false)

			layout := tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(table, 0, 1, true).
				AddItem(buttons, 1, 1, false)

			layout.SetBorder(true).SetTitle(" Running Processes ")
			m.pages.AddPage("running_processes", layout, true, true)
			m.pages.SwitchToPage("running_processes")
		})
	}()
}

// -----------------------------------------------------------------------
func (m *BPMNManager) showProcessSearch() {
	form := tview.NewForm().
		AddInputField("Process ID", "", 20, nil, nil)
	// form.GetFormItem(0).(*tview.InputField).SetFinishedFunc(func(tcell.Key) {
	// processID := form.GetFormItem(0).(*tview.InputField).GetText()
	// m.showProcessDetails(processID)
	// })

	form.AddButton("Load", func() {
		processID := form.GetFormItem(0).(*tview.InputField).GetText()
		m.showProcessDetails(processID)
		//form.SetFocus(0)

	}).
		AddButton("Back", func() {
			m.pages.SwitchToPage("main")
		})

	form.SetBorder(true).SetTitle(" Enter Process ID ").SetBorderColor(tcell.Color102)
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	flex.AddItem(form, 40, 2, true)
	flex.AddItem(m.infoPanel, 0, 3, false)
	m.pages.AddPage("process_selection", flex, true, true)
	m.pages.SwitchToPage("process_selection")
	m.currentPage = "process_selection"
}

// -----------------------------------------------------------------------
func (m *BPMNManager) showProcessSelection() {
	form := tview.NewForm().
		AddInputField("Process ID", "process-123", 30, nil, nil)

	form.AddButton("Load", func() {
		processID := form.GetFormItem(0).(*tview.InputField).GetText()
		m.showProcessDetails(processID)
	}).
		AddButton("Back", func() {
			m.pages.SwitchToPage("main")
		})

	form.SetBorder(true).SetTitle(" Enter Process ID ")
	m.pages.AddPage("process_selection", form, true, true)
	m.pages.SwitchToPage("process_selection")
}

// -----------------------------------------------------------------------
func (m *BPMNManager) showProcessDetails(processID string) {

	details := m.createProcessDetails(processID)
	m.infoPanel.SetBackgroundColor(tcell.ColorDarkGreen).SetBorderColor(tcell.Color102)
	m.infoPanel.SetDynamicColors(true)
	m.infoPanel.SetTitle("Process Details")
	m.infoPanel.SetText(details)

}

// -----------------------------------------------------------------------
func (m *BPMNManager) showProcessDetailsMock(processID string) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("🔄 Loading details for process: %s", processID)).
		AddButtons([]string{"Cancel"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			m.pages.SwitchToPage("main")
		}
	})
	m.pages.AddPage("details_loading", modal, true, true)

	go func() {
		time.Sleep(2 * time.Second)

		m.app.QueueUpdateDraw(func() {
			detailsText := fmt.Sprintf(`Process: Order Processing
ID: %s
Description: Handles customer order fulfillment
Version: 2.1
Status: Active

Activities:
1. Receive Order
2. Validate Payment
3. Prepare Shipment
4. Deliver Order

Statistics:
• Total Instances: 156
• Running: 12
• Completed: 144
• Success Rate: 92%%`, processID)

			modal := tview.NewModal().
				SetText(detailsText).
				AddButtons([]string{"Start Instance", "View History", "Close"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					switch buttonLabel {
					case "Start Instance":
						m.apiClient.StartProcess()
						m.showMessage("New process instance started!")
					case "View History":
						m.showMessage("Process history would be shown here")
					case "Close":
						m.pages.SwitchToPage("main")
					}
				})

			m.pages.AddPage("process_details", modal, true, true)
			m.pages.SwitchToPage("process_details")
		})
	}()
}

// -----------------------------------------------------------------------
func (m *BPMNManager) showSettings() {
	form := tview.NewForm().
		AddInputField("API Base URL", m.baseURL, 50, nil, nil).
		AddInputField("Auth Token", "your-token-here", 50, nil, nil)

	form.AddButton("Save", func() {
		newURL := form.GetFormItem(0).(*tview.InputField).GetText()
		token := form.GetFormItem(1).(*tview.InputField).GetText()

		m.baseURL = newURL
		m.apiClient = api.NewAPIClient(newURL)
		if token != "" {
			m.apiClient.SetAuthToken(token)
		}
		m.showMessage("Settings saved successfully!")
	}).
		AddButton("Cancel", func() {
			m.pages.SwitchToPage("main")
		})

	form.SetBorder(true).SetTitle(" Settings ")
	m.pages.AddPage("settings", form, true, true)
	m.pages.SwitchToPage("settings")
}

// -----------------------------------------------------------------------
func (m *BPMNManager) showMessage(message string) {
	modal := tview.NewModal().
		SetText("✅ " + message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			m.pages.SwitchToPage("main")
		})
	m.pages.AddPage("message", modal, true, true)
	m.pages.SwitchToPage("message")
}

// -----------------------------------------------------------------------
func (m *BPMNManager) Run() error {
	m.setupUI()

	// Set up global keyboard shortcuts
	m.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF5:
			// m.showProcessDetails("123")
			m.updateDashboardPanel()
			return nil
		case tcell.KeyF3:
			m.showProcessSearch()
			return nil
		case tcell.KeyF2:
			if m.currentPage == "process_selection" {
				m.updateDashboardPanel()
			}
			return nil
		case tcell.KeyCtrlC:
			m.app.Stop()
			return nil
		case tcell.KeyEsc:
			if m.currentPage == "process_selection" {
				return nil
			}
			m.updateDashboardPanel()
			return nil
		}
		return event
	})

	return m.app.Run()
}

// -----------------------------------------------------------------------
func main() {
	// Default API URL
	baseURL := "http://192.168.164.150:8086"
	if len(os.Args) > 1 {
		baseURL = os.Args[1]
	}

	fmt.Printf("Starting BPMN Manager with API: %s\n", baseURL)
	fmt.Println("Initializing UI...")

	// resp, err := http.Get(baseURL + "/api/user/tasks/all")
	// if err != nil {
	// 	fmt.Printf("Error making GET request: %v", err)
	// }
	// defer resp.Body.Close() // Ensure the body is closed after reading.

	// // Check if the response status is OK (200)
	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Printf("Error: HTTP status %d", resp.StatusCode)
	// }

	// // Read the response body
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Printf("Error reading response body: %v", err)
	// }

	// // Print the raw response body (optional)
	// // fmt.Println("Raw Response Body:", string(body))

	// // Parse the response body into a slice of Post structs
	// var posts []models.UserTask
	// err = json.Unmarshal(body, &posts)
	// if err != nil {
	// 	fmt.Printf("Error unmarshalling JSON response: %v", err)
	// }

	// for _, task := range posts {
	// 	fmt.Println(task)
	// }

	manager := NewBPMNManager(baseURL)
	if err := manager.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}

// -----------------------------------------------------------------------
