-- name: create_{{.Model.Name | lower}}_table
-- id: {{.ID}}
-- created: {{.Timestamp}}

-- up:
CREATE TABLE IF NOT EXISTS {{.Model.Name | lower}}s (
    id BIGSERIAL PRIMARY KEY,
    {{range .Model.Fields}}
    {{.Name}} {{.Type}} {{if .Required}}NOT NULL{{end}}{{if .Unique}} UNIQUE{{end}},
    {{end}}
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- down:
DROP TABLE IF EXISTS {{.Model.Name | lower}}s;