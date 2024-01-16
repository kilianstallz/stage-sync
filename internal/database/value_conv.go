package database

import (
	"database/sql"
	"github.com/jackc/pgtype"
	numeric "github.com/jackc/pgtype/ext/shopspring-numeric"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func ConvertDbValue(v interface{}) interface{} {
	switch t := v.(type) {
	case sql.NullInt64:
		if t.Valid {
			return t.Int64
		}
		return nil
	case sql.NullFloat64:
		if t.Valid {
			return t.Float64
		}
		return nil
	case sql.NullString:
		if t.Valid {
			return t.String
		}
		return nil
	case sql.NullBool:
		if t.Valid {
			return t.Bool
		}
		return nil
	case sql.NullByte:
		if t.Valid {
			return t.Byte
		}
		return nil
	case sql.NullTime:
		if t.Valid {
			return t.Time
		}
		return nil
	case numeric.Numeric:
		v, err := t.Value()
		if err != nil {
			zap.S().Errorw("Error converting numeric value", "error", err)
			return nil
		}
		return v
	case decimal.NullDecimal:
		if t.Valid {
			return t.Decimal
		}
		v, err := t.Value()
		if err != nil {
			zap.S().Errorw("Error converting numeric value", "error", err)
			return nil
		}
		return v
	case pgtype.Int4Array:
		return pgIntArrayToSlice(t)
	default:
		zap.L().Error("unsupported type", zap.Any("type", v))
		return v
	}
}

func pgIntArrayToSlice(pgArray pgtype.Int4Array) []int {
	result := make([]int, len(pgArray.Elements))
	for i, elem := range pgArray.Elements {
		// Use Int field to access the actual int value
		result[i] = int(elem.Int)
	}
	return result
}
