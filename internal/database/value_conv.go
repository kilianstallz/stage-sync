package database

import (
	"database/sql"
	"go.uber.org/zap"
)

func ConvertDbValue(v interface{}) interface{} {
	switch v.(type) {
	case sql.NullInt64:
		if v.(sql.NullInt64).Valid {
			return v.(sql.NullInt64).Int64
		}
		return nil
	case sql.NullFloat64:
		if v.(sql.NullFloat64).Valid {
			return v.(sql.NullFloat64).Float64
		}
		return nil
	case sql.NullString:
		if v.(sql.NullString).Valid {
			return v.(sql.NullString).String
		}
		return nil
	case sql.NullBool:
		if v.(sql.NullBool).Valid {
			return v.(sql.NullBool).Bool
		}
		return nil
	case sql.NullByte:
		if v.(sql.NullByte).Valid {
			return v.(sql.NullByte).Byte
		}
		return nil
	case sql.NullTime:
		if v.(sql.NullTime).Valid {
			return v.(sql.NullTime).Time
		}
		return nil
	default:
		zap.L().Error("unsupported type", zap.Any("type", v))
		return v
	}
}
