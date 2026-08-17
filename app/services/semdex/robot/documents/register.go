package documents

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/rs/xid"
	adksession "google.golang.org/adk/v2/session"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

const (
	registerStateKey = "document_register"
	maxOpenDocuments = 8
	RootNodeID       = "root"
)

type SourceType string

const (
	SourceTypeLibraryPage SourceType = "library_page"
	SourceTypeThread      SourceType = "thread"
	SourceTypeWeb         SourceType = "web"
)

type Snapshot struct {
	DocumentID       string     `json:"document_id"`
	SourceType       SourceType `json:"source_type"`
	SourceID         string     `json:"source_id"`
	Title            string     `json:"title"`
	Content          string     `json:"content"`
	ActiveNodeID     string     `json:"active_node_id"`
	ActivePage       int        `json:"active_page"`
	ActiveTotalPages int        `json:"active_total_pages"`
	ActiveItemStart  int        `json:"active_item_start"`
	ActiveItemEnd    int        `json:"active_item_end"`
	ActiveTotalItems int        `json:"active_total_items"`
}

type Register struct {
	ActiveDocumentID string     `json:"active_document_id"`
	Documents        []Snapshot `json:"documents"`
}

type DocumentInfo struct {
	DocumentID   string
	SourceType   SourceType
	SourceID     string
	Title        string
	Active       bool
	ActiveNodeID string
	Page         int
	TotalPages   int
	ItemStart    int
	ItemEnd      int
	TotalItems   int
}

var (
	ErrNoActiveDocument = errors.New("no document is open")
	ErrDocumentNotFound = errors.New("open document not found")
	ErrNodeNotFound     = errors.New("document location not found")
)

func Open(state adksession.State, sourceType SourceType, sourceID, title string, content datagraph.Content) (Projection, error) {
	sourceID = strings.TrimSpace(sourceID)
	if sourceID == "" {
		return Projection{}, fmt.Errorf("document source identifier is required")
	}

	stable, err := datagraph.NewRichTextWithNewBlocks(content)
	if err != nil {
		return Projection{}, fmt.Errorf("prepare document structure: %w", err)
	}

	register, err := loadRegister(state)
	if err != nil {
		return Projection{}, err
	}

	if strings.TrimSpace(title) == "" {
		title = sourceID
	}

	snapshot := Snapshot{
		SourceType: sourceType,
		SourceID:   sourceID,
		Title:      title,
		Content:    stable.HTML(),
	}

	replaced := -1
	for i, existing := range register.Documents {
		if existing.SourceType == sourceType && existing.SourceID == sourceID {
			snapshot.DocumentID = existing.DocumentID
			replaced = i
			break
		}
	}
	if snapshot.DocumentID == "" {
		snapshot.DocumentID = "doc_" + xid.New().String()
	}
	projection, err := project(snapshot, RootNodeID, 1)
	if err != nil {
		return Projection{}, err
	}
	setSnapshotCursor(&snapshot, projection)

	if replaced >= 0 {
		register.Documents = append(register.Documents[:replaced], register.Documents[replaced+1:]...)
	}
	register.Documents = append(register.Documents, snapshot)
	if len(register.Documents) > maxOpenDocuments {
		register.Documents = register.Documents[len(register.Documents)-maxOpenDocuments:]
	}
	register.ActiveDocumentID = snapshot.DocumentID

	if err := state.Set(registerStateKey, register); err != nil {
		return Projection{}, fmt.Errorf("save document register: %w", err)
	}

	return projection, nil
}

func Get(state adksession.State, documentID, nodeID string, page int) (Projection, error) {
	register, err := loadRegister(state)
	if err != nil {
		return Projection{}, err
	}
	index, err := resolveSnapshotIndex(register, documentID)
	if err != nil {
		return Projection{}, err
	}
	snapshot := register.Documents[index]
	if strings.TrimSpace(nodeID) == "" {
		nodeID = snapshot.ActiveNodeID
		if !trustedNodeID(nodeID) {
			nodeID = RootNodeID
		}
		if page <= 0 {
			page = max(1, snapshot.ActivePage)
		}
	} else if page <= 0 {
		page = 1
	}
	projection, err := project(snapshot, nodeID, page)
	if err != nil {
		return Projection{}, err
	}
	setSnapshotCursor(&register.Documents[index], projection)
	register.ActiveDocumentID = register.Documents[index].DocumentID
	if err := state.Set(registerStateKey, register); err != nil {
		return Projection{}, fmt.Errorf("save document register: %w", err)
	}
	return projection, nil
}

