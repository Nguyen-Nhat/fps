// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/processingfile"
)

// ProcessingFile is the model entity for the ProcessingFile schema.
type ProcessingFile struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// ClientID holds the value of the "client_id" field.
	ClientID int32 `json:"client_id,omitempty"`
	// DisplayName holds the value of the "display_name" field.
	DisplayName string `json:"display_name,omitempty"`
	// FileURL holds the value of the "file_url" field.
	FileURL string `json:"file_url,omitempty"`
	// ResultFileURL holds the value of the "result_file_url" field.
	ResultFileURL string `json:"result_file_url,omitempty"`
	// Init=1; Processing=2; Failed=3; Finished=4
	Status int16 `json:"status,omitempty"`
	// TotalMapping holds the value of the "total_mapping" field.
	TotalMapping int32 `json:"total_mapping,omitempty"`
	// StatsTotalRow holds the value of the "stats_total_row" field.
	StatsTotalRow int32 `json:"stats_total_row,omitempty"`
	// StatsTotalSuccess holds the value of the "stats_total_success" field.
	StatsTotalSuccess int32 `json:"stats_total_success,omitempty"`
	// ErrorDisplay holds the value of the "error_display" field.
	ErrorDisplay string `json:"error_display,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// CreatedBy holds the value of the "created_by" field.
	CreatedBy string `json:"created_by,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ProcessingFile) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case processingfile.FieldID, processingfile.FieldClientID, processingfile.FieldStatus, processingfile.FieldTotalMapping, processingfile.FieldStatsTotalRow, processingfile.FieldStatsTotalSuccess:
			values[i] = new(sql.NullInt64)
		case processingfile.FieldDisplayName, processingfile.FieldFileURL, processingfile.FieldResultFileURL, processingfile.FieldErrorDisplay, processingfile.FieldCreatedBy:
			values[i] = new(sql.NullString)
		case processingfile.FieldCreatedAt, processingfile.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type ProcessingFile", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ProcessingFile fields.
func (pf *ProcessingFile) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case processingfile.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pf.ID = int(value.Int64)
		case processingfile.FieldClientID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field client_id", values[i])
			} else if value.Valid {
				pf.ClientID = int32(value.Int64)
			}
		case processingfile.FieldDisplayName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field display_name", values[i])
			} else if value.Valid {
				pf.DisplayName = value.String
			}
		case processingfile.FieldFileURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field file_url", values[i])
			} else if value.Valid {
				pf.FileURL = value.String
			}
		case processingfile.FieldResultFileURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field result_file_url", values[i])
			} else if value.Valid {
				pf.ResultFileURL = value.String
			}
		case processingfile.FieldStatus:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				pf.Status = int16(value.Int64)
			}
		case processingfile.FieldTotalMapping:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field total_mapping", values[i])
			} else if value.Valid {
				pf.TotalMapping = int32(value.Int64)
			}
		case processingfile.FieldStatsTotalRow:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field stats_total_row", values[i])
			} else if value.Valid {
				pf.StatsTotalRow = int32(value.Int64)
			}
		case processingfile.FieldStatsTotalSuccess:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field stats_total_success", values[i])
			} else if value.Valid {
				pf.StatsTotalSuccess = int32(value.Int64)
			}
		case processingfile.FieldErrorDisplay:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field error_display", values[i])
			} else if value.Valid {
				pf.ErrorDisplay = value.String
			}
		case processingfile.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pf.CreatedAt = value.Time
			}
		case processingfile.FieldCreatedBy:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field created_by", values[i])
			} else if value.Valid {
				pf.CreatedBy = value.String
			}
		case processingfile.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				pf.UpdatedAt = value.Time
			}
		}
	}
	return nil
}

// Update returns a builder for updating this ProcessingFile.
// Note that you need to call ProcessingFile.Unwrap() before calling this method if this ProcessingFile
// was returned from a transaction, and the transaction was committed or rolled back.
func (pf *ProcessingFile) Update() *ProcessingFileUpdateOne {
	return (&ProcessingFileClient{config: pf.config}).UpdateOne(pf)
}

// Unwrap unwraps the ProcessingFile entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pf *ProcessingFile) Unwrap() *ProcessingFile {
	_tx, ok := pf.config.driver.(*txDriver)
	if !ok {
		panic("ent: ProcessingFile is not a transactional entity")
	}
	pf.config.driver = _tx.drv
	return pf
}

// String implements the fmt.Stringer.
func (pf *ProcessingFile) String() string {
	var builder strings.Builder
	builder.WriteString("ProcessingFile(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pf.ID))
	builder.WriteString("client_id=")
	builder.WriteString(fmt.Sprintf("%v", pf.ClientID))
	builder.WriteString(", ")
	builder.WriteString("display_name=")
	builder.WriteString(pf.DisplayName)
	builder.WriteString(", ")
	builder.WriteString("file_url=")
	builder.WriteString(pf.FileURL)
	builder.WriteString(", ")
	builder.WriteString("result_file_url=")
	builder.WriteString(pf.ResultFileURL)
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", pf.Status))
	builder.WriteString(", ")
	builder.WriteString("total_mapping=")
	builder.WriteString(fmt.Sprintf("%v", pf.TotalMapping))
	builder.WriteString(", ")
	builder.WriteString("stats_total_row=")
	builder.WriteString(fmt.Sprintf("%v", pf.StatsTotalRow))
	builder.WriteString(", ")
	builder.WriteString("stats_total_success=")
	builder.WriteString(fmt.Sprintf("%v", pf.StatsTotalSuccess))
	builder.WriteString(", ")
	builder.WriteString("error_display=")
	builder.WriteString(pf.ErrorDisplay)
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(pf.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("created_by=")
	builder.WriteString(pf.CreatedBy)
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(pf.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// ProcessingFiles is a parsable slice of ProcessingFile.
type ProcessingFiles []*ProcessingFile

func (pf ProcessingFiles) config(cfg config) {
	for _i := range pf {
		pf[_i].config = cfg
	}
}