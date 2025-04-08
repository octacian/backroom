package cage

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/octacian/backroom/api/db"
)

// requestCreateRecord is the request body for creating a new caged record.
type requestCreateRecord struct {
	Key  string   `json:"key"`
	Data db.JSONB `json:"data"`
}

// responseDelete is the response body for deleting caged record(s).
type responseDelete struct {
	Success bool `json:"success"`
	Deleted int  `json:"deleted"`
}

// HandleCreateRecord handles the creation of a new caged record. Expects
// a JSON payload with the record data. Returns the created record as JSON.
// Cage key is passed as a URL parameter called `key`.
func HandleCreateRecord(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		http.Error(w, "Missing cage key", http.StatusBadRequest)
		return
	}

	var req requestCreateRecord
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.Key != key {
		http.Error(w, "Cage key mismatch", http.StatusBadRequest)
		return
	}

	cage := NewRecord(key, req.Data)
	if err := CreateRecord(cage); err != nil {
		http.Error(w, "Failed to create record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cage)
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

	cage, err := GetRecord(uuid)
	if err != nil {
		http.Error(w, "Failed to retrieve record", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cage)
}

// HandleListRecordsByKey handles the retrieval of all caged records by their key.
// Expects the key as a URL parameter. Returns the records as JSON.
func HandleListRecordsByKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	if key == "" {
		http.Error(w, "Missing cage key", http.StatusBadRequest)
		return
	}

	cages, err := ListRecordsByKey(key)
	if err != nil {
		http.Error(w, "Failed to retrieve records", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cages)
}

// HandleListKeys handles the retrieval of all unique cage keys.
// Returns the keys as JSON.
func HandleListKeys(w http.ResponseWriter, r *http.Request) {
	keys, err := ListCageKeys()
	if err != nil {
		http.Error(w, "Failed to retrieve keys", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(keys)
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

	if err := DeleteRecord(uuid); err != nil {
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
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

	deleted, err := DeleteCage(key)
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