func setSnapshotCursor(snapshot *Snapshot, projection Projection) {
	snapshot.ActiveNodeID = projection.NodeID
	snapshot.ActivePage = projection.Page
	snapshot.ActiveTotalPages = projection.TotalPages
	snapshot.ActiveItemStart = projection.ItemStart
	snapshot.ActiveItemEnd = projection.ItemEnd
	snapshot.ActiveTotalItems = projection.TotalItems
}

func List(state adksession.ReadonlyState) ([]DocumentInfo, error) {
	register, err := loadRegister(state)
	if err != nil {
		return nil, err
	}

	items := make([]DocumentInfo, 0, len(register.Documents))
	for _, snapshot := range register.Documents {
		items = append(items, DocumentInfo{
			DocumentID:   snapshot.DocumentID,
			SourceType:   snapshot.SourceType,
			SourceID:     snapshot.SourceID,
			Title:        snapshot.Title,
			Active:       snapshot.DocumentID == register.ActiveDocumentID,
			ActiveNodeID: snapshot.ActiveNodeID,
			Page:         max(1, snapshot.ActivePage),
			TotalPages:   max(1, snapshot.ActiveTotalPages),
			ItemStart:    snapshot.ActiveItemStart,
			ItemEnd:      snapshot.ActiveItemEnd,
			TotalItems:   snapshot.ActiveTotalItems,
		})
	}
	return items, nil
}

func RegisterInstruction(state adksession.ReadonlyState) (string, error) {
	register, err := loadRegister(state)
	if err != nil {
		return "", err
	}

	ids := make([]string, 0, len(register.Documents))
	for _, snapshot := range register.Documents {
		if trustedDocumentID(snapshot.DocumentID) {
			ids = append(ids, snapshot.DocumentID)
		}
	}
	if len(ids) == 0 {
		return "", nil
	}

	active := "none"
	activeNode := "none"
	activePage := 1
	activeTotalPages := 1
	activeItemStart := 0
	activeItemEnd := 0
	activeTotalItems := 0
	if trustedDocumentID(register.ActiveDocumentID) {
		for _, snapshot := range register.Documents {
			if snapshot.DocumentID == register.ActiveDocumentID {
				active = snapshot.DocumentID
				activeNode = RootNodeID
				if trustedNodeID(snapshot.ActiveNodeID) {
					activeNode = snapshot.ActiveNodeID
				}
				activePage = max(1, snapshot.ActivePage)
				activeTotalPages = max(activePage, snapshot.ActiveTotalPages)
				activeItemStart = snapshot.ActiveItemStart
				activeItemEnd = snapshot.ActiveItemEnd
				activeTotalItems = snapshot.ActiveTotalItems
				break
			}
		}
	}

	var instruction strings.Builder
	instruction.WriteString("### Document register\n\n")
	instruction.WriteString(fmt.Sprintf("Active document: `%s`\n", active))
	instruction.WriteString(fmt.Sprintf("Current node: `%s`\n", activeNode))
	if activeTotalPages > 1 {
		instruction.WriteString(fmt.Sprintf("Page: %d/%d\n", activePage, activeTotalPages))
	}
	if activeTotalItems > 0 {
		instruction.WriteString(fmt.Sprintf("Items: %d-%d of %d\n", activeItemStart, activeItemEnd, activeTotalItems))
	}
	instruction.WriteString(fmt.Sprintf("Open documents: %d\n", len(ids)))
	instruction.WriteString("Open document IDs, oldest to newest:\n")
	for _, id := range ids {
		instruction.WriteString("- `")
		instruction.WriteString(id)
		instruction.WriteString("`")
		if id == active {
			instruction.WriteString(" (active)")
		}
		instruction.WriteString("\n")
	}
	instruction.WriteString("Document titles, sources, and content are intentionally omitted from system state. Use document_list or document_get to retrieve them as tool output.")

	return instruction.String(), nil
}

