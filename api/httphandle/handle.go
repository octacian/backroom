package httphandle

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/octacian/backroom/api/cage"
	"github.com/octacian/backroom/api/db"
	"github.com/octacian/backroom/api/hook"
)

// requestCreateRecord is the request body for creating a new caged record.
type requestCreateRecord struct {
	Cage string   `json:"cage"`
	Data db.JSONB `json:"data"`
}

// responseDelete is the response body for deleting caged record(s).
type responseDelete struct {
	Success bool `json:"success"`
	Deleted int  `json:"deleted"`
}

// wrapHookRunner just vibes for me.
func wrapHookRunner(w http.ResponseWriter, action hook.Action, record *cage.Record) bool {
	if err := hook.RunHooksByAction(action, record); err != nil {
		slog.Error("Failed to run hooks", "action", action, "error", err)
		http.Error(w, fmt.Sprintf("Failed to run %s hooks", action), http.StatusInternalServerError)
		return false
	}
	return true
}

// HandleCreateRecord handles the creation of a new caged record. Expects
// a JSON payload with the record data. Returns the created record as JSON.
func HandleCreateRecord(w http.ResponseWriter, r *http.Request) {
	var req requestCreateRecord
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.Cage == "" {
		http.Error(w, "Missing cage key", http.StatusBadRequest)
		return
	}

	record := cage.NewRecord(req.Cage, req.Data)
	if err := cage.CreateRecord(record); err != nil {
		http.Error(w, "Failed to create record", http.StatusInternalServerError)
		return
	}

	// Run hooks after creating the record
	if ok := wrapHookRunner(w, hook.ActionCreate, record); !ok {
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(record)
}

// HandleGetRecord handles the retrieval of a caged record by its UUID.
// Expects the UUID as a URL parameter. Returns the record as JSON.
func HandleGetRecord(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")
	if uuidStr == "" {
		http.Error(w, "Missing UUID", http.StatusBadRequest)
		return
	}

	uuid, err := db.ParseUUID(uuidStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	record, err := cage.GetRecord(uuid)
	if err != nil {
		http.Error(w, "Failed to retrieve record", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(record)
}

// HandleListRecordsByCage handles the retrieval of all caged records by their key.
// Expects the key as a URL parameter. Returns the records as JSON.
func HandleListRecordsByCage(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		http.Error(w, "Missing cage key", http.StatusBadRequest)
		return
	}

	records, err := cage.ListRecordsByCage(key)
	if err != nil {
		http.Error(w, "Failed to retrieve records", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(records)
}

// HandleListCages handles the retrieval of all unique cage keys.
// Returns the keys as JSON.
func HandleListCages(w http.ResponseWriter, r *http.Request) {
	keys, err := cage.ListCages()
	if err != nil {
		http.Error(w, "Failed to retrieve cages", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(keys)
}

// HandleUpdateRecord handles the update of an existing caged record.
// Expects the UUID as a URL parameter and a JSON payload matching requestCreateRecord.
// Returns the updated record as JSON.
func HandleUpdateRecord(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")
	if uuidStr == "" {
		http.Error(w, "Missing UUID", http.StatusBadRequest)
		return
	}

	uuid, err := db.ParseUUID(uuidStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	var req requestCreateRecord
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.Cage == "" {
		http.Error(w, "Missing cage key", http.StatusBadRequest)
		return
	}

	record, err := cage.GetRecord(uuid)
	if err != nil {
		http.Error(w, "Failed to retrieve record", http.StatusInternalServerError)
		return
	}
	record.Data = req.Data

	if err := cage.UpdateRecord(record); err != nil {
		http.Error(w, "Failed to update record", http.StatusInternalServerError)
		return
	}

	// Run update hooks after updating the record
	if ok := wrapHookRunner(w, hook.ActionUpdate, record); !ok {
		return
	}

	json.NewEncoder(w).Encode(record)
}

// HandleDeleteRecord handles the deletion of a caged record by its UUID.
// Expects the UUID as a URL parameter.
// Returns a success message as JSON.
func HandleDeleteRecord(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")
	if uuidStr == "" {
		http.Error(w, "Missing UUID", http.StatusBadRequest)
		return
	}

	uuid, err := db.ParseUUID(uuidStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	record, err := cage.GetRecord(uuid)
	if err != nil {
		http.Error(w, "Failed to retrieve record", http.StatusInternalServerError)
		return
	}

	if err := cage.DeleteRecord(uuid); err != nil {
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	// Run delete hooks after deleting the record
	if ok := wrapHookRunner(w, hook.ActionDelete, record); !ok {
		return
	}

	response := responseDelete{
		Success: true,
		Deleted: 1,
	}
	json.NewEncoder(w).Encode(response)
}

// HandleDeleteRecordsByKey handles the deletion of all caged records by their key.
// Expects the key as a URL parameter.
// Returns a success message as JSON.
func HandleDeleteRecordsByKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		http.Error(w, "Missing cage key", http.StatusBadRequest)
		return
	}

	deleted, err := cage.DeleteCage(key)
	if err != nil {
		http.Error(w, "Failed to delete records", http.StatusInternalServerError)
		return
	}

	response := responseDelete{
		Success: true,
		Deleted: int(deleted),
	}
	json.NewEncoder(w).Encode(response)
}
