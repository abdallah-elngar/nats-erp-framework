package models

import (
    "time"

    "gorm.io/gorm"
)

// {{.Model.Name}}Relations يمثل علاقات {{.Model.Name}}
type {{.Model.Name}}Relations struct {
    {{range .Model.Relations}}
    // {{.Type}} relation with {{.Model}}
    {{.Model}} {{.Model}} `gorm:"foreignKey:{{.ForeignKey}}"`
    {{end}}
}

// GetRelations يعيد علاقات النموذج
func (m *{{.Model.Name}}) GetRelations() *{{.Model.Name}}Relations {
    return &{{.Model.Name}}Relations{
        {{range .Model.Relations}}
        {{.Model}}: m.{{.Model}},
        {{end}}
    }
}

// LoadRelations يحمل العلاقات
func (m *{{.Model.Name}}) LoadRelations(db *gorm.DB) error {
    {{range .Model.Relations}}
    if err := db.Model(m).Association("{{.Model}}").Find(&m.{{.Model}}); err != nil {
        return err
    }
    {{end}}
    return nil
}