func trustedDocumentID(value string) bool {
	if !strings.HasPrefix(value, "doc_") {
		return false
	}
	_, err := xid.FromString(strings.TrimPrefix(value, "doc_"))
	return err == nil
}

func trustedNodeID(value string) bool {
	if value == RootNodeID {
		return true
	}
	if !strings.HasPrefix(value, "sdb_") {
		return false
	}
	_, err := xid.FromString(strings.TrimPrefix(value, "sdb_"))
	return err == nil
}

func Close(state adksession.State, documentID string) (closed string, active string, remaining int, err error) {
	register, err := loadRegister(state)
	if err != nil {
		return "", "", 0, err
	}
	if strings.TrimSpace(documentID) == "" {
		documentID = register.ActiveDocumentID
	}
	if documentID == "" {
		return "", "", 0, ErrNoActiveDocument
	}

	index := -1
	for i, snapshot := range register.Documents {
		if snapshot.DocumentID == documentID {
			index = i
			break
		}
	}
	if index < 0 {
		return "", "", 0, fmt.Errorf("%w: %s", ErrDocumentNotFound, documentID)
	}

	register.Documents = append(register.Documents[:index], register.Documents[index+1:]...)
	if register.ActiveDocumentID == documentID {
		register.ActiveDocumentID = ""
		if len(register.Documents) > 0 {
			register.ActiveDocumentID = register.Documents[len(register.Documents)-1].DocumentID
		}
	}

	if err := state.Set(registerStateKey, register); err != nil {
		return "", "", 0, fmt.Errorf("save document register: %w", err)
	}

	return documentID, register.ActiveDocumentID, len(register.Documents), nil
}

func resolveSnapshot(state adksession.ReadonlyState, documentID string) (Snapshot, error) {
	register, err := loadRegister(state)
	if err != nil {
		return Snapshot{}, err
	}
	index, err := resolveSnapshotIndex(register, documentID)
	if err != nil {
		return Snapshot{}, err
	}
	return register.Documents[index], nil
}

func resolveSnapshotIndex(register Register, documentID string) (int, error) {
	if strings.TrimSpace(documentID) == "" {
		documentID = register.ActiveDocumentID
	}
	if documentID == "" {
		return -1, ErrNoActiveDocument
	}
	for index, snapshot := range register.Documents {
		if snapshot.DocumentID == documentID {
			return index, nil
		}
	}
	return -1, fmt.Errorf("%w: %s", ErrDocumentNotFound, documentID)
}

func loadRegister(state adksession.ReadonlyState) (Register, error) {
	value, err := state.Get(registerStateKey)
	if err != nil {
		if errors.Is(err, adksession.ErrStateKeyNotExist) {
			return Register{}, nil
		}
		return Register{}, fmt.Errorf("read document register: %w", err)
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return Register{}, fmt.Errorf("encode document register: %w", err)
	}
	var register Register
	if err := json.Unmarshal(encoded, &register); err != nil {
		return Register{}, fmt.Errorf("decode document register: %w", err)
	}
	for index := range register.Documents {
		snapshot := &register.Documents[index]
		if !trustedNodeID(snapshot.ActiveNodeID) {
			snapshot.ActiveNodeID = RootNodeID
		}
		snapshot.ActivePage = max(1, snapshot.ActivePage)
		snapshot.ActiveTotalPages = max(snapshot.ActivePage, snapshot.ActiveTotalPages)
		if snapshot.ActiveTotalItems <= 0 {
			snapshot.ActiveItemStart = 0
			snapshot.ActiveItemEnd = 0
			snapshot.ActiveTotalItems = 0
		} else {
			snapshot.ActiveItemStart = max(1, snapshot.ActiveItemStart)
			snapshot.ActiveItemEnd = min(snapshot.ActiveTotalItems, max(snapshot.ActiveItemStart, snapshot.ActiveItemEnd))
		}
	}
	return register, nil
}
