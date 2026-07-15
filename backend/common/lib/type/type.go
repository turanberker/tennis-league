package customtype

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Sadece Tarih formatı için özel bir tip oluşturuyoruz
type Date time.Time

// MarshalJSON metodu ile bu tipin JSON çıktısını özelleştiriyoruz
func (d Date) MarshalJSON() ([]byte, error) {
	// time.Time tipine dönüştürüp formatlıyoruz
	t := time.Time(d)
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))
	return []byte(formatted), nil
}

// UnmarshalJSON metodu (İsteğe bağlı: JSON'dan okuma yaparken de bu formatı kabul etmek için)
func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Tırnak işaretlerini temizliyoruz
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// ==========================================
// DATABASE ENTEGRASYONU (Valuer ve Scanner)
// ==========================================

// 1. Value: Go'dan Database'e veri yazarken çalışır.
// Squirrel veya SQL driver bu metot sayesinde veriyi "time.Time" formatında DB'ye gönderir.
func (d Date) Value() (driver.Value, error) {
	return time.Time(d), nil
}

// 2. Scan: Database'den Go'ya veri okurken (Query/Scan yaparken) çalışır.
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch t := value.(type) {
	case time.Time:
		*d = Date(t)
		return nil
	case string: // Bazı DB driver'ları tarihi string dönebilir
		parsedTime, err := time.Parse("2006-01-02", t)
		if err != nil {
			return err
		}
		*d = Date(parsedTime)
		return nil
	default:
		return fmt.Errorf("customtype.Date için desteklenmeyen tip: %T", value)
	}
}
