package database

import (
	"database/sql"
	numeric "github.com/jackc/pgtype/ext/shopspring-numeric"
	"github.com/shopspring/decimal"
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
	case numeric.Numeric:
		v, err := v.(numeric.Numeric).Value()
		if err != nil {
			zap.S().Errorw("Error converting numeric value", "error", err)
			return nil
		}
		return v
	case decimal.NullDecimal:
		if v.(decimal.NullDecimal).Valid {
			return v.(decimal.NullDecimal).Decimal
		}
		v, err := v.(decimal.NullDecimal).Value()
		if err != nil {
			zap.S().Errorw("Error converting numeric value", "error", err)
			return nil
		}
		return v
	default:
		zap.L().Error("unsupported type", zap.Any("type", v))
		return v
	}
}
